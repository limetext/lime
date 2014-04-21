// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/limetext/lime/backend"
	. "github.com/quarnster/util/text"
	"testing"
)

func TestJoin(t *testing.T) {
	type Test struct {
		text   string
		sel    []Region
		expect string
	}

	tests := []Test{
		{
			"a\n\t  bc",
			[]Region{{1, 1}},
			"a bc",
		},
		{
			"abc\r\n\tde",
			[]Region{{0, 0}},
			"abc de",
		},
		{
			"testing \t\t\n join",
			[]Region{{9, 8}},
			"testing join",
		},
		{
			"test\n join\n command\n whith\n multiple\n regions",
			[]Region{{2, 17}, {34, 40}},
			"test join command whith\n multiple regions",
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

		ed.CommandHandler().RunTextCommand(v, "join", nil)
		if d := v.Buffer().Substr(Region{0, v.Buffer().Size()}); d != test.expect {
			t.Errorf("Test %d:\nExcepted: '%s'\nbut got: '%s'", i, test.expect, d)
		}
	}
}
