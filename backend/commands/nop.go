// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/limetext/lime/backend"
)

type (
	NopApplicationCommand struct {
		BypassUndoCommand
	}

	NopWindowCommand struct {
		BypassUndoCommand
	}

	NopTextCommand struct {
		BypassUndoCommand
	}
)

func (c *NopApplicationCommand) Run() error {
	return nil
}

func (c *NopWindowCommand) Run(w *Window) error {
	return nil
}

func (c *NopTextCommand) Run(v *View, e *Edit) error {
	return nil
}

func init() {
	registerByName([]namedCmd{
		{"nop", &NopApplicationCommand{}},
		{"nop", &NopWindowCommand{}},
		{"nop", &NopTextCommand{}},
	})
}
