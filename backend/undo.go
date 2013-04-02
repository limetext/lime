package backend

import (
	. "lime/backend/primitives"
)

type (
	undoStack struct {
		position int
		actions  []Action
	}
)

func (us *undoStack) Add(a Action, inc bool) {
	if us.position != len(us.actions) {
		us.actions = us.actions[0:us.position]
	}
	us.actions = append(us.actions, a)
	if inc {
		us.position++
	}
}

func (us *undoStack) ApplyAction(a Action) {
	us.Add(a, false)
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
