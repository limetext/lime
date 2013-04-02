package backend

import (
	. "lime/backend/primitives"
	"testing"
)

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
		{0, 1},
		{1, 2},
		{2, 3},
		{3, 4},
	}
	for _, r2 := range r {
		v.selection.Add(r2)
	}

	edit = v.BeginEdit()
	for _, ins := range "4321" {
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

	edit = v.BeginEdit()
	for _, ins := range []string{"world", "hello "} {
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
