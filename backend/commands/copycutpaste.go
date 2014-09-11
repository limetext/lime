// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/limetext/lime/backend"
	"strings"
)

type (
	CopyCommand struct {
		DefaultCommand
	}

	CutCommand struct {
		DefaultCommand
	}

	PasteCommand struct {
		DefaultCommand
	}
)

func (c *CopyCommand) Run(v *View, e *Edit) error {
	// TODO: Copy the entire line if there is no selection.

	s := make([]string, len(v.Sel().Regions()))
	for i, r := range v.Sel().Regions() {
		s[i] = v.Buffer().Substr(r)
	}

	GetEditor().SetClipboard(strings.Join(s, "\n"))

	return nil
}

func (c *CutCommand) Run(v *View, e *Edit) error {
	// TODO: Cut the entire line if there is no selection.

	s := make([]string, len(v.Sel().Regions()))
	for i, r := range v.Sel().Regions() {
		s[i] = v.Buffer().Substr(r)
		v.Erase(e, r)
	}

	GetEditor().SetClipboard(strings.Join(s, "\n"))

	return nil
}

func (c *PasteCommand) Run(v *View, e *Edit) error {
	// TODO: Paste the entire line on the line before the cursor if a
	//		 line was autocopied.
	// TODO: If the number of regions copied matches the number of regions
	//		 selected, paste the regions one at a time.

	ed := GetEditor()

	for _, r := range v.Sel().Regions() {
		v.Replace(e, r, ed.GetClipboard())
	}

	return nil
}

func init() {
	register([]Command{
		&CopyCommand{},
		&CutCommand{},
		&PasteCommand{},
	})
}
