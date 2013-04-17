package commands

import (
	"fmt"
	. "lime/backend"
	"lime/backend/primitives"
	"lime/backend/util"
)

type (
	MoveCommand struct {
		DefaultCommand
	}

	MoveToCommand struct {
		DefaultCommand
	}

	ScrollLinesCommand struct {
		BypassUndoCommand
	}
)

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
	sel.AddAll(r)
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
	case "bof":
		move_action(v, extend, func(r primitives.Region) int {
			return 0
		})
	case "eof":
		move_action(v, extend, func(r primitives.Region) int {
			return v.Buffer().Size()
		})
	default:
		return fmt.Errorf("move_to: Unimplemented 'to' action: %s", to)
	}
	return nil
}

func (c *MoveCommand) Run(v *View, e *Edit, args Args) error {
	p := util.Prof.Enter("move.run.init")
	by, ok := args["by"].(string)
	if !ok {
		return fmt.Errorf("move: Missing or invalid 'by' argument: %v", args)
	}
	extend, ok := args["extend"].(bool)
	fwd, ok := args["forward"].(bool)
	word_begin, ok := args["word_begin"].(bool)
	word_end, ok := args["word_end"].(bool)

	p.Exit()
	p = util.Prof.Enter("move.run.action")
	defer p.Exit()

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
	case "lines":
		move_action(v, extend, func(in primitives.Region) int {
			r, c := v.Buffer().RowCol(in.B)
			_ = r
			if !fwd {
				r--
			} else {
				r++
			}
			return v.Buffer().TextPoint(r, c)
		})
	default:
		return fmt.Errorf("move: Unimplemented 'by' action: %s", by)
	}
	return nil
}

func (c *ScrollLinesCommand) Run(v *View, e *Edit, args Args) error {
	var amount int
	a, ok := args["amount"]
	switch t := a.(type) {
	case int:
		amount = t
	case int64:
		amount = int(t)
	case float64:
		amount = int(t)
	default:
		ok = false
	}
	if !ok {
		return fmt.Errorf("scroll_lines: Missing or invalid 'amount' argument: %v", args)
	}
	fe := GetEditor().Frontend()
	vr := fe.VisibleRegion(v)
	var r int
	if amount >= 0 {
		r, _ = v.Buffer().RowCol(vr.Begin())
		r -= amount
	} else {
		r, _ = v.Buffer().RowCol(vr.End() - 1)
		r -= amount
	}
	r = v.Buffer().TextPoint(r, 0)
	fe.Show(v, primitives.Region{r, r})
	return nil
}

func init() {
	register([]cmd{
		{"move", &MoveCommand{}},
		{"move_to", &MoveToCommand{}},
		{"scroll_lines", &ScrollLinesCommand{}},
	})
}
