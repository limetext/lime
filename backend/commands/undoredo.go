// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.
package commands

import (
	"lime/backend"
)

type (
	UndoCommand struct {
		backend.BypassUndoCommand
		hard bool
	}
	RedoCommand struct {
		backend.BypassUndoCommand
		hard bool
	}
)

func (c *UndoCommand) Run(v *backend.View, e *backend.Edit) error {
	v.UndoStack().Undo(c.hard)
	return nil
}

func (c *RedoCommand) Run(v *backend.View, e *backend.Edit) error {
	v.UndoStack().Redo(c.hard)
	return nil
}

func init() {
	register([]cmd{
		{"undo", &UndoCommand{hard: true}},
		{"redo", &RedoCommand{hard: true}},
		{"soft_undo", &UndoCommand{}},
		{"soft_redo", &RedoCommand{}},
	})
}
