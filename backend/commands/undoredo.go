// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/limetext/lime/backend"
)

type (
	UndoCommand struct {
		BypassUndoCommand
		hard bool
	}
	RedoCommand struct {
		BypassUndoCommand
		hard bool
	}
)

func (c *UndoCommand) Run(v *View, e *Edit) error {
	v.UndoStack().Undo(c.hard)
	return nil
}

func (c *RedoCommand) Run(v *View, e *Edit) error {
	v.UndoStack().Redo(c.hard)
	return nil
}

func init() {
	register([]Command{
		&UndoCommand{hard: true},
		&RedoCommand{hard: true},
	})

	registerByName([]namedCmd{
		{"soft_undo", &UndoCommand{}},
		{"soft_redo", &RedoCommand{}},
	})
}
