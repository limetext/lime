// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/limetext/lime/backend"
	. "github.com/quarnster/util/text"
	"testing"
)

func TestJoin(t *testing.T) {
	type JoinTest struct {
		text   string
		sel    []Region
		expect string
	}

	tests := []JoinTest{
		{
			"a\n\t  bc",
			[]Region{{1, 1}},
			"a bc",
		},
		{
			"abc\r\n\tde",
			[]Region{{0, 0}},
			"abc de",
		},
		{
			"testing \t\t\n join",
			[]Region{{9, 8}},
			"testing join",
		},
		{
			"test\n join\n command\n whith\n multiple\n regions",
			[]Region{{2, 17}, {34, 40}},
			"test join command whith\n multiple regions",
		},
	}

	ed := GetEditor()
	w := ed.NewWindow()

	for i, test := range tests {
		v := w.NewFile()
		e := v.BeginEdit()

		v.Insert(e, 0, test.text)
		v.EndEdit(e)

		v.Sel().Clear()
		for _, r := range test.sel {
			v.Sel().Add(r)
		}

		ed.CommandHandler().RunTextCommand(v, "join", nil)
		if d := v.Buffer().Substr(Region{0, v.Buffer().Size()}); d != test.expect {
			t.Errorf("Test %d:\nExcepted: '%s'\nbut got: '%s'", i, test.expect, d)
		}
	}
}

func TestSelectLines(t *testing.T) {
	type SelectLinesTest struct {
		text    string
		sel     []Region
		forward bool
		expect  []Region
	}

	tests := []SelectLinesTest{
		{
			"abc\ndefg",
			[]Region{{1, 1}},
			true,
			[]Region{{1, 1}, {5, 5}},
		},
		{
			"abcde\nfg",
			[]Region{{4, 4}},
			true,
			[]Region{{4, 4}, {8, 8}},
		},
		{
			"Testing select lines command\nin\nlime text",
			[]Region{{8, 14}, {30, 30}},
			true,
			[]Region{{8, 14}, {30, 30}, {31, 31}, {33, 33}},
		},
		{
			"abc\n\ndefg",
			[]Region{{6, 6}},
			false,
			[]Region{{6, 6}, {4, 4}},
		},
		{
			"Testing select lines command\nin\nlime text",
			[]Region{{30, 36}, {29, 29}},
			false,
			[]Region{{30, 36}, {29, 29}, {0, 0}, {1, 1}},
		},
	}

	ed := GetEditor()
	w := ed.NewWindow()

	for i, test := range tests {
		v := w.NewFile()
		e := v.BeginEdit()

		v.Insert(e, 0, test.text)
		v.EndEdit(e)

		v.Sel().Clear()
		for _, r := range test.sel {
			v.Sel().Add(r)
		}

		ed.CommandHandler().RunTextCommand(v, "select_lines", Args{"forward": test.forward})
		// TODO: Comparing regions
		d := v.Sel()
		if d.Len() != len(test.expect) {
			t.Errorf("Test %d:\nExcepted: '%d' regions, but got: '%d' regions", i, len(test.expect), d.Len())
			t.Errorf("%+v  %+v", test.expect, d.Regions())
		} else {
			var found bool
			for _, r := range test.expect {
				found = false
				for _, r2 := range d.Regions() {
					if r2 == r {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Test %d:\nRegion %+v not found in view regions: %+v", i, r, d.Regions())
				}
			}
		}

	}
}

func TestSwapLine(t *testing.T) {
	type SwapLineTest struct {
		text   string
		sel    []Region
		expect string
	}

	uptests := []SwapLineTest{
		{
			"a\nb",
			[]Region{{2, 2}},
			"b\na",
		},
		{
			"Testing swap line up\ncommand whit multiple\nregions selected\nTesting swap line up\ncommand whit multiple\nregions selected",
			[]Region{{25, 53}, {86, 95}},
			"command whit multiple\nregions selected\nTesting swap line up\ncommand whit multiple\nTesting swap line up\nregions selected",
		},
	}

	dwtests := []SwapLineTest{
		{
			"a\nb",
			[]Region{{1, 1}},
			"b\na",
		},
		{
			"Testing swap line up\ncommand whit multiple\nregions selected\nTesting swap line up\ncommand whit multiple\nregions selected",
			[]Region{{25, 53}, {86, 95}},
			"Testing swap line up\nTesting swap line up\ncommand whit multiple\nregions selected\nregions selected\ncommand whit multiple",
		},
	}

	ed := GetEditor()
	w := ed.NewWindow()

	for i, test := range uptests {
		v := w.NewFile()
		e := v.BeginEdit()

		v.Insert(e, 0, test.text)
		v.EndEdit(e)

		v.Sel().Clear()
		for _, r := range test.sel {
			v.Sel().Add(r)
		}

		ed.CommandHandler().RunTextCommand(v, "swap_line_up", nil)
		if d := v.Buffer().Substr(Region{0, v.Buffer().Size()}); d != test.expect {
			t.Errorf("Test %d:\nExcepted: '%s'\nbut got: '%s'", i, test.expect, d)
		}
	}

	for i, test := range dwtests {
		v := w.NewFile()
		e := v.BeginEdit()

		v.Insert(e, 0, test.text)
		v.EndEdit(e)

		v.Sel().Clear()
		for _, r := range test.sel {
			v.Sel().Add(r)
		}

		ed.CommandHandler().RunTextCommand(v, "swap_line_down", nil)
		if d := v.Buffer().Substr(Region{0, v.Buffer().Size()}); d != test.expect {
			t.Errorf("Test %d:\nExcepted: '%s'\nbut got: '%s'", i, test.expect, d)
		}
	}
}
