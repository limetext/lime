package backend

import (
	. "lime/backend/primitives"
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
	}
	Edit struct {
		CompositeAction
	}
)

func (v *View) setBuffer(b *Buffer) {
	v.buffer = b
	b.AddCallback(v.selection.Adjust)
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
