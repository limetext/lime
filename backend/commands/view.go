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
	for i, v := range w.Views() {
		if v == w.ActiveView() {
			i++
			if i == len(w.Views()) {
				i = 0
			}
			w.SetActiveView(w.Views()[i])
			break
		}
	}

	return nil
}

func (c *PrevViewCommand) Run(w *Window) error {
	for i, v := range w.Views() {
		if v == w.ActiveView() {
			if i == 0 {
				i = len(w.Views())
			}
			i--
			w.SetActiveView(w.Views()[i])
			break
		}
	}

	return nil
}

func init() {
	register([]Command{
		&NextViewCommand{},
		&PrevViewCommand{},
	})
}
