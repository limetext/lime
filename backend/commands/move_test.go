// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/quarnster/util/text"
	. "lime/backend"
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
			t.Errorf("Move test %d failed: %v, %+v", i, sr, test)
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
			t.Errorf("Move test %d failed: %v, %+v", i, sr, test)
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
