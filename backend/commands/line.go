// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "lime/backend"
	"strings"
	"unicode"
)

type (
	// The IndentCommand
	// The UnindentCommand
	// The ReindentCommand
	// The SwapLineUpCommand
	// The SwapLineDownCommand
	// The DublicateLineCommand
	// The DeleteLineCommand

	// The JoinLinesCommand joins the following line
	// to the current line, deleting the end of line
	// and replacing all in between whitespace with
	// a single space. For example, the text:
	// "Hello, World
	// and Goodby!"
	// turns in to:
	// "Hello, World and Goodby!"
	JoinLinesCommand struct {
		DefaultCommand
	}
)

func (c *JoinLinesCommand) Run(v *View, e *Edit) error {
	sel := v.Sel()
	r := sel.Regions()
	for i := 0; i < len(r); i++ {
		r0 := v.Buffer().FullLineR(r[i])
		r1 := v.Buffer().FullLine(v.Buffer().FullLineR(r[i]).End())
		s0 := v.Buffer().Substr(r0)
		s1 := v.Buffer().Substr(r1)
		sa0 := strings.SplitAfter(s0, "\n")
		sa1 := strings.SplitAfter(s1, "\n")
		v.Replace(e, r0, strings.TrimRightFunc(sa0[0], unicode.IsSpace)+" "+strings.TrimLeftFunc(sa1[0], unicode.IsSpace))
		v.Replace(e, r1, "")
	}
	return nil
}

func init() {
	register([]cmd{
		// {"indent", &IndentCommand{}},
		// {"unindent", &UnindentCommand{}},
		// {"reindent", &ReindentCommand{}},
		// {"swap_line_up", &SwapLineUpCommand{}},
		// {"swap_line_down", &SwapLineDownCommand{}},
		// {"dublicate_line", &DublicateLineCommand{}},
		// {"delete_line", &DeleteLineCommand{}},
		{"join_lines", &JoinLinesCommand{}},
	})
}
