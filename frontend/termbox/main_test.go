// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package main

import (
	"github.com/limetext/lime/backend"
	"github.com/limetext/termbox-go"
	. "github.com/quarnster/util/text"
	"testing"
)

func TestPadLineRunes(t *testing.T) {
	var testPadData = []struct {
		line     int
		total    int
		expected []rune
	}{
		{1, 3, []rune{' ', ' ', '1', ' '}},
		{10, 3, []rune{' ', '1', '0', ' '}},
		{100, 3, []rune{'1', '0', '0', ' '}},
	}

	for _, p := range testPadData {
		runes := intToRunes(p.line)
		padded := padLineRunes(runes, p.total)

		for i, r := range p.expected {
			if r != padded[i] {
				t.Error("Expected runes to be padded")
			}
		}
	}
}

func TestIntToRunes(t *testing.T) {
	actual := intToRunes(100)
	expected := []rune{'1', '0', '0'}
	for i, r := range expected {
		if actual[i] != r {
			t.Error("Expected int to be transformed to runes")
		}
	}
}

func TestGetCaretStyle(t *testing.T) {
	var testcases = []struct {
		style    string
		inverse  bool
		expected termbox.Attribute
	}{
		{"", false, termbox.AttrUnderline},
		{"block", false, termbox.AttrReverse},
		{"", true, termbox.AttrReverse},
	}

	for _, tc := range testcases {
		if style := getCaretStyle(tc.style, tc.inverse); style != tc.expected {
			t.Errorf("Expected %s, got %s", tc.expected, style)
		}
	}
}

func TestUpdateVisibleRegion(t *testing.T) {
	var (
		fe tbfe
		e  = backend.GetEditor()
		w  = e.NewWindow()
		v  = w.NewFile()
	)

	fe.layout = make(map[*backend.View]layout)
	fe.layout[v] = layout{0, 0, 100, 100 - console_height - 1, Region{}, 0}
	fe.setupCallbacks(v)

	edit := v.BeginEdit()
	v.Insert(edit, 0, "foo")
	v.EndEdit(edit)

	if end := fe.layout[v].visible.End(); end != 3 {
		t.Fatalf("Expected 3, got %d", end)
	}
}
