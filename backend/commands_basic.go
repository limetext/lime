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
		if r.Size() == 0 {
			v.Insert(e, r.B, chars)
		} else {
			v.Replace(e, r, chars)
		}
	}
	return nil
}

func (c *LeftDeleteCommand) Run(v *View, e *Edit, args Args) error {
	trim_space := false
	tab_size := 4
	if t, ok := v.Settings().Get("translate_tabs_to_spaces", false).(bool); ok && t {
		if t, ok := v.Settings().Get("use_tab_stops", true).(bool); ok && t {
			trim_space = true
			tab_size, ok = v.Settings().Get("tab_size", 4).(int)
			if !ok {
				tab_size = 4
			}
		}
	}

	sel := v.Sel()
	hasNonEmpty := false
	for _, r := range sel.Regions() {
		if !r.Empty() {
			hasNonEmpty = true
			break
		}
	}
	i := 0
	for {
		l := sel.Len()
		if i >= l {
			break
		}
		r := sel.Get(i)
		if r.A == r.B && !hasNonEmpty {
			d := v.buffer.Data()
			if trim_space {
				_, col := v.RowCol(r.A)
				col -= 1
				prev_col := r.A - (col - (col-tab_size+(tab_size-1))&^(tab_size-1))
				if prev_col < 0 {
					prev_col = 0
				}
				for r.A > prev_col && d[r.A-1] == ' ' {
					r.A--
				}
			}
			if r.A == r.B {
				r.A--
			}
		}
		v.Erase(e, r)
		if sel.Len() != l {
			continue
		}
		i++
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
				r[i].B += dir
				r[i].A = r[i].B
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
