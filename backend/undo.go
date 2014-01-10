// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package backend

type (
	UndoStack struct {
		position int
		actions  []*Edit
	}
)

// Adds the provided Edit object to the UndoStack, potentially
// destroying the old redo stack if one had been created.
// TODO(.): It would be nice with branched undo histories
func (us *UndoStack) Add(a *Edit) {
	if us.position != len(us.actions) {
		us.actions = us.actions[0:us.position]
	}
	us.actions = append(us.actions, a)
	us.position++
}

// index returns the real index in the UndoStack of an undo item
// relative to the current position.
//
// When modifying_only is set to true, only actions actually modifying
// the buffer (as opposed to just moving the cursor) are counted as an
// index. Also see comment in Undo.
func (us *UndoStack) index(relative int, modifying_only bool) int {
	dir := -1
	i := us.position
	if relative > 0 {
		dir = 1
	} else {
		i--
	}
	relative *= dir
	for ; i >= 0 && i < len(us.actions) && relative > 0; i += dir {
		if modifying_only {
			if us.actions[i].composite.Len() != 0 {
				relative--
			}
		} else {
			relative--
		}
	}
	if i >= 0 && i < len(us.actions) {
		return i
	} else {
		return -1
	}
}

// Reverts the last action on the UndoStack.
//
// When the argument "hard" is set to true,
// the "last action" will be the last action that
// modified the contents of the buffer (rather than just
// changing the cursor position). In this case, all
// actions between the current action and the last "hard"
// action will be reverted.
func (us *UndoStack) Undo(hard bool) {
	if us.position <= 0 {
		// Nothing to undo
		return
	}
	to := us.index(0, hard)
	if to == -1 {
		to = 0
	}
	for us.position > to {
		us.position--
		us.actions[us.position].Undo()
	}
}

// Re-applies the next action in the undo stack
// if there are any actions on the stack that had
// been undone.
//
// See comment in Undo regarding the use of "hard".
func (us *UndoStack) Redo(hard bool) {
	if us.position >= len(us.actions) {
		// No more actions to redo
		return
	}
	to := us.index(1, hard)
	if to == -1 {
		to = len(us.actions)
	}
	for us.position < to {
		us.actions[us.position].Apply()
		us.position++
	}
}

// Returns the current position in the UndoStack.
func (us *UndoStack) Position() int {
	return us.position
}

// Glues all edits from the position given by mark,
// to the current position in the UndoStack, replacing
// them by a single entry which will now be composite
// of all those other actions.
//
// In other words, after the glue operation
// a single "undo" operation will then undo all of those edits
// and a single redo after that will redo them all again.
func (us *UndoStack) GlueFrom(mark int) {
	if mark >= us.position {
		return
	}
	var e Edit
	e.command = "sequence"
	type entry struct {
		name string
		args Args
	}
	e.v = us.actions[mark].v
	e.savedSel.AddAll(us.actions[mark].savedSel.Regions())

	entries := make([]entry, us.position-mark)
	for i := range entries {
		a := us.actions[i+mark]
		entries[i].name = a.command
		entries[i].args = a.args
		e.composite.Add(a)
	}
	us.position = mark
	us.actions = us.actions[:mark+1]
	e.args = make(Args)
	e.args["commands"] = entries
	us.Add(&e)
}
