// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	"fmt"
	"lime/backend"
)

const lime_cmd_mark = "lime.cmd.mark"

type (
	// The MarkUndoGroupsForGluingCommand marks the current position
	// in the undo stack as the start of commands to glue, potentially
	// overwriting any existing marks.
	MarkUndoGroupsForGluingCommand struct {
		backend.BypassUndoCommand
	}

	// The GlueMarkedUndoGroupsCommand merges commands from the previously
	// marked undo stack location to the current location into a single
	// entry in the undo stack.
	GlueMarkedUndoGroupsCommand struct {
		backend.BypassUndoCommand
	}

	// The MaybeMarkUndoGroupsForGluingCommand is similar to
	// MarkUndoGroupsForGluingCommand with the exception that if there
	// is already a mark set, it is not overwritten.
	MaybeMarkUndoGroupsForGluingCommand struct {
		backend.BypassUndoCommand
	}

	// The UnmarkUndoGroupsForGluingCommand removes the glue mark set by
	// either MarkUndoGroupsForGluingCommand or MaybeMarkUndoGroupsForGluingCommand
	// if it was set.
	UnmarkUndoGroupsForGluingCommand struct {
		backend.BypassUndoCommand
	}
)

func (c *MarkUndoGroupsForGluingCommand) Run(v *backend.View, e *backend.Edit) error {
	v.Settings().Set(lime_cmd_mark, v.UndoStack().Position())
	return nil
}

func (c *UnmarkUndoGroupsForGluingCommand) Run(v *backend.View, e *backend.Edit) error {
	v.Settings().Erase(lime_cmd_mark)
	return nil
}

func (c *GlueMarkedUndoGroupsCommand) Run(v *backend.View, e *backend.Edit) error {
	pos := v.UndoStack().Position()
	mark, ok := v.Settings().Get(lime_cmd_mark).(int)
	if !ok {
		return fmt.Errorf("No mark in the current view")
	}
	if l, p := pos-mark, mark; p != -1 && (l-p) > 1 {
		v.UndoStack().GlueFrom(mark)
	}
	return nil
}

func (c *MaybeMarkUndoGroupsForGluingCommand) Run(v *backend.View, e *backend.Edit) error {
	if !v.Settings().Has(lime_cmd_mark) {
		v.Settings().Set(lime_cmd_mark, v.UndoStack().Position())
	}
	return nil
}

func init() {
	register([]cmd{
		{"mark_undo_groups_for_gluing", &MarkUndoGroupsForGluingCommand{}},
		{"glue_marked_undo_groups", &GlueMarkedUndoGroupsCommand{}},
		{"maybe_mark_undo_groups_for_gluing", &MaybeMarkUndoGroupsForGluingCommand{}},
		{"unmark_undo_groups_for_gluing", &UnmarkUndoGroupsForGluingCommand{}},
	})
}
