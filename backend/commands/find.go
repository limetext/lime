// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/limetext/lime/backend"
	. "github.com/quarnster/util/text"
)

type (
	// The FindUnderExpandCommand extends the selection to the current word
	// if the current selection region is empty.
	// If one character or more is selected, the text buffer is scanned for
	// the next occurrence of the selection and that region too is added to
	// the selection set.
	FindUnderExpandCommand struct {
		DefaultCommand
	}
	// The SingleSelectionCommand merges multiple cursors
	// into a single one.
	SingleSelectionCommand struct {
		DefaultCommand
	}
	// The AllSelectionCommand selects the whole buffer of the current file
	AllSelectionCommand struct {
		DefaultCommand
	}
)

func (c *SingleSelectionCommand) Run(v *View, e *Edit) error {
	r := v.Sel().Get(0)
	v.Sel().Clear()
	v.Sel().Add(r)
	return nil
}

func (c *FindUnderExpandCommand) Run(v *View, e *Edit) error {
	sel := v.Sel()
	rs := sel.Regions()

	he := sel.HasEmpty()
	if he {
		for i, r := range rs {
			if r2 := v.Buffer().Word(r.A); r2.Size() > r.Size() {
				rs[i] = r2
			}
		}
		sel.Clear()
		sel.AddAll(rs)
		return nil
	}
	last := rs[len(rs)-1]
	b := v.Buffer()
	data := b.SubstrR(last)
	next := last
	size := last.Size()
	next.A += size
	next.B += size
	buf := b.SubstrR(Region{next.A, next.B})
	for next.End() < b.Size() {
		buf[size-1] = b.Index(next.B - 1)
		found := true
		for j, r := range buf {
			if r != data[j] {
				found = false
				break
			}
		}
		if found {
			sel.Add(next)
			break
		}
		copy(buf, buf[1:])
		next.A += 1
		next.B += 1
	}
	return nil
}

func (c *AllSelectionCommand) Run(v *View, e *Edit) error {
	r := Region{0, v.Buffer().Size()}
	v.Sel().Clear()
	v.Sel().Add(r)
	return nil
}

func init() {
	registerByName([]namedCmd{
		{"find_under_expand", &FindUnderExpandCommand{}},
		{"single_selection", &SingleSelectionCommand{}},
		{"select_all", &AllSelectionCommand{}},
	})
}
