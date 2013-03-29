package backend

import (
	"testing"
)

func TestView(t *testing.T) {
	var (
		w Window
		v = w.NewView()
	)
	v.Insert(0, "abcd")
	v.selection = RegionSet{
		[]Region{
			{0, 1},
			{1, 2},
			{2, 3},
			{3, 4},
		},
	}
	for _, ins := range "4321" {
		for i := range v.selection.Regions {
			v.Insert(v.selection.Regions[i].Begin(), string(ins))
		}
	}

	if v.buffer.data != "1234a1234b1234c1234d" {
		t.Error(v.buffer.data)
	}
	for i := v.undoStack.position - 1; i > 0; i-- {
		v.undoStack.Undo()
	}
	if v.buffer.data != "abcd" {
		t.Error(v.buffer.data)
	}
	for i := v.undoStack.position; i < len(v.undoStack.actions)-8; i++ {
		v.undoStack.Redo()
	}
	if v.buffer.data != "34a34b34c34d" {
		t.Error(v.buffer.data)
	}
	for _, ins := range []string{"world", "hello "} {
		for i := range v.selection.Regions {
			v.Insert(v.selection.Regions[i].Begin(), ins)
		}
	}

	if v.buffer.data != "hello world34ahello world34bhello world34chello world34d" {
		t.Error(v.buffer.data)
	}
}
