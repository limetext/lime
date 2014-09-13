// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/limetext/lime/backend"
	"github.com/quarnster/util/text"
	"testing"
)

func TestCopy(t *testing.T) {
	ed := GetEditor()
	w := ed.NewWindow()

	tests := []struct {
		buf     string
		regions []text.Region
		expClip string
		expBuf  string
	}{
		{
			"test string",
			[]text.Region{{1, 3}},
			"es",
			"test string",
		},
		{
			"test\nstring",
			[]text.Region{{3, 6}},
			"t\ns",
			"test\nstring",
		},
		// TODO: Multiregion functionality will probably need to work differently.
		{
			"test string",
			[]text.Region{{1, 3}, {5, 6}},
			"es\ns",
			"test string",
		},
		{
			"test\nstring",
			[]text.Region{{1, 3}, {5, 6}},
			"es\ns",
			"test\nstring",
		},
		{
			"test\nstring",
			[]text.Region{{3, 6}, {9, 10}},
			"t\ns\nn",
			"test\nstring",
		},
		{
			"test string",
			[]text.Region{{5, 6}, {1, 3}},
			"es\ns",
			"test string",
		},
	}

	for i, test := range tests {
		v := w.NewFile()
		v.Buffer().Insert(0, test.buf)
		v.Sel().Clear()
		ed.SetClipboard("")

		for _, r := range test.regions {
			v.Sel().Add(r)
		}

		ed.CommandHandler().RunTextCommand(v, "copy", nil)

		if ed.GetClipboard() != test.expClip {
			t.Errorf("Test %d: Expected clipboard to be %v, but got %v", i, test.expClip, ed.GetClipboard())
		}

		b := v.Buffer().Substr(text.Region{0, v.Buffer().Size()})

		if b != test.expBuf {
			t.Errorf("Test %d: Expected buffer to be %v, but got %v", i, test.expBuf, b)
		}
	}
}

func TestCut(t *testing.T) {
	ed := GetEditor()
	w := ed.NewWindow()

	tests := []struct {
		buf     string
		regions []text.Region
		expClip string
		expBuf  string
	}{
		{
			"test string",
			[]text.Region{{1, 3}},
			"es",
			"tt string",
		},
		{
			"test\nstring",
			[]text.Region{{3, 6}},
			"t\ns",
			"testring",
		},
		// TODO: Multiregion functionality will probably need to work differently.
		{
			"test string",
			[]text.Region{{1, 3}, {5, 6}},
			"es\ns",
			"tt tring",
		},
		{
			"test\nstring",
			[]text.Region{{1, 3}, {5, 6}},
			"es\ns",
			"tt\ntring",
		},
		{
			"test\nstring",
			[]text.Region{{3, 6}, {9, 10}},
			"t\ns\nn",
			"testrig",
		},
		{
			"test string",
			[]text.Region{{5, 6}, {1, 3}},
			"es\ns",
			"tt tring",
		},
	}

	for i, test := range tests {
		v := w.NewFile()
		v.Buffer().Insert(0, test.buf)
		v.Sel().Clear()
		ed.SetClipboard("")

		for _, r := range test.regions {
			v.Sel().Add(r)
		}

		ed.CommandHandler().RunTextCommand(v, "cut", nil)

		if ed.GetClipboard() != test.expClip {
			t.Errorf("Test %d: Expected clipboard to be %v, but got %v", i, test.expClip, ed.GetClipboard())
		}

		b := v.Buffer().Substr(text.Region{0, v.Buffer().Size()})

		if b != test.expBuf {
			t.Errorf("Test %d: Expected buffer to be %v, but got %v", i, test.expBuf, b)
		}
	}
}

func TestPaste(t *testing.T) {
	ed := GetEditor()
	w := ed.NewWindow()

	tests := []struct {
		buf     string
		clip    string
		regions []text.Region
		expBuf  string
	}{
		{
			"test string",
			"test",
			[]text.Region{{1, 1}},
			"ttestest string",
		},
		{
			"test string",
			"test",
			[]text.Region{{1, 3}},
			"ttestt string",
		},
		{
			"test\nstring",
			"test",
			[]text.Region{{3, 6}},
			"testesttring",
		},
		// TODO: Multiregion functionality will probably need to work differently.
		{
			"test string",
			"test",
			[]text.Region{{1, 3}, {5, 6}},
			"ttestt testtring",
		},
		{
			"test\nstring",
			"test",
			[]text.Region{{1, 3}, {5, 6}},
			"ttestt\ntesttring",
		},
		{
			"test\nstring",
			"test",
			[]text.Region{{3, 6}, {9, 10}},
			"testesttritestg",
		},
		{
			"test\nstring",
			"test",
			[]text.Region{{9, 10}, {3, 6}},
			"testesttritestg",
		},
	}

	for i, test := range tests {
		v := w.NewFile()
		v.Buffer().Insert(0, test.buf)
		v.Sel().Clear()

		ed.SetClipboard(test.clip)

		for _, r := range test.regions {
			v.Sel().Add(r)
		}

		ed.CommandHandler().RunTextCommand(v, "paste", nil)

		b := v.Buffer().Substr(text.Region{0, v.Buffer().Size()})

		if b != test.expBuf {
			t.Errorf("Test %d: Expected buffer to be %v, but got %v", i, test.expBuf, b)
		}
	}
}
