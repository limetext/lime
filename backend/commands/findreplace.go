// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	"errors"
	. "github.com/limetext/lime/backend"
	. "github.com/limetext/text"
)

type (
	// The FindUnderExpandCommand extends the selection to the current word
	// if the current selection region is empty.
	// If one character or more is selected, the text buffer is scanned for
	// the next occurrence of the selection and that region too is added to
	// the selection set.
	FindUnderExpandCommand struct {
		DefaultCommand
	}
	// The FindNext command searches for the last search term, starting at
	// the end of the last selection in the buffer, and wrapping around. If
	// it finds the term, it clears the current selections and selects the
	// newly-found regions.
	FindNextCommand struct {
		DefaultCommand
	}

	// The ReplaceNextCommand searches for the "old" argument text,
	// and at the first occurance of the text, replaces it with the
	// "new" argument text. If there are multiple regions, the find
	// starts from the max region.
	ReplaceNextCommand struct {
		DefaultCommand
	}
)

var (
	// Remembers the last sequence of runes searched for.
	lastSearch  []rune
	replaceText string
)

func (c *FindUnderExpandCommand) Run(v *View, e *Edit) error {
	sel := v.Sel()
	rs := sel.Regions()

	if sel.HasEmpty() {
		for i, r := range rs {
			if r2 := v.Buffer().Word(r.A); r2.Size() > r.Size() {
				rs[i] = r2
			}
		}
		sel.Clear()
		sel.AddAll(rs)
		b := v.Buffer()
		lastSearch = b.SubstrR(rs[len(rs)-1])
		return nil
	}
	last := rs[len(rs)-1]
	b := v.Buffer()
	lastSearch = b.SubstrR(last)
	r := v.Find(string(lastSearch), last.End(), IGNORECASE|LITERAL)
	if r.A != -1 {
		sel.Add(r)
	}
	return nil
}

func nextSelection(v *View, search string) (Region, error) {
	sel := v.Sel()
	rs := sel.Regions()
	last := 0

	// Regions are not sorted, so finding the last one requires a search.
	for _, r := range rs {
		last = Max(last, r.End())
	}

	// Start the search right after the last selection.
	start := last
	r := v.Find(search, start, IGNORECASE|LITERAL)
	// If not found yet, search from the start of the buffer to our original
	// starting point.
	if r.A == -1 {
		r = v.Find(search, 0, IGNORECASE|LITERAL)
	}
	// If we found our string, select it.
	if r.A != -1 {
		return r, nil
	}
	return Region{-1, -1}, errors.New("Selection not Found")
}

func (c *FindNextCommand) Run(v *View, e *Edit) error {
	/*
		Correct behavior of FindNext:
			- If there is no previous search, do nothing
			- Find the last region in the buffer, start the
			  search immediately after that.
			- If the search term is found, clear any existing
			  selections, and select the newly-found region.
			- Right now this is doing a case-sensitive search. In ST3
			  that's a setting.
	*/

	// If there is no last search term, nothing to do here.
	if len(lastSearch) == 0 {
		return nil
	}
	newr, err := nextSelection(v, string(lastSearch))
	if err != nil {
		return err
	}
	sel := v.Sel()
	sel.Clear()
	sel.Add(newr)
	return nil
}

func (c *ReplaceNextCommand) Run(v *View, e *Edit) error {
	// use selection function from find.go to get the next region
	selection, err := nextSelection(v, string(lastSearch))
	if err != nil {
		return err
	}
	v.Erase(e, selection)
	v.Insert(e, selection.Begin(), replaceText)
	return nil
}

func init() {
	register([]Command{
		&FindUnderExpandCommand{},
		&FindNextCommand{},
		&ReplaceNextCommand{},
	})
}
