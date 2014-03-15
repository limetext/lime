// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/limetext/lime/backend"
	. "github.com/quarnster/util/text"
	"testing"
)

func TestIndent(t *testing.T) {
	type Test struct {
		text string
		expect string
		sel []Region
	}
	tests := []Test{
		{
			"Life is short\nYou need Python\n",
			"\tLife is short\n\tYou need Python\n",
			[]Region{{3, 5}, {9, 10}, {16, 20}},
		},
		{
			"Life is short\nYou need Python\n",
			"\tLife is short\n\tYou need Python\n",
			[]Region{{3, 20}},
		},
	}
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

		ed.CommandHandler().RunTextCommand(v, "indent", nil)
		if d := v.Buffer().Substr(Region{0, v.Buffer().Size()}); d != test.expect {
			t.Errorf("Test %d: Excepted %s, but got %s", i, test.expect, d)
		}
	}
}
