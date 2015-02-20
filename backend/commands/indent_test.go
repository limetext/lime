// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/limetext/lime/backend"
	. "github.com/limetext/text"
	"testing"
)

type indentTest struct {
	text                     string
	translate_tabs_to_spaces interface{}
	tab_size                 interface{}
	sel                      []Region
	expect                   string
}

func runIndentTest(t *testing.T, tests []indentTest, command string) {
	ed := GetEditor()
	w := ed.NewWindow()
	defer w.Close()

	for i, test := range tests {
		v := w.NewFile()
		defer func() {
			v.SetScratch(true)
			v.Close()
		}()

		e := v.BeginEdit()
		v.Insert(e, 0, test.text)
		v.EndEdit(e)

		v.Sel().Clear()
		for _, r := range test.sel {
			v.Sel().Add(r)
		}
		v.Settings().Set("translate_tabs_to_spaces", test.translate_tabs_to_spaces)
		v.Settings().Set("tab_size", test.tab_size)

		ed.CommandHandler().RunTextCommand(v, command, nil)
		if d := v.Buffer().Substr(Region{0, v.Buffer().Size()}); d != test.expect {
			t.Errorf("Test %d: Expected \n%s, but got \n%s", i, test.expect, d)
		}
	}
}

func TestIndent(t *testing.T) {
	tests := []indentTest{
		{ // translate_tabs_to_spaces = false
			// indent should be "\t"
			"a\n b\n  c\n   d\n",
			false,
			4,
			[]Region{{0, 1}},
			"\ta\n b\n  c\n   d\n",
		},
		{ // translate_tabs_to_spaces = nil
			// indent should be "\t"
			"a\n b\n  c\n   d\n",
			nil,
			1,
			[]Region{{0, 1}},
			"\ta\n b\n  c\n   d\n",
		},
		{ // translate_tabs_to_spaces = true and tab_size = 2
			// indent should be "  "
			"a\n b\n  c\n   d\n",
			true,
			2,
			[]Region{{0, 1}},
			"  a\n b\n  c\n   d\n",
		},
		{ // translate_tabs_to_spaces = true and tab_size = nil
			// indent should be "    "
			"a\n b\n  c\n   d\n",
			true,
			nil,
			[]Region{{0, 1}},
			"    a\n b\n  c\n   d\n",
		},
		{ // region include the 1st line and the 4th line
			// indent should add to the begining of 1st and 4th line
			"a\n b\n  c\n   d\n",
			false,
			1,
			[]Region{{0, 1}, {11, 12}},
			"\ta\n b\n  c\n\t   d\n",
		},
		{ // region selected reversely
			// should perform indent
			"a\n b\n  c\n   d\n",
			false,
			1,
			[]Region{{3, 0}},
			"\ta\n\t b\n  c\n   d\n",
		},
	}

	runIndentTest(t, tests, "indent")
}

func TestUnindent(t *testing.T) {
	tests := []indentTest{
		{ // translate_tabs_to_spaces = false
			// indent should be "\t"
			"\ta\n  b\n      c\n\t  d\n",
			false,
			4,
			[]Region{{0, 19}},
			"a\nb\n  c\n  d\n",
		},
		{ // translate_tabs_to_spaces = nil
			// indent should be "\t"
			"\ta\n b\n  c\n   d\n",
			nil,
			1,
			[]Region{{0, 1}},
			"a\n b\n  c\n   d\n",
		},
		{ // translate_tabs_to_spaces = true and tab_size = 2
			// indent should be "  "
			"  a\n b\n  c\n   d\n",
			true,
			2,
			[]Region{{0, 1}},
			"a\n b\n  c\n   d\n",
		},
		{ // translate_tabs_to_spaces = true and tab_size = nil
			// indent should be "    "
			"    a\n b\n  c\n   d\n",
			true,
			nil,
			[]Region{{0, 1}},
			"a\n b\n  c\n   d\n",
		},
		{ // region include the 1st line and the 4th line
			// unindent should remove from the begining of 1st and 4th line
			"\ta\n b\n  c\n \t   d\n",
			false,
			1,
			[]Region{{0, 1}, {11, 12}},
			"a\n b\n  c\n\t   d\n",
		},
		{ // region selected reversely
			// should perform unindent
			"\ta\n\t b\n  c\n   d\n",
			false,
			4,
			[]Region{{3, 0}},
			"a\n b\n  c\n   d\n",
		},
		{ // empty strings
			// should continue
			"",
			false,
			nil,
			[]Region{{0, 0}},
			"",
		},
	}

	runIndentTest(t, tests, "unindent")
}
