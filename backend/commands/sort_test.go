// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/limetext/lime/backend"
	. "github.com/quarnster/util/text"
	"testing"
)

type test struct {
	text             string
	caseSensitive    bool
	reverse          bool
	removeDuplicates bool
	sel              []Region
	expect           string
}

func TestSortLines(t *testing.T) {
	tests := []test{
		{ // Case sensitive
			"B\nc\na",
			true,
			false,
			false,
			[]Region{{0, 5}},
			"B\na\nc",
		},
		{ // Case insensitive
			"text\nSublime\nlime",
			false,
			false,
			false,
			[]Region{{0, 17}},
			"lime\nSublime\ntext",
		},
		{ // Reverse
			"b\nc\na",
			true,
			true,
			false,
			[]Region{{0, 5}},
			"c\nb\na",
		},
		{ // Noncontinuous selection
			"b\nc\na",
			true,
			false,
			false,
			[]Region{{0, 1}, {4, 5}},
			"a\nc\nb",
		},
		{ // Remove duplicates
			"a\nb\na",
			true,
			false,
			true,
			[]Region{{0, 5}},
			"a\nb\n",
		},
		{ // No duplicates removal
			"c\nb\na\nc\n",
			true,
			false,
			false,
			[]Region{{0, 8}},
			"a\nb\nc\nc\n",
		},
	}

	runSortTest(t, tests, "sort_lines")
}

func TestSortSelection(t *testing.T) {
	tests := []test{
		{ // Case sensitive
			"Bca",
			true,
			false,
			false,
			[]Region{{0, 1}, {1, 2}, {2, 3}},
			"Bac",
		},
		{ // Case insensitive
			"textSublimelime",
			false,
			false,
			false,
			[]Region{{0, 4}, {4, 11}, {11, 15}},
			"limeSublimetext",
		},
		{ // Reverse
			"bca",
			true,
			true,
			false,
			[]Region{{0, 1}, {1, 2}, {2, 3}},
			"cba",
		},
		{ // Noncontinuous selection
			"bca",
			true,
			false,
			false,
			[]Region{{0, 1}, {2, 3}},
			"acb",
		},
		{ // Remove duplicates
			"aba",
			true,
			false,
			true,
			[]Region{{0, 1}, {1, 2}, {2, 3}},
			"ab",
		},
		{ // No duplicates removal
			"cbac",
			true,
			false,
			false,
			[]Region{{0, 1}, {1, 2}, {2, 3}, {3, 4}},
			"abcc",
		},
	}

	runSortTest(t, tests, "sort_selection")
}

func runSortTest(t *testing.T, tests []test, command string) {
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

		args := map[string]interface{}{
			"case_sensitive":    test.caseSensitive,
			"reverse":           test.reverse,
			"remove_duplicates": test.removeDuplicates,
		}
		ed.CommandHandler().RunTextCommand(v, command, args)

		if d := v.Buffer().Substr(Region{0, v.Buffer().Size()}); d != test.expect {
			t.Errorf("Test %d: Excepted %#v,\n but got %#v", i, test.expect, d)
		}
	}
}
