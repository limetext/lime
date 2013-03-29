package backend

type (
	View struct {
		HasSettings
		filename  string
		window    *Window
		buffer    *Buffer
		selection RegionSet
		undoStack undoStack
	}
	undoStack struct {
		position int
		actions  []Action
	}
)

func (v *View) setBuffer(b *Buffer) {
	v.buffer = b
	b.callbacks = append(b.callbacks, v.selection.adjust)
}

func (v *View) Window() *Window {
	return v.window
}

func (v *View) Buffer() *Buffer {
	return v.buffer
}

func (v *View) Insert(point int, value string) {
	ia := NewInsertAction(v.buffer, point, value)
	v.undoStack.Apply(ia)
}

func (us *undoStack) Apply(a Action) {
	if us.position != len(us.actions) {
		us.actions = us.actions[0:us.position]
	}
	us.actions = append(us.actions, a)
	us.Redo()
}

func (us *undoStack) Undo() {
	if us.position > 0 {
		us.position--
		us.actions[us.position].Undo()
	}
}

func (us *undoStack) Redo() {
	if us.position < len(us.actions) {
		us.actions[us.position].Apply()
		us.position++
	}
}
