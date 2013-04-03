package backend

import (
	"code.google.com/p/log4go"
	"fmt"
)

type (
	InsertCommand struct {
		DefaultCommand
	}

	LeftDeleteCommand struct {
		DefaultCommand
	}
	MoveCommand struct {
		DefaultCommand
	}
	UndoCommand struct {
		DefaultCommand
	}
	RedoCommand struct {
		DefaultCommand
	}
)

func (c *InsertCommand) Run(v *View, e *Edit, args Args) error {
	sel := v.Sel()
	chars, ok := args["characters"].(string)
	if !ok {
		return fmt.Errorf("insert: Missing or invalid characters argument: %v", args)
	}
	for i := 0; i < sel.Len(); i++ {
		r := sel.Get(i)
		v.Insert(e, r.B, chars)
	}
	return nil
}

func (c *LeftDeleteCommand) Run(v *View, e *Edit, args Args) error {
	sel := v.Sel()
	for i := 0; i < sel.Len(); i++ {
		r := sel.Get(i)
		r.A, r.B = r.Begin()-1, r.End()
		v.Erase(e, r)
	}
	return nil
}

func (c *MoveCommand) Run(v *View, e *Edit, args Args) error {
	by, ok := args["by"].(string)
	if !ok {
		return fmt.Errorf("move: Missing or invalid 'by' argument: %v", args)
	}
	extend, ok := args["extend"].(bool)
	fwd, ok := args["forward"].(bool)
	sel := v.Sel()
	r := sel.Regions()

	switch by {
	case "characters":
		dir := 1
		if !fwd {
			dir = -1
		}
		if extend {
			for i := range r {
				r[i].B = r[i].B + dir
			}
		} else {
			for i := range r {
				r[i].A += dir
				r[i].B += dir
			}
		}
		sel.Clear()
		for i := range r {
			sel.Add(r[i])
		}
	default:
		return fmt.Errorf("move: Unimplemented 'by' action: %s", by)
	}
	return nil
}

func (c *UndoCommand) Run(w *Window, args Args) error {
	act := w.ActiveView()
	if act == nil {
		return fmt.Errorf("undo: no active view")
	}
	act.undoStack.Undo()
	return nil
}

func (c *RedoCommand) Run(w *Window, args Args) error {
	act := w.ActiveView()
	if act == nil {
		return fmt.Errorf("redo: no active view")
	}
	act.undoStack.Redo()
	return nil
}

func init() {
	log4go.Debug("Registering commands...")
	e := GetEditor()
	type Cmd struct {
		name string
		cmd  Command
	}
	cmds := []Cmd{
		{"insert", &InsertCommand{}},
		{"left_delete", &LeftDeleteCommand{}},
		{"move", &MoveCommand{}},
		{"undo", &UndoCommand{}},
		{"redo", &RedoCommand{}},
	}
	for i := range cmds {
		if err := e.CommandHandler().Register(cmds[i].name, cmds[i].cmd); err != nil {
			log4go.Error("Failed to register command %s: %s", cmds[i].name, err)
		}
	}
	log4go.Debug("Commands registered")
}
