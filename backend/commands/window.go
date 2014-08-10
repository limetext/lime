// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/limetext/lime/backend"
)

type (
	NewWindowCommand struct {
		DefaultCommand
	}

	CloseWindowCommand struct {
		DefaultCommand
	}

	NewWindowAppCommand struct {
		DefaultCommand
	}

	CloseWindowAppCommand struct {
		DefaultCommand
	}
)

func (c *NewWindowCommand) Run(w *Window) error {
	ed := GetEditor()
	ed.SetActiveWindow(ed.NewWindow())
	return nil
}

func (c *CloseWindowCommand) Run(w *Window) error {
	ed := GetEditor()
	ed.ActiveWindow().Close()
	return nil
}

func (c *NewWindowAppCommand) Run() error {
	ed := GetEditor()
	ed.SetActiveWindow(ed.NewWindow())
	return nil
}

func (c *CloseWindowAppCommand) Run() error {
	ed := GetEditor()
	ed.ActiveWindow().Close()
	return nil
}

func (c *NewWindowAppCommand) IsChecked() bool {
	return false
}

func (c *CloseWindowAppCommand) IsChecked() bool {
	return false
}

func init() {
	register([]cmd{
		{"new_window", &NewWindowCommand{}},
		{"close_window", &CloseWindowCommand{}},
		{"new_window", &NewWindowAppCommand{}},
		{"close_window", &CloseWindowAppCommand{}},
	})
}
