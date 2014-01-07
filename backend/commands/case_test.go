// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/quarnster/util/text"
	. "lime/backend"
	"testing"
)

func TestTitleCase(t *testing.T) {

	// Please note the bizarre  capitalization of the first L in he'Ll...  This is due to a bug in go's strings
	// library.  I'm going to try to get them to fix it...  If not, maybe we'll have
	// to write our own Title Casing function.
	type Test struct {
		in_region []Region
		in        string
		exp       string
	}

	tests := []Test{
		/*single selection*/
		{
			[]Region{{24, 51}},
			"Give a man a match, and he'll be warm for a minute, but set him on fire, and he'll be warm for the rest of his life.",
			"Give a man a match, and He'Ll Be Warm For A Minute, but set him on fire, and he'll be warm for the rest of his life.",
		},
		/*multiple selection*/
		{
			[]Region{{0, 17}, {52, 71}},
			"Give a man a match, and he'll be warm for a minute, but set him on fire, and he'll be warm for the rest of his life.",
			"Give A Man A Match, and he'll be warm for a minute, But Set Him On Fire, and he'll be warm for the rest of his life.",
		},

		/*no selection*/
		{
			nil,
			"Give a man a match, and he'll be warm for a minute, but set him on fire, and he'll be warm for the rest of his life.",
			"Give a man a match, and he'll be warm for a minute, but set him on fire, and he'll be warm for the rest of his life.",
		},
		/*unicode*/
		{
			[]Region{{0, 12}},
			"ничего себе!",
			"Ничего Себе!",
		},
	}

	ed := GetEditor()
	w := ed.NewWindow()

	for i, test := range tests {
		v := w.NewFile()
		e := v.BeginEdit()
		v.Insert(e, 0, test.in)
		v.EndEdit(e)

		v.Sel().Clear()
		if test.in_region != nil {
			for _, r := range test.in_region {
				v.Sel().Add(r)
			}
		}
		ed.CommandHandler().RunTextCommand(v, "title_case", nil)
		sr := v.Buffer().Substr(Region{0, v.Buffer().Size()})
		if sr != test.exp {
			t.Errorf("Title Case test %d failed: %v, %+v", i, sr, test)
		}
	}
}

func TestSwapCase(t *testing.T) {

	type Test struct {
		in_region []Region
		in        string
		exp       string
	}

	tests := []Test{
		{
			[]Region{{0, 13}},
			"Hello, World!",
			"hELLO, wORLD!",
		},
		{
			[]Region{{0, 11}},
			"ПрИвЕт, МиР",
			"пРиВеТ, мИр",
		},
	}

	ed := GetEditor()
	w := ed.NewWindow()

	for i, test := range tests {
		v := w.NewFile()
		e := v.BeginEdit()
		v.Insert(e, 0, test.in)
		v.EndEdit(e)

		v.Sel().Clear()
		if test.in_region != nil {
			for _, r := range test.in_region {
				v.Sel().Add(r)
			}
		}
		ed.CommandHandler().RunTextCommand(v, "swap_case", nil)
		sr := v.Buffer().Substr(Region{0, v.Buffer().Size()})
		if sr != test.exp {
			t.Errorf("Swap Case test %d failed: %v, %+v", i, sr, test)
		}
	}
}
