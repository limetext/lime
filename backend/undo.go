package backend

type (
	undoStack struct {
		position int
		actions  []*Edit
		mark     int
	}
)

func (us *undoStack) Add(a *Edit, inc bool) {
	if us.position != len(us.actions) {
		us.actions = us.actions[0:us.position]
	}
	us.actions = append(us.actions, a)
	if inc {
		us.position++
	}
}

func (us *undoStack) index(relative int, modifying_only bool) int {
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

func (us *undoStack) Undo(hard bool) {
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

func (us *undoStack) Redo(hard bool) {
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
