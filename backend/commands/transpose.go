// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/limetext/lime/backend"
	. "github.com/limetext/text"
)

type (
	// Transpose: Swap the characters on either side of the cursor,
	// then move the cursor forward one character.
	TransposeCommand struct {
		DefaultCommand
	}
)

func (c *TransposeCommand) Run(v *View, e *Edit) error {
	/*
		Correct behavior of Transpose:
			- ...is actually surprisingly complicated.
			- Transpose behaves differently depending on whether any non-empty
			  region is selected.
			- If there are no non-empty regions, it will swap the characters on
			  either side of the cursor(s), then move all cursors forward one
			  character.
			- If one region is selected, do nothing.
			- If a region is selected and there is another cursor position,
			  expand the cursor to a word and swap them.
			- If two regions are selected, swap them.
			- If more than two regions are selected, rotate them forward.
			- See transpose_test.go for examples.
	*/

	rsnew := RegionSet{}
	rs := v.Sel().Regions()

	if v.Sel().HasNonEmpty() {
		// Build a list of transpose regions based on the current selected regions
		trs := RegionSet{}
		for _, r := range rs {
			if r.Empty() {
				trs.Add(v.Buffer().Word(r.A))
			} else {
				trs.Add(r)
			}
		}
		if trs.Len() < 2 {
			return nil
		}

		srcr := trs.Regions()[trs.Len()-1]
		stxt := v.Buffer().Substr(srcr)
		slen := srcr.Size()
		for i := 0; i < trs.Len(); i++ {
			r := trs.Regions()[i]
			dtxt := v.Buffer().Substr(r)
			dlen := r.Size()
			v.Replace(e, r, stxt)
			trs.Adjust(r.Begin()+1, slen-dlen)
			rsnew.Add(Region{r.Begin(), r.Begin() + slen})
			stxt, slen = dtxt, dlen
		}

	} else {
		for i, r := range rs {
			if i > 0 && r.A-1 == v.Sel().Regions()[i-1].A {
				continue
			}
			rsnew.Add(Region{r.A + 1, r.B + 1})
			if r.A == 0 || r.A >= v.Buffer().Size() {
				continue
			}
			r1 := Region{r.A - 1, r.A}
			r2 := Region{r.A, r.A + 1}
			s1 := v.Buffer().Substr(r1)
			s2 := v.Buffer().Substr(r2)
			v.Replace(e, r1, s2)
			v.Replace(e, r2, s1)
		}
	}

	// Rebuild the active selections
	v.Sel().Clear()
	for _, r := range rsnew.Regions() {
		v.Sel().Add(Region{r.A, r.B})
	}

	return nil
}

func init() {
	register([]Command{
		&TransposeCommand{},
	})
}
