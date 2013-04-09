package backend

import (
	"fmt"
	"github.com/quarnster/completion/util"
	"io/ioutil"
	. "lime/backend/primitives"
	"math/rand"
	"reflect"
	"testing"
)

func TestView(t *testing.T) {
	var (
		w Window
		v = w.NewFile()
	)
	edit := v.BeginEdit()
	v.Insert(edit, 0, "abcd")
	v.EndEdit(edit)
	v.selection.Clear()
	r := []Region{
		{0, 0},
		{1, 1},
		{2, 2},
		{3, 3},
	}
	for _, r2 := range r {
		v.selection.Add(r2)
	}

	edit = v.BeginEdit()
	for _, ins := range "1234" {
		for i := 0; i < v.selection.Len(); i++ {
			v.Insert(edit, v.selection.Get(i).Begin(), string(ins))
		}
	}
	v.EndEdit(edit)

	if v.buffer.String() != "1234a1234b1234c1234d" {
		t.Error(v.buffer.String())
	}
	v.undoStack.Undo(true)
	if v.buffer.String() != "abcd" {
		t.Error("expected 'abcd', but got: ", v.buffer.String())
	}
	v.undoStack.Redo(true)
	if v.buffer.String() != "1234a1234b1234c1234d" {
		t.Error("expected '1234a1234b1234c1234d', but got: ", v.buffer.String())
	}

	v.selection.Clear()
	r = []Region{
		{0, 0},
		{5, 5},
		{10, 10},
		{15, 15},
	}
	for _, r2 := range r {
		v.selection.Add(r2)
	}

	edit = v.BeginEdit()
	for _, ins := range []string{"hello ", "world"} {
		for i := 0; i < v.selection.Len(); i++ {
			v.Insert(edit, v.selection.Get(i).Begin(), ins)
		}
	}
	v.EndEdit(edit)

	if v.buffer.String() != "hello world1234ahello world1234bhello world1234chello world1234d" {
		t.Error(v.buffer.String())
	}
	v.undoStack.Undo(true)

	if v.buffer.String() != "1234a1234b1234c1234d" {
		t.Error("expected '1234a1234b1234c1234d', but got: ", v.buffer.String())
	}
	v.undoStack.Undo(true)
	if v.buffer.String() != "abcd" {
		t.Error("expected 'abcd', but got: ", v.buffer.String())
	}
	v.undoStack.Undo(true)
	if v.buffer.String() != "" {
		t.Error("expected '', but got: ", v.buffer.String())
	}
	v.undoStack.Redo(true)
	if v.buffer.String() != "abcd" {
		t.Error("expected 'abcd', but got: ", v.buffer.String())
	}

	v.undoStack.Redo(true)
	if v.buffer.String() != "1234a1234b1234c1234d" {
		t.Error("expected '1234a1234b1234c1234d', but got: ", v.buffer.String())
	}

	v.undoStack.Redo(true)
	if v.buffer.String() != "hello world1234ahello world1234bhello world1234chello world1234d" {
		t.Error(v.buffer.String())
	}
}

func TestUndoRedoCommands(t *testing.T) {
	ch := GetEditor().CommandHandler()
	var (
		w Window
		v = w.NewFile()
	)
	edit := v.BeginEdit()
	v.Insert(edit, 0, "abcd")
	v.EndEdit(edit)
	v.selection.Clear()
	r := []Region{
		{0, 0},
		{1, 1},
		{2, 2},
		{3, 3},
	}
	for _, r2 := range r {
		v.selection.Add(r2)
	}

	edit = v.BeginEdit()
	for _, ins := range "1234" {
		for i := 0; i < v.selection.Len(); i++ {
			v.Insert(edit, v.selection.Get(i).Begin(), string(ins))
		}
	}
	v.EndEdit(edit)

	if v.buffer.String() != "1234a1234b1234c1234d" {
		t.Error(v.buffer.String())
	}
	ch.RunTextCommand(v, "undo", nil)
	if v.buffer.String() != "abcd" {
		t.Error("expected 'abcd', but got: ", v.buffer.String())
	}
	ch.RunTextCommand(v, "redo", nil)
	if v.buffer.String() != "1234a1234b1234c1234d" {
		t.Error("expected '1234a1234b1234c1234d', but got: ", v.buffer.String())
	}

	v.selection.Clear()
	r = []Region{
		{0, 0},
		{5, 5},
		{10, 10},
		{15, 15},
	}
	for _, r2 := range r {
		v.selection.Add(r2)
	}

	edit = v.BeginEdit()
	for _, ins := range []string{"hello ", "world"} {
		for i := 0; i < v.selection.Len(); i++ {
			v.Insert(edit, v.selection.Get(i).Begin(), ins)
		}
	}
	v.EndEdit(edit)

	if v.buffer.String() != "hello world1234ahello world1234bhello world1234chello world1234d" {
		t.Error(v.buffer.String())
	}
	ch.RunTextCommand(v, "undo", nil)

	if v.buffer.String() != "1234a1234b1234c1234d" {
		t.Error("expected '1234a1234b1234c1234d', but got: ", v.buffer.String())
	}
	ch.RunTextCommand(v, "undo", nil)
	if v.buffer.String() != "abcd" {
		t.Error("expected 'abcd', but got: ", v.buffer.String())
	}
	ch.RunTextCommand(v, "undo", nil)
	if v.buffer.String() != "" {
		t.Error("expected '', but got: ", v.buffer.String())
	}
	v.undoStack.Redo(true)
	if v.buffer.String() != "abcd" {
		t.Error("expected 'abcd', but got: ", v.buffer.String())
	}

	v.undoStack.Redo(true)
	if v.buffer.String() != "1234a1234b1234c1234d" {
		t.Error("expected '1234a1234b1234c1234d', but got: ", v.buffer.String())
	}

	v.undoStack.Redo(true)
	if v.buffer.String() != "hello world1234ahello world1234bhello world1234chello world1234d" {
		t.Error(v.buffer.String())
	}
}

func TestErase(t *testing.T) {
	var (
		w Window
		v = w.NewFile()
		s = v.Sel()
	)
	edit := v.BeginEdit()
	v.Insert(edit, 0, "1234abcd5678abcd")
	v.EndEdit(edit)
	s.Clear()
	v.Sel().Add(Region{4, 8})
	v.Sel().Add(Region{12, 16})

	edit = v.BeginEdit()
	for i := 0; i < s.Len(); i++ {
		v.Erase(edit, s.Get(i))
	}
	v.EndEdit(edit)
	if !reflect.DeepEqual(s.Regions(), []Region{{4, 4}, {8, 8}}) {
		t.Error(s)
	}
	if d := v.buffer.String(); d != "12345678" {
		t.Error(d)
	}
}

func TestScopeName(t *testing.T) {
	var (
		w Window
		v = w.NewFile()
	)
	const (
		in      = "textmate/testdata/main.go"
		expfile = "testdata/scopename.res"
		syntax  = "textmate/testdata/Go.tmLanguage"
	)
	v.SetSyntaxFile(syntax)
	if d, err := ioutil.ReadFile(in); err != nil {
		t.Fatal(err)
	} else {
		e := v.BeginEdit()
		v.Insert(e, 0, string(d))
		v.EndEdit(e)
		last := ""
		str := ""
		lasti := 0
		for i := 0; i < v.buffer.Size(); i++ {
			if name := v.ScopeName(i); name != last {
				if last != "" {
					str += fmt.Sprintf("%d-%d: %s\n", lasti, i, last)
					lasti = i
				}
				last = name
			}
		}
		if i := v.Buffer().Size(); lasti != i {
			str += fmt.Sprintf("%d-%d: %s\n", lasti, i, last)
		}
		if d, err := ioutil.ReadFile(expfile); err != nil {
			if err := ioutil.WriteFile(expfile, []byte(str), 0644); err != nil {
				t.Error(err)
			}
		} else if diff := util.Diff(string(d), str); diff != "" {
			t.Error(diff)
		}

	}
}

func BenchmarkScopeNameLinear(b *testing.B) {
	var (
		w Window
		v = w.NewFile()
	)
	const (
		in     = "textmate/language_test.go"
		syntax = "textmate/testdata/Go.tmLanguage"
	)
	b.StopTimer()
	v.SetSyntaxFile(syntax)
	if d, err := ioutil.ReadFile(in); err != nil {
		b.Fatal(err)
	} else {
		e := v.BeginEdit()
		v.Insert(e, 0, string(d))
		v.EndEdit(e)
		b.StartTimer()
		for j := 0; j < b.N; j++ {
			for i := 0; i < v.buffer.Size(); i++ {
				v.ScopeName(i)
			}
		}
	}
}

func BenchmarkScopeNameRandom(b *testing.B) {
	var (
		w Window
		v = w.NewFile()
	)
	const (
		in     = "textmate/language_test.go"
		syntax = "textmate/testdata/Go.tmLanguage"
	)
	b.StopTimer()
	v.SetSyntaxFile(syntax)
	if d, err := ioutil.ReadFile(in); err != nil {
		b.Fatal(err)
	} else {
		e := v.BeginEdit()
		v.Insert(e, 0, string(d))
		v.EndEdit(e)
		p := rand.Perm(b.N)
		b.StartTimer()
		for _, i := range p {
			v.ScopeName(i)
		}
	}
}
