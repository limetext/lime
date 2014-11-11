// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/limetext/lime/backend"
	. "github.com/limetext/text"
	"testing"
)

func TestTranspose(t *testing.T) {
	ed := GetEditor()

	w := ed.NewWindow()
	defer w.Close()

	v := w.NewFile()
	defer func() {
		v.SetScratch(true)
		v.Close()
	}()

	type Test struct {
		start      string
		regions    []Region
		expect     string
		newregions []Region
	}

	// Test results produced using ST3, and commands like:
	// >>> rs.add_all([sublime.Region(0,0),sublime.Region(4,7),sublime.Region(9,9),sublime.Region(15,16)])
	// >>> v.run_command("transpose")
	// >>> print([a for a in v.sel()])
	tests := []Test{
		{
			// Simple test with just one cursor position
			"one",
			[]Region{{1, 1}},
			"noe",
			[]Region{{2, 2}},
		},
		{
			// Test with several cursors, including one at the beginning of
			// the buffer, which doesn't transpose, and one beyond the end.
			"one two three four",
			[]Region{{0, 0}, {2, 2}, {5, 5}, {20, 20}},
			"oen wto three four",
			[]Region{{1, 1}, {3, 3}, {6, 6}, {21, 21}},
		},
		{
			// Similar test, but with two adjacent cursors. The second one gets
			// dropped, and doesn't transpose.
			"one two three four",
			[]Region{{0, 0}, {1, 1}, {5, 5}},
			"one wto three four",
			[]Region{{1, 1}, {6, 6}},
		},
		{
			// Test with a single region. This should do nothing.
			"one two three four",
			[]Region{{6, 10}},
			"one two three four",
			[]Region{{6, 10}},
		},
		{
			// Test with two regions of different sizes
			"one two three four",
			[]Region{{4, 7}, {8, 13}},
			"one three two four",
			[]Region{{4, 9}, {10, 13}},
		},
		{
			// Test with four regions
			"one two three four",
			[]Region{{0, 3}, {4, 7}, {8, 13}, {14, 18}},
			"four one two three",
			[]Region{{0, 4}, {5, 8}, {9, 12}, {13, 18}},
		},
		{
			// Test with one region and three cursors. The newline at the end
			// of these lines is a workaround for a bug in the Buffer.Word()
			// call, which currently has problems if it finds EOF at the end
			// of the word.
			"one two three four\n",
			[]Region{{0, 0}, {4, 7}, {9, 9}, {16, 16}},
			"four one two three\n",
			[]Region{{0, 4}, {5, 8}, {9, 12}, {13, 18}},
		},
		{
			// Test with two regions and two cursors
			"one two three four",
			[]Region{{0, 0}, {4, 7}, {9, 9}, {15, 16}},
			"o one two fthreeur",
			[]Region{{0, 1}, {2, 5}, {6, 9}, {11, 16}},
		},
	}

	for i, test := range tests {
		// Load the starting text into the buffer
		e := v.BeginEdit()
		v.Erase(e, Region{0, v.Buffer().Size()})
		v.Insert(e, 0, test.start)
		v.EndEdit(e)

		// Add the starting selections
		v.Sel().Clear()
		for _, r := range test.regions {
			v.Sel().Add(r)
		}

		ed.CommandHandler().RunTextCommand(v, "transpose", nil)

		b := v.Buffer().Substr(Region{0, v.Buffer().Size()})
		if b != test.expect {
			t.Errorf("Test %d: Expected %q; got %q", i, test.expect, b)
		}
		rs := v.Sel()
		if rs.Len() == 0 {
			t.Errorf("Test %d: No regions after transpose!", i)
		}
		for ir, r := range v.Sel().Regions() {
			if r != test.newregions[ir] {
				t.Logf("Expected: %s", test.newregions)
				t.Logf("Got     : %s", v.Sel().Regions())
				t.Errorf("Test %d: Selected regions wrong after transpose", i)
				break
			}
		}
	}
}
