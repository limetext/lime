package backend

import (
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
	v.selection = RegionSet{
		[]Region{
			{0, 1},
			{1, 2},
			{2, 3},
			{3, 4},
		},
	}

	edit = v.BeginEdit()
	for _, ins := range "4321" {
		for i := range v.selection.regions {
			v.Insert(edit, v.selection.regions[i].Begin(), string(ins))
		}
	}
	v.EndEdit(edit)

	if v.buffer.data != "1234a1234b1234c1234d" {
		t.Error(v.buffer.data)
	}
	v.undoStack.Undo()
	if v.buffer.data != "abcd" {
		t.Error("expected 'abcd', but got: ", v.buffer.data)
	}
	v.undoStack.Redo()
	if v.buffer.data != "1234a1234b1234c1234d" {
		t.Error("expected '1234a1234b1234c1234d', but got: ", v.buffer.data)
	}

	edit = v.BeginEdit()
	for _, ins := range []string{"world", "hello "} {
		for i := range v.selection.regions {
			v.Insert(edit, v.selection.regions[i].Begin(), ins)
		}
	}
	v.EndEdit(edit)

	if v.buffer.data != "hello world1234ahello world1234bhello world1234chello world1234d" {
		t.Error(v.buffer.data)
	}
	v.undoStack.Undo()

	if v.buffer.data != "1234a1234b1234c1234d" {
		t.Error("expected '1234a1234b1234c1234d', but got: ", v.buffer.data)
	}
	v.undoStack.Undo()
	if v.buffer.data != "abcd" {
		t.Error("expected 'abcd', but got: ", v.buffer.data)
	}
	v.undoStack.Undo()
	if v.buffer.data != "" {
		t.Error("expected '', but got: ", v.buffer.data)
	}
	v.undoStack.Redo()
	if v.buffer.data != "abcd" {
		t.Error("expected 'abcd', but got: ", v.buffer.data)
	}

	v.undoStack.Redo()
	if v.buffer.data != "1234a1234b1234c1234d" {
		t.Error("expected '1234a1234b1234c1234d', but got: ", v.buffer.data)
	}

	v.undoStack.Redo()
	if v.buffer.data != "hello world1234ahello world1234bhello world1234chello world1234d" {
		t.Error(v.buffer.data)
	}
}
