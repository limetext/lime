// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/limetext/lime/backend"
)

type (
	// The ReplaceNextCommand searches for the "old" argument text,
	// and at the first occurance of the text, replaces it with the
	// "new" argument text. If there are multiple regions, the find
	// starts from the max region.
	ReplaceNextCommand struct {
		DefaultCommand
	}
)

var replaceText string

func SetReplaceSearchText(replace string) {
	replaceText = replace
}

func (c *ReplaceNextCommand) Run(v *View, e *Edit) error {
	// use selection function from find.go to get the next region
	lastSearch := GetLastSearch()
	selection, err := GetNextSelection(v, e, string(lastSearch))
	if err != nil {
		return err
	}
	v.Erase(e, selection)
	v.Insert(e, selection.Begin(), replaceText)
	return nil
}

func init() {
	register([]Command{
		&ReplaceNextCommand{},
	})
}
