// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/quarnster/util/text"
	. "lime/backend"
	"testing"
)

type LineTest struct {
	in_region []Region
	in        string
	exp       string
}

func RunLineTest(command string, testsuite *[]LineTest, t *testing.T) {
	ed := GetEditor()
	w := ed.NewWindow()

	for i, test := range *testsuite {
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
		ed.CommandHandler().RunTextCommand(v, command, nil)
		sr := v.Buffer().Substr(Region{0, v.Buffer().Size()})
		if sr != test.exp {
			t.Errorf("%s test %d failed: %v, %+v", command, i, sr, test)
		}
	}
}

func TestJoinLine(t *testing.T) {
	tests := []LineTest{
		{
			[]Region{{0, 0}},
			`Hello, World 			
			and Goodby!`,
			`Hello, World and Goodby!`,
		},
		{
			[]Region{{0, 0}},
			"Hello, World! \nWhat is this \n ",
			"Hello, World! What is this \n ",
		},
		{
			[]Region{{0, 0}},
			`Hello, World
			and Goodby,
			and Goodby,
			and Goodby!`,
			`Hello, World and Goodby,
			and Goodby,
			and Goodby!`,
		},
	}
	RunLineTest("join_lines", &tests, t)
}
