package commands

import (
	. "lime/backend"
	. "lime/backend/primitives"
	"testing"
)

func TestUndoRedoCommands(t *testing.T) {
	ch := GetEditor().CommandHandler()
	var (
		w Window
		v = w.NewFile()
	)
	edit := v.BeginEdit()
	v.Insert(edit, 0, "abcd")
	v.EndEdit(edit)
	v.Sel().Clear()
	r := []Region{
		{0, 0},
		{1, 1},
		{2, 2},
		{3, 3},
	}
	for _, r2 := range r {
		v.Sel().Add(r2)
	}

	edit = v.BeginEdit()
	for _, ins := range "1234" {
		for i := 0; i < v.Sel().Len(); i++ {
			v.Insert(edit, v.Sel().Get(i).Begin(), string(ins))
		}
	}
	v.EndEdit(edit)

	if v.Buffer().String() != "1234a1234b1234c1234d" {
		t.Error(v.Buffer().String())
	}
	ch.RunTextCommand(v, "undo", nil)
	if v.Buffer().String() != "abcd" {
		t.Error("expected 'abcd', but got: ", v.Buffer().String())
	}
	ch.RunTextCommand(v, "redo", nil)
	if v.Buffer().String() != "1234a1234b1234c1234d" {
		t.Error("expected '1234a1234b1234c1234d', but got: ", v.Buffer().String())
	}

	v.Sel().Clear()
	r = []Region{
		{0, 0},
		{5, 5},
		{10, 10},
		{15, 15},
	}
	for _, r2 := range r {
		v.Sel().Add(r2)
	}

	edit = v.BeginEdit()
	for _, ins := range []string{"hello ", "world"} {
		for i := 0; i < v.Sel().Len(); i++ {
			v.Insert(edit, v.Sel().Get(i).Begin(), ins)
		}
	}
	v.EndEdit(edit)

	if v.Buffer().String() != "hello world1234ahello world1234bhello world1234chello world1234d" {
		t.Error(v.Buffer().String())
	}
	ch.RunTextCommand(v, "undo", nil)

	if v.Buffer().String() != "1234a1234b1234c1234d" {
		t.Error("expected '1234a1234b1234c1234d', but got: ", v.Buffer().String())
	}
	ch.RunTextCommand(v, "undo", nil)
	if v.Buffer().String() != "abcd" {
		t.Error("expected 'abcd', but got: ", v.Buffer().String())
	}
	ch.RunTextCommand(v, "undo", nil)
	if v.Buffer().String() != "" {
		t.Error("expected '', but got: ", v.Buffer().String())
	}
	v.UndoStack().Redo(true)
	if v.Buffer().String() != "abcd" {
		t.Error("expected 'abcd', but got: ", v.Buffer().String())
	}

	v.UndoStack().Redo(true)
	if v.Buffer().String() != "1234a1234b1234c1234d" {
		t.Error("expected '1234a1234b1234c1234d', but got: ", v.Buffer().String())
	}

	v.UndoStack().Redo(true)
	if v.Buffer().String() != "hello world1234ahello world1234bhello world1234chello world1234d" {
		t.Error(v.Buffer().String())
	}
}
