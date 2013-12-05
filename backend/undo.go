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

func (us *UndoStack) Add(a *Edit) {
	if us.position != len(us.actions) {
		us.actions = us.actions[0:us.position]
	}
	us.actions = append(us.actions, a)
	us.position++
}

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

func (us *UndoStack) Undo(hard bool) {
	if us.position > 0 {
		to := us.index(0, hard)
		if to == -1 {
			to = 0
		}
		for us.position > to {
			us.position--
			us.actions[us.position].Undo()
		}
	}
}

func (us *UndoStack) Redo(hard bool) {
	if us.position < len(us.actions) {
		to := us.index(1, hard)
		if to == -1 {
			to = len(us.actions)
		}
		for us.position < to {
			us.actions[us.position].Apply()
			us.position++
		}
	}
}

func (us *UndoStack) Position() int {
	return us.position
}

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
