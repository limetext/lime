// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/limetext/lime/backend"
	. "github.com/limetext/text"
	"testing"
)

type replaceTest struct {
	cursors []Region
	in      string
	exp     string
}

func runReplaceTest(tests []replaceTest, t *testing.T, commands ...string) {
	ed := GetEditor()
	w := ed.NewWindow()
	defer w.Close()

	v := w.NewFile()
	defer func() {
		v.SetScratch(true)
		v.Close()
	}()

	for i, test := range tests {
		e := v.BeginEdit()
		v.Insert(e, 0, test.in)
		v.EndEdit(e)
		v.Sel().Clear()

		for _, r := range test.cursors {
			v.Sel().Add(r)
		}
		SetReplaceSearchText("f")
		for _, command := range commands {
			ed.CommandHandler().RunTextCommand(v, command, nil)
		}
		if out := v.Buffer().Substr(Region{0, v.Buffer().Size()}); out != test.exp {
			t.Errorf("Test %d failed: %s, %+v", i, out, test)
		}
		e = v.BeginEdit()
		v.Erase(e, Region{0, v.Buffer().Size()})
		v.EndEdit(e)
	}
}

func TestReplaceNext(t *testing.T) {
	tests := []replaceTest{
		{
			[]Region{{1, 1}, {2, 2}, {3, 3}},
			"abc abc bac abc abc",
			"abc f bac abc abc",
		},
		{
			[]Region{{0, 0}, {4, 4}, {8, 8}, {12, 13}},
			"abc abc bac abc abc",
			"abc abc bac abc f",
		},
		{
			[]Region{{12, 13}, {8, 8}, {4, 4}, {1, 0}},
			"abc abc bac abc abc",
			"abc abc bac abc f",
		},
		{
			[]Region{{15, 15}},
			"abc abc bac abc abc",
			"abc abc bac abc f",
		},
		{
			[]Region{{0, 0}},
			"abc abc bac abc abc",
			"abc f bac abc abc",
		},
	}

	runReplaceTest(tests, t, "find_under_expand", "replace_next")
}
