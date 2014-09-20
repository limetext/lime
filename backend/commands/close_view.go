// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/limetext/lime/backend"
)

type (
	CloseViewCommand struct {
		DefaultCommand
	}
)

func (c *CloseViewCommand) Run(w *Window) error {
	w.ActiveView().Close()
	return nil
}

func init() {
	register([]Command{
		&CloseViewCommand{},
	})
}
