// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	"fmt"
	. "github.com/limetext/lime/backend"
	"github.com/limetext/lime/backend/util"
	"github.com/limetext/text"
	"strings"
)

const (
	// Beginning of line
	BOL MoveToType = iota
	// End of line
	EOL
	// Beginning of file
	BOF
	// End of file
	EOF
	// Current level close bracket
	Brackets
)

const (
	// Move by Characters
	Characters MoveByType = iota
	// Move by Stops (TODO(.): what exactly is a stop?)
	Stops
	// Move by Lines
	Lines
	// Move by Words
	Words
	// Move by Word Ends
	WordEnds
	// Move by Sub Words
	SubWords
	// Move by Sub Word Ends
	SubWordEnds
	// Move by Page
	Pages
)

type (
	// The MoveCommand moves the current selection
	MoveCommand struct {
		DefaultCommand
		// Specifies the type of "move" operation
		By MoveByType
		// Whether the current selection should be extended or not
		Extend bool
		// Whether to move forward or backwards
		Forward bool
		// Used together with By=Stops, go to word begin
		WordBegin bool
		// Used together with By=Stops, go to word end
		WordEnd bool
		// Used together with By=Stops, go to punctuation begin
		PunctBegin bool
		// Used together with By=Stops, go to punctuation end
		PunctEnd bool
		// Used together with By=Stops, extends "word_separators" defined by settings
		Separators string
		// Used together with By=Stops, go to an empty line
		EmptyLine bool
		// Used together with By=Stops, TODO: ???
		ClipToLine bool
	}

	// Specifies the type of "move" operation
	MoveByType int

	// Specifies the type of "move_to" operation to perform
	MoveToType int

	// The MoveToCommand moves or extends the current selection to the specified location
	MoveToCommand struct {
		DefaultCommand
		// The type of "move_to" operation to perform
		To MoveToType
		// Whether the current selection should be extended or not
		Extend bool
	}

	// The ScrollLinesCommand moves the viewpoint "Amount" lines from the current location
	ScrollLinesCommand struct {
		BypassUndoCommand
		// The number of lines to scroll (positive or negative direction)
		Amount int
	}
)

func move_action(v *View, extend bool, transform func(r text.Region) int) {
	sel := v.Sel()
	r := sel.Regions()
	bs := v.Buffer().Size()
	for i := range r {
		r[i].B = transform(r[i])
		if r[i].B < 0 {
			r[i].B = 0
		} else if r[i].B > bs {
			// Yes > the size, and not size-1 because the cursor being at "size"
			// is the position it will be at when we are appending
			// to the buffer
			r[i].B = bs
		}

		if !extend {
			r[i].A = r[i].B
		}
	}
	sel.Clear()
	sel.AddAll(r)
}

func (mt *MoveToType) Set(v interface{}) error {
	switch to := v.(string); to {
	case "eol":
		*mt = EOL
	case "bol":
		*mt = BOL
	case "bof":
		*mt = BOF
	case "eof":
		*mt = EOF
	case "brackets":
		*mt = Brackets
	default:
		return fmt.Errorf("move_to: Unimplemented 'to' type: %s", to)
	}
	return nil
}

func (c *MoveToCommand) Run(v *View, e *Edit) error {
	switch c.To {
	case EOL:
		move_action(v, c.Extend, func(r text.Region) int {
			line := v.Buffer().Line(r.B)
			return line.B
		})
	case BOL:
		move_action(v, c.Extend, func(r text.Region) int {
			line := v.Buffer().Line(r.B)
			return line.A
		})
	case BOF:
		move_action(v, c.Extend, func(r text.Region) int {
			return 0
		})
	case EOF:
		move_action(v, c.Extend, func(r text.Region) int {
			return v.Buffer().Size()
		})
	case Brackets:
		move_action(v, c.Extend, func(r text.Region) (pos int) {
			var (
				of          int
				co          = 1
				str, br, rv string
				opening     = "([{"
				closing     = ")]}"
			)
			pos = r.B

			// next and before character
			n := v.Buffer().Substr(text.Region{r.B, r.B + 1})
			b := v.Buffer().Substr(text.Region{r.B, r.B - 1})
			if strings.ContainsAny(n, opening) {
				// TODO: Maybe it's better to use sth like view.FindByClass or even
				// view.FindByClass() function itself instead of getting whole text
				// and looping through it. With using view.FindByClass() function
				// backward we won't need to reverse the text anymore
				str = v.Buffer().Substr(text.Region{r.B + 1, v.Buffer().Size()})
				br = n
				rv = revert(n)
				of = 2
			} else if strings.ContainsAny(b, closing) {
				// TODO: same as above
				str = v.Buffer().Substr(text.Region{0, r.B - 1})
				br = b
				rv = revert(b)
				str = reverse(str)
				co = -1
				of = -2
			} else if strings.ContainsAny(n, closing) {
				// TODO: same as above
				str = v.Buffer().Substr(text.Region{0, r.B - 1})
				br = n
				rv = revert(n)
				str = reverse(str)
				co = -1
				of = -1
			} else {
				// TODO: same as above
				str = v.Buffer().Substr(text.Region{r.B, v.Buffer().Size()})
				bef := v.Buffer().Substr(text.Region{0, r.B})
				if p := strings.LastIndexAny(bef, opening); p == -1 {
					return
				} else {
					br = string(bef[p])
					rv = revert(br)
				}
			}
			count := 1
			for i, c := range str {
				if ch := string(c); ch == br {
					count++
				} else if ch == rv {
					count--
				}
				if count == 0 {
					return i*co + r.B + of
				}
			}
			return
		})
	default:
		return fmt.Errorf("move_to: Unimplemented 'to' action: %d", c.To)
	}
	return nil
}

func (m *MoveByType) Set(v interface{}) error {
	switch by := v.(string); by {
	case "lines":
		*m = Lines
	case "characters":
		*m = Characters
	case "stops":
		*m = Stops
	case "words":
		*m = Words
	case "word_ends":
		*m = WordEnds
	case "subwords":
		*m = SubWords
	case "subword_ends":
		*m = SubWordEnds
	case "pages":
		*m = Pages
	default:
		return fmt.Errorf("move: Unimplemented 'by' action: %s", by)
	}
	return nil
}

func (c *MoveCommand) Run(v *View, e *Edit) error {
	p := util.Prof.Enter("move.run.init")

	p.Exit()
	p = util.Prof.Enter("move.run.action")
	defer p.Exit()

	switch c.By {
	case Characters:
		dir := 1
		if !c.Forward {
			dir = -1
		}
		move_action(v, c.Extend, func(r text.Region) int {
			return r.B + dir
		})
	case Stops:
		move_action(v, c.Extend, func(r text.Region) int {
			var next text.Region
			word := v.Buffer().Word(r.B)
			if c.WordEnd && c.Forward && r.B < word.End() {
				next = word
			} else if c.WordBegin && !c.Forward && r.B > word.Begin() {
				next = word
			} else if c.Forward {
				next = v.Buffer().Word(word.B + 1)
			} else {
				next = v.Buffer().Word(word.A - 1)
				next = v.Buffer().Word(next.A - 1)
			}

			if c.WordBegin {
				return next.A
			} else if c.WordEnd {
				return next.B
			}
			return r.B
		})
	case Lines:
		move_action(v, c.Extend, func(in text.Region) int {
			r, col := v.Buffer().RowCol(in.B)
			_ = r
			if !c.Forward {
				r--
			} else {
				r++
			}
			return v.Buffer().TextPoint(r, col)
		})
	case Words:
		move_action(v, c.Extend, func(in text.Region) int {
			return v.FindByClass(in.B, c.Forward, CLASS_WORD_START|
				CLASS_LINE_END|CLASS_LINE_START)
		})
	case WordEnds:
		move_action(v, c.Extend, func(in text.Region) int {
			return v.FindByClass(in.B, c.Forward, CLASS_WORD_END|
				CLASS_LINE_END|CLASS_LINE_START)
		})
	case Pages:
		// TODO: Should know how many lines does the frontend show in one page
	}
	return nil
}

func revert(c string) string {
	switch c {
	case "(":
		return ")"
	case ")":
		return "("
	case "[":
		return "]"
	case "]":
		return "["
	case "{":
		return "}"
	case "}":
		return "{"
	}
	return ""
}

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func (c *ScrollLinesCommand) Run(v *View, e *Edit) error {
	fe := GetEditor().Frontend()
	vr := fe.VisibleRegion(v)
	var r int
	if c.Amount >= 0 {
		r, _ = v.Buffer().RowCol(vr.Begin())
		r -= c.Amount
	} else {
		r, _ = v.Buffer().RowCol(vr.End() - 1)
		r -= c.Amount
	}
	r = v.Buffer().TextPoint(r, 0)
	fe.Show(v, text.Region{A: r, B: r})
	return nil
}

func init() {
	register([]Command{
		&MoveCommand{},
		&MoveToCommand{},
		&ScrollLinesCommand{},
	})
}
