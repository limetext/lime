package backend

import (
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
		CompositeAction
	}
)

func (v *View) setBuffer(b *Buffer) {
	v.buffer = b
	b.AddCallback(v.selection.Adjust)
	// TODO(q): Dynamically load the correct syntax file
	v.syntax.Language = &textmate.Language{}
	b.AddCallback(v.reparse)
}

func (v *View) reparse(a, b int) {
	v.syntax.Parse(v.buffer.Data())
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
	edit.AddExec(NewInsertAction(v.buffer, point, value))
}

func (v *View) Erase(edit *Edit, r Region) {
	edit.AddExec(NewEraseAction(v.buffer, r))
}

func (v *View) Replace(edit *Edit, r Region, value string) {
	edit.AddExec(NewReplaceAction(v.buffer, r, value))
}

func (v *View) BeginEdit() *Edit {
	return &Edit{}
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
