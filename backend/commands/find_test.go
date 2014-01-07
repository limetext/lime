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
		{[]Region{{1, 1}, {2, 2}, {3, 3}, {6, 6}}, []Region{{1, 1}}},
		{[]Region{{2, 2}, {3, 3}, {6, 6}}, []Region{{2, 2}}},
		{[]Region{{5, 5}}, []Region{{5, 5}}},
	}

	for i, test := range tests {
		v.Sel().Clear()
		for _, r := range test.in {
			v.Sel().Add(r)
		}
		ed.CommandHandler().RunTextCommand(v, "single_selection", nil)
		if sr := v.Sel().Regions(); !reflect.DeepEqual(sr, test.exp) {
			t.Errorf("Single Selection test %d failed: %v, %+v", i, sr, test)
		}
	}
}
