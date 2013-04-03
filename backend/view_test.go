package backend

import (
	"encoding/json"
	"fmt"
	"github.com/quarnster/completion/util"
	"io/ioutil"
	. "lime/backend/primitives"
	"testing"
)

func TestRowCol(t *testing.T) {
	var (
		w Window
		v = w.NewView()
	)

	if d, err := ioutil.ReadFile("./testdata/unittest.json"); err != nil {
		t.Fatal(err)
	} else {
		type Test struct {
			Offset          int
			Line            int
			Column          int
			LineUntilOffset string
			LineAtOffset    string
			WordAtOffset    string
		}
		var tests []Test
		if err := json.Unmarshal(d, &tests); err != nil {
			t.Fatal(err)
		} else {
			if d, err := ioutil.ReadFile("./testdata/unittest.cpp"); err != nil {
				t.Fatal(err)
			} else {
				e := v.BeginEdit()
				v.Insert(e, 0, string(d))
				v.EndEdit(e)
			}

			for _, test := range tests {
				var a Test
				a.Line, a.Column = v.RowCol(test.Offset)
				a.LineAtOffset = v.Substr(v.Line(test.Offset))
				a.WordAtOffset = v.Substr(v.Word(test.Offset))
				a.Offset = v.TextPoint(test.Line, test.Column)
				if a.Line != test.Line {
					t.Errorf("Line mismatch: %d != %d", a.Line, test.Line)
				}
				if a.Column != test.Column {
					t.Errorf("Column mismatch: %d != %d", a.Column, test.Column)
				}
				if a.Offset != test.Offset {
					t.Errorf("Offset mismatch: %d != %d", a.Offset, test.Offset)
				}
				if a.LineAtOffset != test.LineAtOffset {
					t.Errorf("LineAtOffset mismatch: '%s' != '%s'", a.LineAtOffset, test.LineAtOffset)
				}
				if a.WordAtOffset != test.WordAtOffset {
					t.Errorf("WordAtOffset mismatch: '%s' != '%s'", a.WordAtOffset, test.WordAtOffset)
				}
			}
		}
	}
}

func TestView(t *testing.T) {
	var (
		w Window
		v = w.NewView()
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

	if v.buffer.Data() != "1234a1234b1234c1234d" {
		t.Error(v.buffer.Data())
	}
	v.undoStack.Undo()
	if v.buffer.Data() != "abcd" {
		t.Error("expected 'abcd', but got: ", v.buffer.Data())
	}
	v.undoStack.Redo()
	if v.buffer.Data() != "1234a1234b1234c1234d" {
		t.Error("expected '1234a1234b1234c1234d', but got: ", v.buffer.Data())
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

	if v.buffer.Data() != "hello world1234ahello world1234bhello world1234chello world1234d" {
		t.Error(v.buffer.Data())
	}
	v.undoStack.Undo()

	if v.buffer.Data() != "1234a1234b1234c1234d" {
		t.Error("expected '1234a1234b1234c1234d', but got: ", v.buffer.Data())
	}
	v.undoStack.Undo()
	if v.buffer.Data() != "abcd" {
		t.Error("expected 'abcd', but got: ", v.buffer.Data())
	}
	v.undoStack.Undo()
	if v.buffer.Data() != "" {
		t.Error("expected '', but got: ", v.buffer.Data())
	}
	v.undoStack.Redo()
	if v.buffer.Data() != "abcd" {
		t.Error("expected 'abcd', but got: ", v.buffer.Data())
	}

	v.undoStack.Redo()
	if v.buffer.Data() != "1234a1234b1234c1234d" {
		t.Error("expected '1234a1234b1234c1234d', but got: ", v.buffer.Data())
	}

	v.undoStack.Redo()
	if v.buffer.Data() != "hello world1234ahello world1234bhello world1234chello world1234d" {
		t.Error(v.buffer.Data())
	}
}

func TestScopeName(t *testing.T) {
	var (
		w Window
		v = w.NewView()
	)
	const (
		in      = "textmate/language_test.go"
		expfile = "./testdata/scopename.res"
		syntax  = "../3rdparty/bundles/go.tmbundle/Syntaxes/Go.tmLanguage"
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
		for i := 0; i < v.Size(); i++ {
			if name := v.ScopeName(i); name != last {
				last = name
				str += fmt.Sprintf("%d: %s\n", i, name)
			}
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
