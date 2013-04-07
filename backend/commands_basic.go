package backend

import (
	"code.google.com/p/log4go"
	"fmt"
	"lime/backend/primitives"
	"regexp"
)

type (
	InsertCommand struct {
		DefaultCommand
	}

	LeftDeleteCommand struct {
		DefaultCommand
	}

	RightDeleteCommand struct {
		DefaultCommand
	}

	MoveCommand struct {
		DefaultCommand
	}

	MoveToCommand struct {
		DefaultCommand
	}

	UndoCommand struct {
		DefaultCommand
		bypassUndoCommand
		hard bool
	}
	RedoCommand struct {
		DefaultCommand
		hard bool
		bypassUndoCommand
	}
	ScrollLinesCommand struct {
		DefaultCommand
		bypassUndoCommand
	}
	MarkUndoGroupsForGluingCommand struct {
		DefaultCommand
		bypassUndoCommand
	}

	GlueMarkedUndoGroupsCommand struct {
		DefaultCommand
		bypassUndoCommand
	}

	MaybeMarkUndoGroupsForGluingCommand struct {
		DefaultCommand
		bypassUndoCommand
	}

	UnmarkUndoGroupsForGluingCommand struct {
		DefaultCommand
		bypassUndoCommand
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
				_, col := v.Buffer().RowCol(r.A)
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

func (c *RightDeleteCommand) Run(v *View, e *Edit, args Args) error {
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
			r.B++
		}
		v.Erase(e, r)
		if sel.Len() != l {
			continue
		}
		i++
	}
	return nil
}

func move_action(v *View, extend bool, transform func(r primitives.Region) int) {
	sel := v.Sel()
	r := sel.Regions()
	for i := range r {
		r[i].B = transform(r[i])
		if !extend {
			r[i].A = r[i].B
		}
	}
	sel.Clear()
	for i := range r {
		sel.Add(r[i])
	}
}

func (c *MoveToCommand) Run(v *View, e *Edit, args Args) error {
	to, _ := args["to"].(string)
	extend, _ := args["extend"].(bool)

	switch to {
	case "eol":
		move_action(v, extend, func(r primitives.Region) int {
			line := v.Buffer().Line(r.B)
			return line.B
		})
	case "bol":
		move_action(v, extend, func(r primitives.Region) int {
			line := v.Buffer().Line(r.B)
			return line.A
		})
	case "eof":
		move_action(v, extend, func(r primitives.Region) int {
			return v.buffer.Size()
		})
	default:
		return fmt.Errorf("move_to: Unimplemented 'to' action: %s", to)
	}
	return nil
}

var _move_stops_re = regexp.MustCompile(`\b`)

func (c *MoveCommand) Run(v *View, e *Edit, args Args) error {
	by, ok := args["by"].(string)
	if !ok {
		return fmt.Errorf("move: Missing or invalid 'by' argument: %v", args)
	}
	extend, ok := args["extend"].(bool)
	fwd, ok := args["forward"].(bool)
	word_begin, ok := args["word_begin"].(bool)
	word_end, ok := args["word_end"].(bool)

	switch by {
	case "characters":
		dir := 1
		if !fwd {
			dir = -1
		}
		move_action(v, extend, func(r primitives.Region) int {
			return r.B + dir
		})
	case "stops":
		move_action(v, extend, func(r primitives.Region) int {
			var next primitives.Region
			word := v.Buffer().Word(r.B)
			if word_end && fwd && r.B < word.End() {
				next = word
			} else if word_begin && !fwd && r.B > word.Begin() {
				next = word
			} else if fwd {
				next = v.Buffer().Word(word.B + 1)
			} else {
				next = v.Buffer().Word(word.A - 1)
				next = v.Buffer().Word(next.A - 1)
			}

			if word_begin {
				return next.A
			} else if word_end {
				return next.B
			}
			return r.B
		})
	default:
		return fmt.Errorf("move: Unimplemented 'by' action: %s", by)
	}
	return nil
}

func (c *UndoCommand) Run(v *View, e *Edit, args Args) error {
	v.undoStack.Undo(c.hard)
	return nil
}

func (c *RedoCommand) Run(v *View, e *Edit, args Args) error {
	v.undoStack.Redo(c.hard)
	return nil
}

func (c *MarkUndoGroupsForGluingCommand) Run(v *View, e *Edit, args Args) error {
	v.undoStack.mark = v.undoStack.position
	return nil
}

func (c *UnmarkUndoGroupsForGluingCommand) Run(v *View, e *Edit, args Args) error {
	v.undoStack.mark = -1
	return nil
}

func (c *GlueMarkedUndoGroupsCommand) Run(v *View, e *Edit, args Args) error {
	if l, p := v.undoStack.position-v.undoStack.mark, v.undoStack.mark; p != -1 && (l-p) > 1 {
		e.command = "sequence"
		e.bypassUndo = true
		type entry struct {
			name string
			args Args
		}
		entries := make([]entry, v.undoStack.position-v.undoStack.mark)
		for i := range entries {
			a := v.undoStack.actions[i+v.undoStack.mark]
			entries[i].name = a.command
			entries[i].args = a.args
			e.composite.Add(a)
		}
		v.undoStack.position = v.undoStack.mark
		v.undoStack.actions = v.undoStack.actions[:v.undoStack.mark+1]
		e.args = make(Args)
		e.args["commands"] = entries
		v.undoStack.Add(e, true)
	}
	v.undoStack.mark = -1
	return nil
}

func (c *MaybeMarkUndoGroupsForGluingCommand) Run(v *View, e *Edit, args Args) error {
	if v.undoStack.mark == -1 {
		v.undoStack.mark = v.undoStack.position
	}
	return nil
}

func (c *ScrollLinesCommand) Run(v *View, e *Edit, args Args) error {
	amount, ok := args["amount"].(int)
	if !ok {
		return fmt.Errorf("scroll_lines: Missing or invalid 'amount' argument: %v", args)
	}
	fe := GetEditor().Frontend()
	r, _ := v.Buffer().RowCol(fe.VisibleRegion(v).Begin())
	r -= amount
	r = v.Buffer().TextPoint(r, 1)
	fe.Show(v, primitives.Region{r, r})
	return nil
}

func initBasicCommands() {
	log4go.Debug("Registering commands...")
	e := GetEditor()
	type Cmd struct {
		name string
		cmd  Command
	}
	cmds := []Cmd{
		{"insert", &InsertCommand{}},
		{"left_delete", &LeftDeleteCommand{}},
		{"right_delete", &RightDeleteCommand{}},
		{"move", &MoveCommand{}},
		{"move_to", &MoveToCommand{}},
		{"undo", &UndoCommand{hard: true}},
		{"redo", &RedoCommand{hard: true}},
		{"soft_undo", &UndoCommand{}},
		{"soft_redo", &RedoCommand{}},
		{"mark_undo_groups_for_gluing", &MarkUndoGroupsForGluingCommand{}},
		{"glue_marked_undo_groups", &GlueMarkedUndoGroupsCommand{}},
		{"maybe_mark_undo_groups_for_gluing", &MaybeMarkUndoGroupsForGluingCommand{}},
		{"unmark_undo_groups_for_gluing", &UnmarkUndoGroupsForGluingCommand{}},
		{"scroll_lines", &ScrollLinesCommand{}},
	}
	for i := range cmds {
		if err := e.CommandHandler().Register(cmds[i].name, cmds[i].cmd); err != nil {
			log4go.Error("Failed to register command %s: %s", cmds[i].name, err)
		}
	}
	log4go.Debug("Commands registered")
}
