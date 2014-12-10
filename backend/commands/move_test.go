// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/limetext/lime/backend"
	. "github.com/limetext/text"
	"reflect"
	"testing"
)

type MoveTest struct {
	in      []Region
	by      string
	extend  bool
	forward bool
	exp     []Region
	args    Args
}

func runMoveTest(tests []MoveTest, t *testing.T, text string) {
	ed := GetEditor()

	w := ed.NewWindow()
	defer w.Close()

	v := w.NewFile()
	defer func() {
		v.SetScratch(true)
		v.Close()
	}()

	e := v.BeginEdit()
	v.Insert(e, 0, text)
	v.EndEdit(e)

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

func TestMove(t *testing.T) {
	tests := []MoveTest{
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
		// Try moving outside the buffer
		{
			[]Region{{0, 0}},
			"lines",
			false,
			false,
			[]Region{{0, 0}},
			nil,
		},
		{
			[]Region{{36, 36}},
			"lines",
			false,
			true,
			[]Region{{36, 36}},
			nil,
		},
		{
			[]Region{{0, 0}},
			"characters",
			false,
			false,
			[]Region{{0, 0}},
			nil,
		},
		{
			[]Region{{36, 36}},
			"characters",
			false,
			true,
			[]Region{{36, 36}},
			nil,
		},
	}
	runMoveTest(tests, t, "Hello World!\nTest123123\nAbrakadabra\n")

	tests = []MoveTest{
		{
			[]Region{{100, 100}},
			"lines",
			false,
			false,
			[]Region{{27, 27}},
			nil,
		},
	}
	runMoveTest(tests, t, "abc")
}

func TestMoveTo(t *testing.T) {
	/*
	   Correct behavior of MoveTo:
	       - Moves each cursor directly to the indicated position.
	       - If extend, the selection will be extended in the direction of movement
	*/

	ed := GetEditor()
	w := ed.NewWindow()
	defer w.Close()

	v := w.NewFile()
	defer func() {
		v.SetScratch(true)
		v.Close()
	}()

	e := v.BeginEdit()

	v.Insert(e, 0, "Hello World!\nTest123123\nAbrakadabra\n")
	v.EndEdit(e)

	type MoveToTest struct {
		in     []Region
		to     string
		extend bool
		exp    []Region
	}

	singleCursor := []Region{{16, 16}}

	sameLineCursors := []Region{{16, 16}, {17, 17}}
	sameLineCursorsReversed := []Region{{17, 17}, {16, 16}}

	diffLineCursors := []Region{{3, 3}, {17, 17}}
	diffLineCursorsReversed := []Region{{17, 17}, {3, 3}}

	singleForwardSelection := []Region{{15, 18}}
	singleBackwardSelection := []Region{{18, 15}}

	sameLineForwardSelections := []Region{{15, 18}, {20, 21}}
	sameLineForwardSelectionsReversed := []Region{{20, 21}, {15, 18}}
	sameLineBackwardSelections := []Region{{18, 15}, {21, 20}}
	sameLineBackwardSelectionsReversed := []Region{{21, 20}, {18, 15}}
	sameLineForwardThenBackwardSelections := []Region{{15, 18}, {21, 20}}
	sameLineForwardThenBackwardSelectionsReversed := []Region{{21, 20}, {15, 18}}
	sameLineBackwardThenForwardSelections := []Region{{18, 15}, {20, 21}}
	sameLineBackwardThenForwardSelectionsReversed := []Region{{20, 21}, {18, 15}}

	diffLineForwardSelections := []Region{{4, 6}, {20, 21}}
	diffLineForwardSelectionsReversed := []Region{{20, 21}, {4, 6}}
	diffLineBackwardSelections := []Region{{6, 4}, {21, 20}}
	diffLineBackwardSelectionsReversed := []Region{{21, 20}, {6, 4}}
	diffLineForwardThenBackwardSelections := []Region{{4, 6}, {21, 20}}
	diffLineForwardThenBackwardSelectionsReversed := []Region{{21, 20}, {4, 6}}
	diffLineBackwardThenForwardSelections := []Region{{6, 4}, {20, 21}}
	diffLineBackwardThenForwardSelectionsReversed := []Region{{20, 21}, {6, 4}}

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
		{
			singleForwardSelection,
			"bof",
			false,
			[]Region{{0, 0}},
		},
		{
			singleBackwardSelection,
			"bof",
			false,
			[]Region{{0, 0}},
		},
		{
			sameLineForwardSelections,
			"bof",
			false,
			[]Region{{0, 0}},
		},
		{
			sameLineForwardSelectionsReversed,
			"bof",
			false,
			[]Region{{0, 0}},
		},
		{
			sameLineBackwardSelections,
			"bof",
			false,
			[]Region{{0, 0}},
		},
		{
			sameLineBackwardSelectionsReversed,
			"bof",
			false,
			[]Region{{0, 0}},
		},
		{
			sameLineForwardThenBackwardSelections,
			"bof",
			false,
			[]Region{{0, 0}},
		},
		{
			sameLineForwardThenBackwardSelectionsReversed,
			"bof",
			false,
			[]Region{{0, 0}},
		},
		{
			sameLineBackwardThenForwardSelections,
			"bof",
			false,
			[]Region{{0, 0}},
		},
		{
			sameLineBackwardThenForwardSelectionsReversed,
			"bof",
			false,
			[]Region{{0, 0}},
		},
		{
			diffLineForwardSelections,
			"bof",
			false,
			[]Region{{0, 0}},
		},
		{
			diffLineForwardSelectionsReversed,
			"bof",
			false,
			[]Region{{0, 0}},
		},
		{
			diffLineBackwardSelections,
			"bof",
			false,
			[]Region{{0, 0}},
		},
		{
			diffLineBackwardSelectionsReversed,
			"bof",
			false,
			[]Region{{0, 0}},
		},
		{
			diffLineForwardThenBackwardSelections,
			"bof",
			false,
			[]Region{{0, 0}},
		},
		{
			diffLineForwardThenBackwardSelectionsReversed,
			"bof",
			false,
			[]Region{{0, 0}},
		},
		{
			diffLineBackwardThenForwardSelections,
			"bof",
			false,
			[]Region{{0, 0}},
		},
		{
			diffLineBackwardThenForwardSelectionsReversed,
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
		{
			singleForwardSelection,
			"bof",
			true,
			[]Region{{15, 0}},
		},
		{
			singleBackwardSelection,
			"bof",
			true,
			[]Region{{18, 0}},
		},
		{
			sameLineForwardSelections,
			"bof",
			true,
			[]Region{{20, 0}},
		},
		{
			sameLineForwardSelectionsReversed,
			"bof",
			true,
			[]Region{{20, 0}},
		},
		{
			sameLineBackwardSelections,
			"bof",
			true,
			[]Region{{21, 0}},
		},
		{
			sameLineBackwardSelectionsReversed,
			"bof",
			true,
			[]Region{{21, 0}},
		},
		{
			sameLineForwardThenBackwardSelections,
			"bof",
			true,
			[]Region{{21, 0}},
		},
		{
			sameLineForwardThenBackwardSelectionsReversed,
			"bof",
			true,
			[]Region{{21, 0}},
		},
		{
			sameLineBackwardThenForwardSelections,
			"bof",
			true,
			[]Region{{20, 0}},
		},
		{
			sameLineBackwardThenForwardSelectionsReversed,
			"bof",
			true,
			[]Region{{20, 0}},
		},
		{
			diffLineForwardSelections,
			"bof",
			true,
			[]Region{{20, 0}},
		},
		{
			diffLineForwardSelectionsReversed,
			"bof",
			true,
			[]Region{{20, 0}},
		},
		{
			diffLineBackwardSelections,
			"bof",
			true,
			[]Region{{21, 0}},
		},
		{
			diffLineBackwardSelectionsReversed,
			"bof",
			true,
			[]Region{{21, 0}},
		},
		{
			diffLineForwardThenBackwardSelections,
			"bof",
			true,
			[]Region{{21, 0}},
		},
		{
			diffLineForwardThenBackwardSelectionsReversed,
			"bof",
			true,
			[]Region{{21, 0}},
		},
		{
			diffLineBackwardThenForwardSelections,
			"bof",
			true,
			[]Region{{20, 0}},
		},
		{
			diffLineBackwardThenForwardSelectionsReversed,
			"bof",
			true,
			[]Region{{20, 0}},
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
		{
			singleForwardSelection,
			"eof",
			false,
			[]Region{{vbufflen, vbufflen}},
		},
		{
			singleBackwardSelection,
			"eof",
			false,
			[]Region{{vbufflen, vbufflen}},
		},
		{
			sameLineForwardSelections,
			"eof",
			false,
			[]Region{{vbufflen, vbufflen}},
		},
		{
			sameLineForwardSelectionsReversed,
			"eof",
			false,
			[]Region{{vbufflen, vbufflen}},
		},
		{
			sameLineBackwardSelections,
			"eof",
			false,
			[]Region{{vbufflen, vbufflen}},
		},
		{
			sameLineBackwardSelectionsReversed,
			"eof",
			false,
			[]Region{{vbufflen, vbufflen}},
		},
		{
			sameLineForwardThenBackwardSelections,
			"eof",
			false,
			[]Region{{vbufflen, vbufflen}},
		},
		{
			sameLineForwardThenBackwardSelectionsReversed,
			"eof",
			false,
			[]Region{{vbufflen, vbufflen}},
		},
		{
			sameLineBackwardThenForwardSelections,
			"eof",
			false,
			[]Region{{vbufflen, vbufflen}},
		},
		{
			sameLineBackwardThenForwardSelectionsReversed,
			"eof",
			false,
			[]Region{{vbufflen, vbufflen}},
		},
		{
			diffLineForwardSelections,
			"eof",
			false,
			[]Region{{vbufflen, vbufflen}},
		},
		{
			diffLineForwardSelectionsReversed,
			"eof",
			false,
			[]Region{{vbufflen, vbufflen}},
		},
		{
			diffLineBackwardSelections,
			"eof",
			false,
			[]Region{{vbufflen, vbufflen}},
		},
		{
			diffLineBackwardSelectionsReversed,
			"eof",
			false,
			[]Region{{vbufflen, vbufflen}},
		},
		{
			diffLineForwardThenBackwardSelections,
			"eof",
			false,
			[]Region{{vbufflen, vbufflen}},
		},
		{
			diffLineForwardThenBackwardSelectionsReversed,
			"eof",
			false,
			[]Region{{vbufflen, vbufflen}},
		},
		{
			diffLineBackwardThenForwardSelections,
			"eof",
			false,
			[]Region{{vbufflen, vbufflen}},
		},
		{
			diffLineBackwardThenForwardSelectionsReversed,
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
		{
			singleForwardSelection,
			"eof",
			true,
			[]Region{{15, vbufflen}},
		},
		{
			singleBackwardSelection,
			"eof",
			true,
			[]Region{{18, vbufflen}},
		},
		{
			sameLineForwardSelections,
			"eof",
			true,
			[]Region{{15, vbufflen}},
		},
		{
			sameLineForwardSelectionsReversed,
			"eof",
			true,
			[]Region{{15, vbufflen}},
		},
		{
			sameLineBackwardSelections,
			"eof",
			true,
			[]Region{{18, vbufflen}},
		},
		{
			sameLineBackwardSelectionsReversed,
			"eof",
			true,
			[]Region{{18, vbufflen}},
		},
		{
			sameLineForwardThenBackwardSelections,
			"eof",
			true,
			[]Region{{15, vbufflen}},
		},
		{
			sameLineForwardThenBackwardSelectionsReversed,
			"eof",
			true,
			[]Region{{15, vbufflen}},
		},
		{
			sameLineBackwardThenForwardSelections,
			"eof",
			true,
			[]Region{{18, vbufflen}},
		},
		{
			sameLineBackwardThenForwardSelectionsReversed,
			"eof",
			true,
			[]Region{{18, vbufflen}},
		},
		{
			diffLineForwardSelections,
			"eof",
			true,
			[]Region{{4, vbufflen}},
		},
		{
			diffLineForwardSelectionsReversed,
			"eof",
			true,
			[]Region{{4, vbufflen}},
		},
		{
			diffLineBackwardSelections,
			"eof",
			true,
			[]Region{{6, vbufflen}},
		},
		{
			diffLineBackwardSelectionsReversed,
			"eof",
			true,
			[]Region{{6, vbufflen}},
		},
		{
			diffLineForwardThenBackwardSelections,
			"eof",
			true,
			[]Region{{4, vbufflen}},
		},
		{
			diffLineForwardThenBackwardSelectionsReversed,
			"eof",
			true,
			[]Region{{4, vbufflen}},
		},
		{
			diffLineBackwardThenForwardSelections,
			"eof",
			true,
			[]Region{{6, vbufflen}},
		},
		{
			diffLineBackwardThenForwardSelectionsReversed,
			"eof",
			true,
			[]Region{{6, vbufflen}},
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
		{
			singleForwardSelection,
			"bol",
			false,
			[]Region{{13, 13}},
		},
		{
			singleBackwardSelection,
			"bol",
			false,
			[]Region{{13, 13}},
		},
		{
			sameLineForwardSelections,
			"bol",
			false,
			[]Region{{13, 13}},
		},
		{
			sameLineForwardSelectionsReversed,
			"bol",
			false,
			[]Region{{13, 13}},
		},
		{
			sameLineBackwardSelections,
			"bol",
			false,
			[]Region{{13, 13}},
		},
		{
			sameLineBackwardSelectionsReversed,
			"bol",
			false,
			[]Region{{13, 13}},
		},
		{
			sameLineForwardThenBackwardSelections,
			"bol",
			false,
			[]Region{{13, 13}},
		},
		{
			sameLineForwardThenBackwardSelectionsReversed,
			"bol",
			false,
			[]Region{{13, 13}},
		},
		{
			sameLineBackwardThenForwardSelections,
			"bol",
			false,
			[]Region{{13, 13}},
		},
		{
			sameLineBackwardThenForwardSelectionsReversed,
			"bol",
			false,
			[]Region{{13, 13}},
		},
		{
			diffLineForwardSelections,
			"bol",
			false,
			[]Region{{0, 0}, {13, 13}},
		},
		{
			diffLineForwardSelectionsReversed,
			"bol",
			false,
			[]Region{{13, 13}, {0, 0}},
		},
		{
			diffLineBackwardSelections,
			"bol",
			false,
			[]Region{{0, 0}, {13, 13}},
		},
		{
			diffLineBackwardSelectionsReversed,
			"bol",
			false,
			[]Region{{13, 13}, {0, 0}},
		},
		{
			diffLineForwardThenBackwardSelections,
			"bol",
			false,
			[]Region{{0, 0}, {13, 13}},
		},
		{
			diffLineForwardThenBackwardSelectionsReversed,
			"bol",
			false,
			[]Region{{13, 13}, {0, 0}},
		},
		{
			diffLineBackwardThenForwardSelections,
			"bol",
			false,
			[]Region{{0, 0}, {13, 13}},
		},
		{
			diffLineBackwardThenForwardSelectionsReversed,
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
		{
			singleForwardSelection,
			"bol",
			true,
			[]Region{{15, 13}},
		},
		{
			singleBackwardSelection,
			"bol",
			true,
			[]Region{{18, 13}},
		},
		{
			sameLineForwardSelections,
			"bol",
			true,
			[]Region{{20, 13}},
		},
		{
			sameLineForwardSelectionsReversed,
			"bol",
			true,
			[]Region{{20, 13}},
		},
		{
			sameLineBackwardSelections,
			"bol",
			true,
			[]Region{{21, 13}},
		},
		{
			sameLineBackwardSelectionsReversed,
			"bol",
			true,
			[]Region{{21, 13}},
		},
		{
			sameLineForwardThenBackwardSelections,
			"bol",
			true,
			[]Region{{21, 13}},
		},
		{
			sameLineForwardThenBackwardSelectionsReversed,
			"bol",
			true,
			[]Region{{21, 13}},
		},
		{
			sameLineBackwardThenForwardSelections,
			"bol",
			true,
			[]Region{{20, 13}},
		},
		{
			sameLineBackwardThenForwardSelectionsReversed,
			"bol",
			true,
			[]Region{{20, 13}},
		},
		{
			diffLineForwardSelections,
			"bol",
			true,
			[]Region{{4, 0}, {20, 13}},
		},
		{
			diffLineForwardSelectionsReversed,
			"bol",
			true,
			[]Region{{20, 13}, {4, 0}},
		},
		{
			diffLineBackwardSelections,
			"bol",
			true,
			[]Region{{6, 0}, {21, 13}},
		},
		{
			diffLineBackwardSelectionsReversed,
			"bol",
			true,
			[]Region{{21, 13}, {6, 0}},
		},
		{
			diffLineForwardThenBackwardSelections,
			"bol",
			true,
			[]Region{{4, 0}, {21, 13}},
		},
		{
			diffLineForwardThenBackwardSelectionsReversed,
			"bol",
			true,
			[]Region{{21, 13}, {4, 0}},
		},
		{
			diffLineBackwardThenForwardSelections,
			"bol",
			true,
			[]Region{{6, 0}, {20, 13}},
		},
		{
			diffLineBackwardThenForwardSelectionsReversed,
			"bol",
			true,
			[]Region{{20, 13}, {6, 0}},
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
		{
			singleForwardSelection,
			"eol",
			false,
			[]Region{{23, 23}},
		},
		{
			singleBackwardSelection,
			"eol",
			false,
			[]Region{{23, 23}},
		},
		{
			sameLineForwardSelections,
			"eol",
			false,
			[]Region{{23, 23}},
		},
		{
			sameLineForwardSelectionsReversed,
			"eol",
			false,
			[]Region{{23, 23}},
		},
		{
			sameLineBackwardSelections,
			"eol",
			false,
			[]Region{{23, 23}},
		},
		{
			sameLineBackwardSelectionsReversed,
			"eol",
			false,
			[]Region{{23, 23}},
		},
		{
			sameLineForwardThenBackwardSelections,
			"eol",
			false,
			[]Region{{23, 23}},
		},
		{
			sameLineForwardThenBackwardSelectionsReversed,
			"eol",
			false,
			[]Region{{23, 23}},
		},
		{
			sameLineBackwardThenForwardSelections,
			"eol",
			false,
			[]Region{{23, 23}},
		},
		{
			sameLineBackwardThenForwardSelectionsReversed,
			"eol",
			false,
			[]Region{{23, 23}},
		},
		{
			diffLineForwardSelections,
			"eol",
			false,
			[]Region{{12, 12}, {23, 23}},
		},
		{
			diffLineForwardSelectionsReversed,
			"eol",
			false,
			[]Region{{23, 23}, {12, 12}},
		},
		{
			diffLineBackwardSelections,
			"eol",
			false,
			[]Region{{12, 12}, {23, 23}},
		},
		{
			diffLineBackwardSelectionsReversed,
			"eol",
			false,
			[]Region{{23, 23}, {12, 12}},
		},
		{
			diffLineForwardThenBackwardSelections,
			"eol",
			false,
			[]Region{{12, 12}, {23, 23}},
		},
		{
			diffLineForwardThenBackwardSelectionsReversed,
			"eol",
			false,
			[]Region{{23, 23}, {12, 12}},
		},
		{
			diffLineBackwardThenForwardSelections,
			"eol",
			false,
			[]Region{{12, 12}, {23, 23}},
		},
		{
			diffLineBackwardThenForwardSelectionsReversed,
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
		{
			singleForwardSelection,
			"eol",
			true,
			[]Region{{15, 23}},
		},
		{
			singleBackwardSelection,
			"eol",
			true,
			[]Region{{18, 23}},
		},
		{
			sameLineForwardSelections,
			"eol",
			true,
			[]Region{{15, 23}},
		},
		{
			sameLineForwardSelectionsReversed,
			"eol",
			true,
			[]Region{{15, 23}},
		},
		{
			sameLineBackwardSelections,
			"eol",
			true,
			[]Region{{18, 23}},
		},
		{
			sameLineBackwardSelectionsReversed,
			"eol",
			true,
			[]Region{{18, 23}},
		},
		{
			sameLineForwardThenBackwardSelections,
			"eol",
			true,
			[]Region{{15, 23}},
		},
		{
			sameLineForwardThenBackwardSelectionsReversed,
			"eol",
			true,
			[]Region{{15, 23}},
		},
		{
			sameLineBackwardThenForwardSelections,
			"eol",
			true,
			[]Region{{18, 23}},
		},
		{
			sameLineBackwardThenForwardSelectionsReversed,
			"eol",
			true,
			[]Region{{18, 23}},
		},
		{
			diffLineForwardSelections,
			"eol",
			true,
			[]Region{{4, 12}, {20, 23}},
		},
		{
			diffLineForwardSelectionsReversed,
			"eol",
			true,
			[]Region{{20, 23}, {4, 12}},
		},
		{
			diffLineBackwardSelections,
			"eol",
			true,
			[]Region{{6, 12}, {21, 23}},
		},
		{
			diffLineBackwardSelectionsReversed,
			"eol",
			true,
			[]Region{{21, 23}, {6, 12}},
		},
		{
			diffLineForwardThenBackwardSelections,
			"eol",
			true,
			[]Region{{4, 12}, {21, 23}},
		},
		{
			diffLineForwardThenBackwardSelectionsReversed,
			"eol",
			true,
			[]Region{{21, 23}, {4, 12}},
		},
		{
			diffLineBackwardThenForwardSelections,
			"eol",
			true,
			[]Region{{6, 12}, {20, 23}},
		},
		{
			diffLineBackwardThenForwardSelectionsReversed,
			"eol",
			true,
			[]Region{{20, 23}, {6, 12}},
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
	defer w.Close()

	v := w.NewFile()
	defer func() {
		v.SetScratch(true)
		v.Close()
	}()

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
