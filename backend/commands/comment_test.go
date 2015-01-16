// Copyright 2015 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/limetext/lime/backend"
	"github.com/limetext/text"
	"testing"
)

func TestToggleComment(t *testing.T) {
	tests := []struct {
		r   []text.Region
		in  string
		exp string
	}{
		{
			[]text.Region{{0, 3}},

			"test",
			"// test",
		},
		{
			[]text.Region{{0, 6}},

			"// test",
			"test",
		},
		{
			[]text.Region{{0, 5}},

			"//test",
			"test",
		},
		{
			[]text.Region{{0, 8}},

			"//   test",
			"  test",
		},
		{
			[]text.Region{{0, 7}},

			"    test",
			"    // test",
		},
		{
			[]text.Region{{0, 10}},

			"    // test",
			"    test",
		},
		{
			[]text.Region{{0, 9}},

			"    //test",
			"    test",
		},
		{
			[]text.Region{{0, 12}},

			"    //   test",
			"      test",
		},
		{
			[]text.Region{{0, 8}},

			"\t    test",
			"\t    // test",
		},
		{
			[]text.Region{{0, 11}},

			"\t    // test",
			"\t    test",
		},
		{
			[]text.Region{{0, 10}},

			"\t    //test",
			"\t    test",
		},
		{
			[]text.Region{{0, 13}},

			"\t    //   test",
			"\t      test",
		},
	}

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
		v.Insert(e, 0, test.in)
		v.EndEdit(e)

		v.Sel().Clear()
		if test.r != nil {
			for _, r := range test.r {
				v.Sel().Add(r)
			}
		}
		ed.CommandHandler().RunTextCommand(v, "toggle_comment", nil)
		sr := v.Buffer().Substr(text.Region{0, v.Buffer().Size()})
		if sr != test.exp {
			t.Errorf("%s test %d failed: %v, %+v", "toggle_comment", i, sr, test)
		}
	}
}
