// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/limetext/lime/backend"
	"io/ioutil"
	"path"
)

type (
	NewFileCommand struct {
		DefaultCommand
	}

	OpenFileCommand struct {
		DefaultCommand
		Path string
	}

	OpenDirCommand struct {
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

func (c *OpenDirCommand) Run(w *Window) error {
	files, err := ioutil.ReadDir(c.Path)
	if err != nil {
		return err
	}
	for _, f := range files {
		if !f.IsDir() {
			w.OpenFile(path.Join(c.Path, f.Name()), 0)
		}
	}
	return nil
}

func init() {
	register([]Command{
		&NewFileCommand{},
		&OpenFileCommand{},
		&OpenDirCommand{},
	})
}
