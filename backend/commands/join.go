// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/limetext/lime/backend"
	"strings"
)

type (
	// JoinCommand removes every new line in the
	// selections and the first new line after
	JoinCommand struct {
		DefaultCommand
	}
)

func (c *JoinCommand) Run(v *View, e *Edit) error {
	sel := v.Sel()
	for i := 0; i < sel.Len(); i++ {
		r := sel.Get(i)
		// Removing new line and triming in the selection
		t := v.Buffer().Substr(r)
		t = strings.Replace(t, "\r", "\n", -1)
		slice := strings.Split(t, "\n")
		t = ""
		for i, s := range slice {
			if i == 0 {
				t += s
				continue
			}
			t += " " + strings.TrimLeft(s, " \t")
		}
		v.Replace(e, r, t)
		// Removing the first new line after selection
		liner := v.Buffer().FullLine(r.End())
		line := v.Buffer().Substr(liner)
		line = strings.Replace(line, "\n", "", -1)
		line = strings.Replace(line, "\r", "", -1)
		line = strings.TrimRight(line, " \t")
		// Triming the line after
		nextline := liner.End() + 1
		nextliner := v.Buffer().FullLine(nextline)
		nline := v.Buffer().Substr(nextliner)
		if nline != "" {
			v.Replace(e, nextliner, " "+strings.TrimLeft(nline, " \t"))
		}
		v.Replace(e, liner, line)
	}

	return nil
}

func init() {
	register([]cmd{
		{"join", &JoinCommand{}},
	})
}
