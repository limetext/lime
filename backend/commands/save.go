// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/limetext/lime/backend"
)

type (
	SaveCommand struct {
		DefaultCommand
	}

	SaveAsCommand struct {
		DefaultCommand
		Name string
	}

	SaveAllCommand struct {
		DefaultCommand
	}
)

func (c *SaveCommand) Run(v *View, e *Edit) error {
	err := v.Save()
	if err != nil {
		return err
	}
	return nil
}

func (c *SaveAsCommand) Run(v *View, e *Edit) error {
	err := v.SaveAs(c.Name)
	if err != nil {
		return err
	}
	return nil
}

func (c *SaveAllCommand) Run(w *Window) error {
	for _, v := range w.Views() {
		if err := v.Save(); err != nil {
			return err
		}
	}
	return nil
}

func init() {
	register([]cmd{
		{"save", &SaveCommand{}},
		{"save_as", &SaveAsCommand{}},
		{"save_all", &SaveAllCommand{}},
	})
}
