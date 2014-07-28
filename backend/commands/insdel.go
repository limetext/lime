// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/limetext/lime/backend"
	"github.com/quarnster/util/text"
	"strings"
)

type (
	// The InsertCommand inserts the given characters, at all
	// of the current selection locations, possibly replacing
	// text if the selection area covers one or more characters.
	InsertCommand struct {
		DefaultCommand
		// The characters to insert
		Characters string
	}

	// The LeftDeleteCommand deletes characters to the left of the
	// current selection or the current selection if it is not empty.
	LeftDeleteCommand struct {
		DefaultCommand
	}

	// The RightDeleteCommand deletes characters to the right of the
	// current selection or the current selection if it is not empty.
	RightDeleteCommand struct {
		DefaultCommand
	}

	// The DeleteWordCommand deletes one word to right or left
	// depending on forward variable
	DeleteWordCommand struct {
		DefaultCommand
		Forward bool
	}
)

func (c *InsertCommand) Run(v *View, e *Edit) error {
	sel := v.Sel()
	for i := 0; i < sel.Len(); i++ {
		r := sel.Get(i)
		if r.Size() == 0 {
			v.Insert(e, r.B, c.Characters)
		} else {
			v.Replace(e, r, c.Characters)
		}
	}
	return nil
}

func (c *LeftDeleteCommand) Run(v *View, e *Edit) error {
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
	hasNonEmpty := sel.HasNonEmpty()
	i := 0
	for {
		l := sel.Len()
		if i >= l {
			break
		}
		r := sel.Get(i)
		if r.A == r.B && !hasNonEmpty {
			if trim_space {
				_, col := v.Buffer().RowCol(r.A)
				prev_col := r.A - (col - (col-tab_size+(tab_size-1))&^(tab_size-1))
				if prev_col < 0 {
					prev_col = 0
				}
				d := v.Buffer().SubstrR(text.Region{prev_col, r.A})
				i := len(d) - 1
				for r.A > prev_col && i >= 0 && d[i] == ' ' {
					r.A--
					i--
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

func (c *RightDeleteCommand) Run(v *View, e *Edit) error {
	sel := v.Sel()
	hasNonEmpty := sel.HasNonEmpty()
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

func (c *DeleteWordCommand) Run(v *View, e *Edit) error {
	var class int
	if c.Forward {
		class = CLASS_WORD_END | CLASS_PUNCTUATION_END | CLASS_LINE_START
	} else {
		class = CLASS_WORD_START | CLASS_PUNCTUATION_START | CLASS_LINE_END | CLASS_LINE_START
	}

	sel := v.Sel()
	var rs []text.Region
	for i := 0; i < sel.Len(); i++ {
		r := sel.Get(i)
		if r.Empty() {
			p := c.findByClass(r.A, class, v)
			if c.Forward {
				r = text.Region{r.A, p}
			} else {
				r = text.Region{p, r.A}
			}
		}
		rs = append(rs, r)
	}
	sel.Clear()
	sel.AddAll(rs)
	if c.Forward {
		GetEditor().CommandHandler().RunTextCommand(v, "right_delete", nil)
	} else {
		GetEditor().CommandHandler().RunTextCommand(v, "left_delete", nil)
	}
	return nil
}

func (c *DeleteWordCommand) findByClass(point int, class int, v *View) int {
	var end, d int
	if c.Forward {
		d = 1
		end = v.Buffer().Size()
		if point > end {
			point = end
		}
		s := v.Buffer().Substr(text.Region{point, point + 2})
		if strings.Contains(s, "\t") && strings.Contains(s, " ") {
			class = CLASS_WORD_START | CLASS_PUNCTUATION_START | CLASS_LINE_END
		}
	} else {
		d = -1
		end = 0
		if point < end {
			point = end
		}
		s := v.Buffer().Substr(text.Region{point - 2, point})
		if strings.Contains(s, "\t") && strings.Contains(s, " ") {
			class = CLASS_WORD_END | CLASS_PUNCTUATION_END | CLASS_LINE_START
		}
	}
	point += d
	for ; point != end; point += d {
		if v.Classify(point)&class != 0 {
			return point
		}
	}
	return point
}

func init() {
	register([]cmd{
		{"insert", &InsertCommand{}},
		{"left_delete", &LeftDeleteCommand{}},
		{"right_delete", &RightDeleteCommand{}},
		{"delete_word", &DeleteWordCommand{}},
	})
}
