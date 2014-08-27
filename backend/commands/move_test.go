// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/limetext/lime/backend"
	. "github.com/quarnster/util/text"
	"reflect"
	"testing"
)

func TestMove(t *testing.T) {
	ed := GetEditor()

	w := ed.NewWindow()
	v := w.NewFile()
	e := v.BeginEdit()
	v.Insert(e, 0, "Hello World!\nTest123123\nAbrakadabra\n")
	v.EndEdit(e)

	type Test struct {
		in      []Region
		by      string
		extend  bool
		forward bool
		exp     []Region
		args    Args
	}

	tests := []Test{
		{
			[]Region{{1, 1}, {3, 3}, {6, 6}},
			"characters",
			false,
			true,
			[]Region{{2, 2}, {4, 4}, {7, 7}},
			nil,
		},
		{
			[]Region{{1, 1}, {3, 3}, {6, 6}},
			"characters",
			false,
			false,
			[]Region{{0, 0}, {2, 2}, {5, 5}},
			nil,
		},
		{
			[]Region{{1, 1}, {3, 3}, {10, 6}},
			"characters",
			false,
			true,
			[]Region{{2, 2}, {4, 4}, {7, 7}},
			nil,
		},
		{
			[]Region{{1, 1}, {3, 3}, {10, 6}},
			"characters",
			false,
			false,
			[]Region{{0, 0}, {2, 2}, {5, 5}},
			nil,
		},
		{
			[]Region{{1, 1}, {3, 3}, {10, 6}},
			"characters",
			true,
			true,
			[]Region{{1, 2}, {3, 4}, {10, 7}},
			nil,
		},
		{
			[]Region{{1, 1}, {3, 3}, {10, 6}},
			"characters",
			true,
			false,
			[]Region{{1, 0}, {3, 2}, {10, 5}},
			nil,
		},
		{
			[]Region{{1, 3}, {3, 5}, {10, 7}},
			"characters",
			true,
			true,
			[]Region{{1, 6}, {10, 8}},
			nil,
		},
		{
			[]Region{{1, 1}},
			"stops",
			true,
			true,
			[]Region{{1, 5}},
			Args{"word_end": true},
		},
		{
			[]Region{{1, 1}},
			"stops",
			false,
			true,
			[]Region{{6, 6}},
			Args{"word_begin": true},
		},
		{
			[]Region{{6, 6}},
			"stops",
			false,
			false,
			[]Region{{0, 0}},
			Args{"word_begin": true},
		},
		{
			[]Region{{34, 34}},
			"lines",
			false,
			false,
			[]Region{{23, 23}},
			nil,
		},
		{
			[]Region{{23, 23}},
			"lines",
			false,
			false,
			[]Region{{10, 10}},
			nil,
		},
		{
			[]Region{{100, 100}},
			"lines",
			false,
			false,
			[]Region{{24, 24}},
			nil,
		},
	}
	for i, test := range tests {
		v.Sel().Clear()
		for _, r := range test.in {
			v.Sel().Add(r)
		}
		args := Args{"by": test.by, "extend": test.extend, "forward": test.forward}
		if test.args != nil {
			for k, v := range test.args {
				args[k] = v
			}
		}
		ed.CommandHandler().RunTextCommand(v, "move", args)
		if sr := v.Sel().Regions(); !reflect.DeepEqual(sr, test.exp) {
			t.Errorf("Test %d failed. Expected %v, but got %v: %+v", i, test.exp, sr, test)
		}
	}

	e = v.BeginEdit()
	v.Insert(e, v.Buffer().Size(), "abc")
	v.EndEdit(e)

	tests = []Test{
		{
			[]Region{{100, 100}},
			"lines",
			false,
			false,
			[]Region{{27, 27}},
			nil,
		},
	}
	for i, test := range tests {
		v.Sel().Clear()
		for _, r := range test.in {
			v.Sel().Add(r)
		}
		args := Args{"by": test.by, "extend": test.extend, "forward": test.forward}
		if test.args != nil {
			for k, v := range test.args {
				args[k] = v
			}
		}

		ed.CommandHandler().RunTextCommand(v, "move", args)
		if sr := v.Sel().Regions(); !reflect.DeepEqual(sr, test.exp) {
			t.Errorf("Test %d failed. Expected %v, but got %v: %+v", i, test.exp, sr, test)
		}
	}
}

func TestMoveTo(t *testing.T) {
	/*
	   Correct behavior of MoveTo:
	       - Moves each cursor directly to the indicated position.
	       - If extend, the selection will be extended in the direction of movement
	*/

	ed := GetEditor()
	w := ed.NewWindow()
	v := w.NewFile()
	e := v.BeginEdit()

	v.Insert(e, 0, "Hello World!\nTest123123\nAbrakadabra\n")
	v.EndEdit(e)

	type MoveToTest struct {
		in     []Region
		to     string
		extend bool
		exp    []Region
	}

	singleCursor 			:= []Region{{16, 16}}
	sameLineCursors 		:= []Region{{16, 16}, {17, 17}}
	sameLineCursorsReversed := []Region{{17, 17}, {16, 16}}
	diffLineCursors 		:= []Region{{3, 3},   {17, 17}}
	diffLineCursorsReversed := []Region{{17, 17}, {3, 3}}

	vbufflen := v.Buffer().Size()

	tests := []MoveToTest{
		// BOF move
		{
			singleCursor,
			"bof",
			false,
			[]Region{{0, 0}},
		},
		{
			sameLineCursors,
			"bof",
			false,
			[]Region{{0, 0}},
		},
		{
			sameLineCursorsReversed,
			"bof",
			false,
			[]Region{{0, 0}},
		},
		{
			diffLineCursors,
			"bof",
			false,
			[]Region{{0, 0}},
		},
		{
			diffLineCursorsReversed,
			"bof",
			false,
			[]Region{{0, 0}},
		},

		// BOF extend
		{
			singleCursor,
			"bof",
			true,
			[]Region{{16, 0}},
		},
		{
			sameLineCursors,
			"bof",
			true,
			[]Region{{17, 0}},
		},
		{
			sameLineCursorsReversed,
			"bof",
			true,
			[]Region{{17, 0}},
		},
		{
			diffLineCursors,
			"bof",
			true,
			[]Region{{17, 0}},
		},
		{
			diffLineCursorsReversed,
			"bof",
			true,
			[]Region{{17, 0}},
		},

		// EOF move
		{
			singleCursor,
			"eof",
			false,
			[]Region{{vbufflen, vbufflen}},
		},
		{
			sameLineCursors,
			"eof",
			false,
			[]Region{{vbufflen, vbufflen}},
		},
		{
			sameLineCursorsReversed,
			"eof",
			false,
			[]Region{{vbufflen, vbufflen}},
		},
		{
			diffLineCursors,
			"eof",
			false,
			[]Region{{vbufflen, vbufflen}},
		},
		{
			diffLineCursorsReversed,
			"eof",
			false,
			[]Region{{vbufflen, vbufflen}},
		},

		// EOF extend
		{
			singleCursor,
			"eof",
			true,
			[]Region{{16, vbufflen}},
		},
		{
			sameLineCursors,
			"eof",
			true,
			[]Region{{16, vbufflen}},
		},
		{
			sameLineCursorsReversed,
			"eof",
			true,
			[]Region{{16, vbufflen}},
		},
		{
			diffLineCursors,
			"eof",
			true,
			[]Region{{3, vbufflen}},
		},
		{
			diffLineCursorsReversed,
			"eof",
			true,
			[]Region{{3, vbufflen}},
		},

		// BOL move
		{
			singleCursor,
			"bol",
			false,
			[]Region{{13, 13}},
		},
		{
			sameLineCursors,
			"bol",
			false,
			[]Region{{13, 13}},
		},
		{
			sameLineCursorsReversed,
			"bol",
			false,
			[]Region{{13, 13}},
		},
		{
			diffLineCursors,
			"bol",
			false,
			[]Region{{0, 0}, {13, 13}},
		},
		{
			diffLineCursorsReversed,
			"bol",
			false,
			[]Region{{13, 13}, {0, 0}},
		},

		// BOL extend
		{
			singleCursor,
			"bol",
			true,
			[]Region{{16, 13}},
		},
		{
			sameLineCursors,
			"bol",
			true,
			[]Region{{17, 13}},
		},
		{
			sameLineCursorsReversed,
			"bol",
			true,
			[]Region{{17, 13}},
		},
		{
			diffLineCursors,
			"bol",
			true,
			[]Region{{3, 0}, {17, 13}},
		},
		{
			diffLineCursorsReversed,
			"bol",
			true,
			[]Region{{17, 13}, {3, 0}},
		},

		// EOL move
		{
			singleCursor,
			"eol",
			false,
			[]Region{{23, 23}},
		},
		{
			sameLineCursors,
			"eol",
			false,
			[]Region{{23, 23}},
		},
		{
			sameLineCursorsReversed,
			"eol",
			false,
			[]Region{{23, 23}},
		},
		{
			diffLineCursors,
			"eol",
			false,
			[]Region{{12, 12}, {23, 23}},
		},
		{
			diffLineCursorsReversed,
			"eol",
			false,
			[]Region{{23, 23}, {12, 12}},
		},

		// EOL extend
		{
			singleCursor,
			"eol",
			true,
			[]Region{{16, 23}},
		},
		{
			sameLineCursors,
			"eol",
			true,
			[]Region{{16, 23}},
		},
		{
			sameLineCursorsReversed,
			"eol",
			true,
			[]Region{{16, 23}},
		},
		{
			diffLineCursors,
			"eol",
			true,
			[]Region{{3, 12}, {17, 23}},
		},
		{
			diffLineCursorsReversed,
			"eol",
			true,
			[]Region{{17, 23}, {3, 12}},
		},
	}

	for i, test := range tests {
		v.Sel().Clear()
		for _, r := range test.in {
			v.Sel().Add(r)
		}
		args := Args{"to": test.to, "extend": test.extend}
		ed.CommandHandler().RunTextCommand(v, "move_to", args)
		if sr := v.Sel().Regions(); !reflect.DeepEqual(sr, test.exp) {
			t.Errorf("Test %d failed. Expected %v, but got %v: %+v", i, test.exp, sr, test)
		}
	}
}

type scfe struct {
	DummyFrontend
	show Region
}

func (f *scfe) VisibleRegion(v *View) Region {
	s := v.Buffer().Line(v.Buffer().TextPoint(3*3, 1))
	e := v.Buffer().Line(v.Buffer().TextPoint(6*3, 1))
	return Region{s.Begin(), e.End()}
}

func (f *scfe) Show(v *View, r Region) {
	f.show = r
}

func TestScrollLines(t *testing.T) {
	var fe scfe
	ed := GetEditor()
	ed.SetFrontend(&fe)
	ch := ed.CommandHandler()
	w := ed.NewWindow()
	v := w.NewFile()
	e := v.BeginEdit()
	for i := 0; i < 10; i++ {
		v.Insert(e, 0, "Hello World!\nTest123123\nAbrakadabra\n")
	}
	v.EndEdit(e)
	ch.RunTextCommand(v, "scroll_lines", Args{"amount": 0})

	if c := v.Buffer().Line(v.Buffer().TextPoint(3*3, 1)); fe.show.Begin() != c.Begin() {
		t.Errorf("Expected %v, but got %v", c, fe.show)
	}

	ch.RunTextCommand(v, "scroll_lines", Args{"amount": 1})
	if c := v.Buffer().Line(v.Buffer().TextPoint(3*3-1, 1)); fe.show.Begin() != c.Begin() {
		t.Errorf("Expected %v, but got %v", c, fe.show)
	}
	t.Log(fe.VisibleRegion(v), v.Buffer().Line(v.Buffer().TextPoint(6*3+1, 1)))
	ch.RunTextCommand(v, "scroll_lines", Args{"amount": -1})
	if c := v.Buffer().Line(v.Buffer().TextPoint(6*3+1, 1)); fe.show.Begin() != c.Begin() {
		t.Errorf("Expected %v, but got %v", c, fe.show)
	}
}
