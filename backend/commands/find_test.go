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
	in  []Region
	exp []Region
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

	e := v.BeginEdit()
	v.Insert(e, 0, "Hello World!\nTest123123\nAbrakadabra\n")
	v.EndEdit(e)

	for i, test := range tests {
		v.Sel().Clear()
		for _, r := range test.in {
			v.Sel().Add(r)
		}
		for _, command := range commands {
			ed.CommandHandler().RunTextCommand(v, command, nil)
		}
		if sr := v.Sel().Regions(); !reflect.DeepEqual(sr, test.exp) {
			t.Errorf("Test %d failed: %v, %+v", i, sr, test)
		}
	}
}

func TestFindUnderExpand(t *testing.T) {
	tests := []findTest{
		{
			[]Region{{0, 0}},
			[]Region{{0, 5}},
		},
		{
			[]Region{{19, 20}},
			[]Region{{19, 20}, {22, 23}},
		},
	}

	runFindTest(tests, t, "find_under_expand")
}

func TestFindNext(t *testing.T) {
	tests := []findTest{
		{
			[]Region{{17, 20}},
			[]Region{{17, 20}},
		},
		{
			[]Region{{21, 23}},
			[]Region{{18, 20}},
		},
	}

	runFindTest(tests, t, "find_under_expand", "find_next")
}
