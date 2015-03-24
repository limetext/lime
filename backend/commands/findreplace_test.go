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

type findTest struct {
	text string
	in   []Region
	exp  []Region
	fw   bool
}

func runFindTest(tests []findTest, t *testing.T, commands ...string) {
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
		v.Insert(e, 0, test.text)
		v.EndEdit(e)
		v.Sel().Clear()
		for _, r := range test.in {
			v.Sel().Add(r)
		}

		v.Settings().Set("find_wrap", tests[i].fw)

		for _, command := range commands {
			ed.CommandHandler().RunTextCommand(v, command, nil)
		}
		if sr := v.Sel().Regions(); !reflect.DeepEqual(sr, test.exp) {
			t.Errorf("Test %d: Expected %s, but got %s", i, test.exp, sr)
		}
		e = v.BeginEdit()
		v.Erase(e, Region{0, v.Buffer().Size()})
		v.EndEdit(e)
	}
}

func TestFindUnderExpand(t *testing.T) {
	tests := []findTest{
		{
			"Hello World!\nTest123123\nAbrakadabra\n",
			[]Region{{0, 0}},
			[]Region{{0, 5}},
			true,
		},
		{
			"Hello World!\nTest123123\nAbrakadabra\n",
			[]Region{{19, 20}},
			[]Region{{19, 20}, {22, 23}},
			true,
		},
	}

	runFindTest(tests, t, "find_under_expand")
}

func TestFindNext(t *testing.T) {
	tests := []findTest{
		{
			"Hello World!\nTest123123\nAbrakadabra\n",
			[]Region{{17, 20}},
			[]Region{{17, 20}},
			true,
		},
		// test find_wrap setting true
		{
			"Hello World!\nTest123123\nAbrakadabra\n",
			[]Region{{21, 23}},
			[]Region{{18, 20}},
			true,
		},
		// test find_wrap setting false
		{
			"Hello World!\nTest123123\nAbrakadabra\n",
			[]Region{{21, 23}},
			[]Region{{21, 23}},
			false,
		},
	}

	runFindTest(tests, t, "find_under_expand", "find_next")
}

type replaceTest struct {
	cursors []Region
	in      string
	exp     string
	fw      bool
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

		v.Settings().Set("find_wrap", tests[i].fw)

		replaceText = "f"
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
			true,
		},
		{
			[]Region{{0, 0}, {4, 4}, {8, 8}, {12, 13}},
			"abc abc bac abc abc",
			"abc abc bac abc f",
			true,
		},
		{
			[]Region{{12, 13}, {8, 8}, {4, 4}, {1, 0}},
			"abc abc bac abc abc",
			"abc abc bac abc f",
			true,
		},
		{
			[]Region{{15, 15}},
			"abc abc bac abc abc",
			"abc abc bac abc f",
			true,
		},
		{
			[]Region{{0, 0}},
			"abc abc bac abc abc",
			"abc f bac abc abc",
			true,
		},
		// test find_wrap setting true
		{
			[]Region{{16, 19}},
			"abc abc bac abc abc",
			"f abc bac abc abc",
			true,
		},
		// test find_wrap setting false
		{
			[]Region{{16, 19}},
			"abc abc bac abc abc",
			"abc abc bac abc abc",
			false,
		},
	}

	runReplaceTest(tests, t, "find_under_expand", "replace_next")
}
