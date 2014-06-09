// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/limetext/lime/backend"
	"strings"
)

type (
	// The IndentCommand increments indentation of selection.
	IndentCommand struct {
		DefaultCommand
	}
)

func (c *IndentCommand) Run(v *View, e *Edit) error {
	indent := "\t"
	if t, ok := v.Settings().Get("translate_tabs_to_spaces", false).(bool); ok && t {
		indent = strings.Repeat(" ", getTabSize(v))
	}
	sel := v.Sel()

	// Keep track of the indented lines
	// (go has no set collections, use that instead - struct{} doesn't take up space)
	indented_rows := map[int]struct{}{}

	for i := 0; i < sel.Len(); i++ {
		r := sel.Get(i)
		start_row, _ := v.Buffer().RowCol(r.A)
		end_row, _ := v.Buffer().RowCol(r.B)
		for row := start_row; row <= end_row; row++ {
			if _, ok := indented_rows[row]; !ok {
				// Insert an indent at the beginning of the line
				pos := v.Buffer().TextPoint(row, 0)
				v.Insert(e, pos, indent)
				indented_rows[row] = struct{}{}
			}
		}
	}
	return nil
}

// Return the tab size from the settings, defaulting to 4 if not found.
func getTabSize(v *View) int {
	tab_size := 4
	if t, ok := v.Settings().Get("tab_size", tab_size).(int); ok {
		tab_size = t
	}
	return tab_size
}

func init() {
	register([]cmd{
		{"indent", &IndentCommand{}},
	})
}
