package backend

import (
	"fmt"
	"github.com/quarnster/parser"
	"io/ioutil"
	"lime/backend/loaders"
	. "lime/backend/primitives"
	"lime/backend/textmate"
	"sort"
)

type (
	View struct {
		HasSettings
		HasId
		name      string
		window    *Window
		buffer    *Buffer
		selection RegionSet
		undoStack undoStack
		scratch   bool
		syntax    textmate.LanguageParser
	}
	Edit struct {
		composite CompositeAction
		savedSel  RegionSet
		v         *View
	}
)

func newEdit(v *View) *Edit {
	ret := &Edit{
		v: v,
	}
	for _, r := range v.Sel().Regions() {
		ret.savedSel.Add(r)
	}
	return ret
}

func (e *Edit) Apply() {
	e.composite.Apply()
}

func (e *Edit) Undo() {
	e.composite.Undo()
	e.v.Sel().Clear()
	for _, r := range e.savedSel.Regions() {
		e.v.Sel().Add(r)
	}
}

func (v *View) setBuffer(b *Buffer) error {
	if v.buffer != nil {
		return fmt.Errorf("There is already a buffer set")
	}
	v.buffer = b
	// TODO(q): Dynamically load the correct syntax file
	v.syntax.Language = &textmate.Language{}
	b.AddCallback(v.flush)
	return nil
}

func (v *View) flush(a, b int) {
	v.selection.Adjust(a, b)
	// TODO(q): A full reparse every time the buffer changes is overkill.
	// It would be better if the nodes are just adjusted as appropriate, together with a
	// minimal parse of the new data
	v.syntax.Parse(v.buffer.Data())
	OnModified.Call(v)
	OnSelectionModified.Call(v)
}

func (v *View) SetSyntaxFile(f string) error {
	var lang textmate.Language
	if d, err := ioutil.ReadFile(f); err != nil {
		return err
	} else if err := loaders.LoadPlist(d, &lang); err != nil {
		return err
	} else {
		v.syntax.Language = &lang
		v.syntax.Parse(v.buffer.Data())
	}
	return nil
}

func (v *View) Sel() *RegionSet {
	return &v.selection
}

func (v *View) Window() *Window {
	return v.window
}

func (v *View) Buffer() *Buffer {
	return v.buffer
}

func (v *View) Insert(edit *Edit, point int, value string) {
	edit.composite.AddExec(NewInsertAction(v.buffer, point, value))
}

func (v *View) Erase(edit *Edit, r Region) {
	edit.composite.AddExec(NewEraseAction(v.buffer, r))
}

func (v *View) Replace(edit *Edit, r Region, value string) {
	edit.composite.AddExec(NewReplaceAction(v.buffer, r, value))
}

func (v *View) BeginEdit() *Edit {
	return newEdit(v)
}

func (v *View) EndEdit(e *Edit) {
	if !v.scratch {
		v.undoStack.Add(e, true)
	}
}

func (v *View) Size() int {
	return v.buffer.Size()
}

func (v *View) Substr(r Region) string {
	return v.buffer.Substr(r)
}

func (v *View) SetScratch(s bool) {
	v.scratch = s
}

func (v *View) IsScratch() bool {
	return v.scratch
}

func findScope(search parser.Range, node *parser.Node, in string) string {
	idx := sort.Search(len(node.Children), func(i int) bool {
		return node.Children[i].Range.Start >= search.Start || node.Children[i].Range.Contains(search)
	})
	for idx < len(node.Children) {
		c := node.Children[idx]
		if c.Range.Start > search.End {
			break
		}
		if c.Range.Contains(search) {
			in += " " + node.Name
			return findScope(search, node.Children[idx], in)
		}
		idx++
	}
	if node.Range.Contains(search) {
		return in + " " + node.Name
	}
	return in
}

func (v *View) ScopeName(point int) string {
	return findScope(parser.Range{point, point + 1}, v.syntax.RootNode(), "")
}
