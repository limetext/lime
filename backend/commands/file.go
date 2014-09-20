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

	OpenFileCommand struct {
		DefaultCommand
		Path string
	}
)

func (c *NewFileCommand) Run(w *Window) error {
	ed := GetEditor()
	ed.ActiveWindow().NewFile()
	return nil
}

func (o *OpenFileCommand) Run(w *Window) error {
	w.OpenFile(o.Path, 0)
	return nil
}

func init() {
	register([]Command{
		&NewFileCommand{},
		&OpenFileCommand{},
	})
}
