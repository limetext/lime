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

func TestSingleSelection(t *testing.T) {
	/*
		Correct behavior of SingleSelect:
			- Remove all selection regions but the first.
	*/

	ed := GetEditor()

	w := ed.NewWindow()
	v := w.NewFile()
	e := v.BeginEdit()
	v.Insert(e, 0, "Hello World!\nTest123123\nAbrakadabra\n")
	v.EndEdit(e)

	type Test struct {
		in  []Region
		exp []Region
	}

	tests := []Test{
		{
			[]Region{{1, 1}, {2, 2}, {3, 3}, {6, 6}},
			[]Region{{1, 1}},
		},
		{
			[]Region{{2, 2}, {3, 3}, {6, 6}},
			[]Region{{2, 2}},
		},
		{
			[]Region{{5, 5}},
			[]Region{{5, 5}},
		},
	}

	for i, test := range tests {
		v.Sel().Clear()
		for _, r := range test.in {
			v.Sel().Add(r)
		}
		ed.CommandHandler().RunTextCommand(v, "single_selection", nil)
		if sr := v.Sel().Regions(); !reflect.DeepEqual(sr, test.exp) {
			t.Errorf("Test %d failed: %v, %+v", i, sr, test)
		}
	}
}

func TestFindUnderExpand(t *testing.T) {
	ed := GetEditor()

	w := ed.NewWindow()
	v := w.NewFile()
	e := v.BeginEdit()
	v.Insert(e, 0, "Hello World!\nTest123123\nAbrakadabra\n")
	v.EndEdit(e)

	type Test struct {
		in  []Region
		exp []Region
	}

	tests := []Test{
		{
			[]Region{{0, 0}},
			[]Region{{0, 5}},
		},
		{
			[]Region{{19, 20}},
			[]Region{{19, 20}, {22, 23}},
		},
	}

	for i, test := range tests {
		v.Sel().Clear()
		for _, r := range test.in {
			v.Sel().Add(r)
		}
		ed.CommandHandler().RunTextCommand(v, "find_under_expand", nil)
		if sr := v.Sel().Regions(); !reflect.DeepEqual(sr, test.exp) {
			t.Errorf("Test %d failed: %v, %+v", i, sr, test)
		}
	}
}

func TestSelectAll(t *testing.T) {
	/*
		Correct behavior of SelectAll:
			- Select a single region of (0, view.buffersize())
	*/

	ed := GetEditor()

	w := ed.NewWindow()
	v := w.NewFile()
	e := v.BeginEdit()
	v.Insert(e, 0, "Hello World!\nTest123123\nAbrakadabra\n")
	v.EndEdit(e)

	type Test struct {
		in  []Region
		exp []Region
	}

	tests := []Test{
		{
			[]Region{{1, 1}, {2, 2}, {3, 3}, {6, 6}},
			[]Region{{0, 36}},
		},
		{
			[]Region{{2, 2}, {3, 3}, {6, 6}},
			[]Region{{0, 36}},
		},
		{
			[]Region{{5, 5}},
			[]Region{{0, 36}},
		},
	}

	for i, test := range tests {
		v.Sel().Clear()
		for _, r := range test.in {
			v.Sel().Add(r)
		}
		ed.CommandHandler().RunTextCommand(v, "select_all", nil)
		if sr := v.Sel().Regions(); !reflect.DeepEqual(sr, test.exp) {
			t.Errorf("Test %d failed: %v, %+v", i, sr, test)
		}
	}
}
