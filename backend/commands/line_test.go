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
			"abc\n\ndefg",
			[]Region{{6, 6}},
			false,
			[]Region{{6, 6}, {4, 4}},
		},
		{
			"abcde\nfg",
			[]Region{{4, 4}},
			true,
			[]Region{{4, 4}, {8, 8}},
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
		if d := v.Buffer().Substr(Region{0, v.Buffer().Size()}); d != test.expect {
			t.Errorf("Test %d:\nExcepted: '%s'\nbut got: '%s'", i, test.expect, d)
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
			"Testing swap line up\nommand whit multiple\nregions selected\nTesting swap line up\nommand whit multiple\nregions selected",
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
