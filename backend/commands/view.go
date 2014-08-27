// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/limetext/lime/backend"
)

type (
	NextViewCommand struct {
		DefaultCommand
	}

	PrevViewCommand struct {
		DefaultCommand
	}
)

func (c *NextViewCommand) Run(w *Window) error {
	ed := GetEditor()
	ed.ActiveWindow().NextView()
	return nil
}

func (c *PrevViewCommand) Run(w *Window) error {
	ed := GetEditor()
	ed.ActiveWindow().PrevView()
	return nil
}

func init() {
	register([]Command{
		&NextViewCommand{},
		&PrevViewCommand{},
	})
}
