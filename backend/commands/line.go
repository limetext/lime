// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/quarnster/util/text"
	. "lime/backend"
	"strings"
	"unicode"
	"unicode/utf8"
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
	r := Region{0, v.Buffer().Size()}
	s := v.Buffer().Substr(r)
	sa := strings.SplitAfter(s, "\n")
	r = Region{0, utf8.RuneCountInString(sa[0]) + utf8.RuneCountInString(sa[1])}
	v.Replace(e, r, strings.TrimRightFunc(sa[0], unicode.IsSpace)+" "+strings.TrimLeftFunc(sa[1], unicode.IsSpace))
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
