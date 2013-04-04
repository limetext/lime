package backend

import (
	"bytes"
	"fmt"
	"github.com/quarnster/parser"
	"io/ioutil"
	"lime/backend/loaders"
	. "lime/backend/primitives"
	"lime/backend/textmate"
	"regexp"
	"sort"
	"strings"
)

type (
	View struct {
		HasSettings
		HasId
		name          string
		window        *Window
		buffer        *Buffer
		selection     RegionSet
		undoStack     undoStack
		scratch       bool
		syntax        textmate.LanguageParser
		lastScopeNode *parser.Node
		lastScopeBuf  bytes.Buffer
		lastScopeName string
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
	v.lastScopeNode = nil
	v.lastScopeBuf.Reset()
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

func (v *View) RowCol(point int) (row, col int) {
	lines := strings.Split(v.Buffer().Data()[:point], "\n")
	if l := len(lines); l == 0 {
		return 1, 1
	} else {
		return l, len(lines[l-1]) + 1
	}
}

func (v *View) TextPoint(row, col int) int {
	lines := strings.Split(v.buffer.Data(), "\n")
	if row == 0 {
		return 0
	}
	if col == 0 {
		col = 1
	}
	if row == 1 {
		col -= 1
	} else if row > len(lines) {
		return v.Size()
	}
	offset := len(strings.Join(lines[:row-1], "\n")) + col
	return offset
}

func (v *View) Line(offset int) Region {
	if offset < 0 {
		return Region{0, 0}
	} else if s := v.Size(); offset >= s {
		return Region{s, s}
	} else if v.buffer.Data()[offset] == '\n' {
		return Region{offset, offset}
	}
	data := v.buffer.Data()
	s := offset
	for s > 0 && data[s-1] != '\n' {
		s--
	}
	e := offset + 1
	for e < len(data) && data[e] != '\n' {
		e++
	}
	return Region{s, e}
}

var (
	vwre1 = regexp.MustCompile(`\b\w*$`)
	vwre2 = regexp.MustCompile(`^\w*`)
)

func (v *View) Word(offset int) Region {
	_, col := v.RowCol(offset)
	lr := v.Line(offset)
	line := v.Substr(lr)
	begin := 0
	end := len(line)

	if col > len(line) {
		col = len(line)
	}
	if m := vwre1.FindStringIndex(line[:col]); m != nil {
		begin = m[0]
	} else {
		return Region{offset, offset}
	}
	if m := vwre2.FindStringIndex(line[begin:]); m != nil {
		end = begin + m[1]
	}
	return Region{lr.Begin() + begin, lr.Begin() + end}
}

func (v *View) Insert(edit *Edit, point int, value string) int {
	if t, ok := v.Settings().Get("translate_tabs_to_spaces", false).(bool); ok && t && strings.Contains(value, "\t") {
		tab_size, ok := v.Settings().Get("tab_size", 4).(int)
		if !ok {
			tab_size = 4
		}

		lines := strings.Split(value, "\n")
		for i, li := range lines {
			for {
				if idx := strings.Index(li, "\t"); idx != -1 {
					ai := idx
					if i == 0 {
						_, col := v.RowCol(point)
						ai = col
					}
					add := 1 + ((ai + (tab_size - 1)) &^ (tab_size - 1))
					spaces := ""
					for j := ai; j < add; j++ {
						spaces += " "
					}
					li = li[:idx] + spaces + li[idx+1:]
					continue
				}
				break
			}
			lines[i] = li
		}
		value = strings.Join(lines, "\n")
	}
	edit.composite.AddExec(NewInsertAction(v.buffer, point, value))
	return len(value)
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

func (v *View) findScope(search parser.Range, node *parser.Node) *parser.Node {
	idx := sort.Search(len(node.Children), func(i int) bool {
		return node.Children[i].Range.Start >= search.Start || node.Children[i].Range.Contains(search)
	})
	for idx < len(node.Children) {
		c := node.Children[idx]
		if c.Range.Start > search.End {
			break
		}
		if c.Range.Contains(search) {
			if node.Name != " " && node != v.lastScopeNode {
				v.lastScopeBuf.WriteByte(' ')
				v.lastScopeBuf.WriteString(node.Name)
			}
			return v.findScope(search, node.Children[idx])
		}
		idx++
	}
	if node != v.lastScopeNode && node.Range.Contains(search) && node.Name != "" {
		v.lastScopeBuf.WriteByte(' ')
		v.lastScopeBuf.WriteString(node.Name)
		return node
	}
	return nil
}

func (v *View) ScopeName(point int) string {
	if v.syntax.Language == nil {
		return ""
	}

	search := parser.Range{point, point + 1}
	if v.lastScopeNode != nil && v.lastScopeNode.Range.Contains(search) {
		if len(v.lastScopeNode.Children) != 0 {
			if no := v.findScope(search, v.lastScopeNode); no != v.lastScopeNode {
				v.lastScopeNode = no
				v.lastScopeName = v.lastScopeBuf.String()
			}
		}
	} else {
		v.lastScopeNode = nil
		v.lastScopeBuf.Reset()
		v.lastScopeBuf.WriteString(v.syntax.Language.ScopeName)
		v.lastScopeNode = v.findScope(search, v.syntax.RootNode())
		v.lastScopeName = v.lastScopeBuf.String()
	}
	return v.lastScopeName
}
