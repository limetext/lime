package commands

import (
	"fmt"
	"lime/backend"
)

const lime_cmd_mark = "lime.cmd.mark"

type (
	MarkUndoGroupsForGluingCommand struct {
		backend.DefaultCommand
		backend.BypassUndoCommand
	}

	GlueMarkedUndoGroupsCommand struct {
		backend.DefaultCommand
		backend.BypassUndoCommand
	}

	MaybeMarkUndoGroupsForGluingCommand struct {
		backend.DefaultCommand
		backend.BypassUndoCommand
	}

	UnmarkUndoGroupsForGluingCommand struct {
		backend.DefaultCommand
		backend.BypassUndoCommand
	}
)

func (c *MarkUndoGroupsForGluingCommand) Run(v *backend.View, e *backend.Edit, args backend.Args) error {
	v.Settings().Set(lime_cmd_mark, v.UndoStack().Position())
	return nil
}

func (c *UnmarkUndoGroupsForGluingCommand) Run(v *backend.View, e *backend.Edit, args backend.Args) error {
	v.Settings().Erase(lime_cmd_mark)
	return nil
}

func (c *GlueMarkedUndoGroupsCommand) Run(v *backend.View, e *backend.Edit, args backend.Args) error {
	pos := v.UndoStack().Position()
	if mark, ok := v.Settings().Get(lime_cmd_mark).(int); ok {
		if l, p := pos-mark, mark; p != -1 && (l-p) > 1 {
			v.UndoStack().GlueFrom(mark)
		}
	} else {
		return fmt.Errorf("No mark in the current view")
	}
	return nil
}

func (c *MaybeMarkUndoGroupsForGluingCommand) Run(v *backend.View, e *backend.Edit, args backend.Args) error {
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
