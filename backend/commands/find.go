// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	"errors"
	. "github.com/limetext/lime/backend"
	. "github.com/limetext/text"
	"strings"
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
	// The SingleSelectionCommand merges multiple cursors
	// into a single one.
	SingleSelectionCommand struct {
		DefaultCommand
	}
	// The SelectAllCommand selects the whole buffer of the current file
	SelectAllCommand struct {
		DefaultCommand
	}
)

func (c *SingleSelectionCommand) Run(v *View, e *Edit) error {
	/*
		Correct behavior of SingleSelect:
			- Remove all selection regions but the first.
	*/

	r := v.Sel().Get(0)
	v.Sel().Clear()
	v.Sel().Add(r)
	return nil
}

// Remembers the last sequence of runes searched for.
var lastSearch []rune

func GetLastSearch() []rune {
	return lastSearch
}

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
	next := last
	size := last.Size()
	next.A += size
	next.B += size
	buf := b.SubstrR(Region{next.A, next.B})
	for next.End() < b.Size() {
		buf[size-1] = b.Index(next.B - 1)
		found := true
		for j, r := range buf {
			if r != lastSearch[j] {
				found = false
				break
			}
		}
		if found {
			sel.Add(next)
			break
		}
		copy(buf, buf[1:])
		next.A += 1
		next.B += 1
	}
	return nil
}

func GetNextSelection(v *View, e *Edit, search string) (Region, error) {
	sel := v.Sel()
	rs := sel.Regions()
	last := 0

	// Ranges are not sorted, so finding the last one requires a search.
	for _, r := range rs {
		last = Max(last, r.End())
	}

	b := v.Buffer()
	// Start the search right after the last selection.
	start := last
	r := Region{start, b.Size()}
	st := b.Substr(r)
	p := search
	size := len(p)
	found := strings.Index(st, p)
	// If not found yet, search from the start of the buffer to our original
	// starting point.
	if found == -1 {
		r = Region{0, start}
		st = b.Substr(r)
		found = strings.Index(st, p)
	}
	// If we found our string, select it.
	if found != -1 {
		newr := Region{r.A + found, r.A + found + size}
		return newr, nil
	}
	return Region{}, errors.New("Selection not Found")
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
	newr, err := GetNextSelection(v, e, string(lastSearch))
	if err != nil {
		return err
	}
	sel := v.Sel()
	sel.Clear()
	sel.Add(newr)
	return nil
}

func (c *SelectAllCommand) Run(v *View, e *Edit) error {
	/*
		Correct behavior of SelectAll:
			- Select a single region of (0, view.buffersize())
	*/

	r := Region{0, v.Buffer().Size()}
	v.Sel().Clear()
	v.Sel().Add(r)
	return nil
}

func init() {
	register([]Command{
		&FindUnderExpandCommand{},
		&FindNextCommand{},
		&SingleSelectionCommand{},
		&SelectAllCommand{},
	})
}
