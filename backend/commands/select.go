// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/limetext/lime/backend"
	. "github.com/limetext/text"
)

type (
	// The SingleSelectionCommand merges multiple cursors
	// into a single one.
	SingleSelectionCommand struct {
		DefaultCommand
	}
	// The SelectAllCommand selects the whole buffer of the current file
	SelectAllCommand struct {
		DefaultCommand
	}
)

func (c *SingleSelectionCommand) Run(v *View, e *Edit) error {
	/*
		Correct behavior of SingleSelect:
			- Remove all selection regions but the first.
	*/

	r := v.Sel().Get(0)
	v.Sel().Clear()
	v.Sel().Add(r)
	return nil
}

func (c *SelectAllCommand) Run(v *View, e *Edit) error {
	/*
		Correct behavior of SelectAll:
			- Select a single region of (0, view.buffersize())
	*/

	r := Region{0, v.Buffer().Size()}
	v.Sel().Clear()
	v.Sel().Add(r)
	return nil
}

func init() {
	register([]Command{
		&SingleSelectionCommand{},
		&SelectAllCommand{},
	})
}
