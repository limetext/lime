// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/limetext/lime/backend"
)

type (
	NewFileCommand struct {
		DefaultCommand
	}
)

func (c *NewFileCommand) Run(w *Window) error {
	ed := GetEditor()
	ed.ActiveWindow().NewFile()
	return nil
}

func init() {
	register([]Command{
		&NewFileCommand{},
	})
}
