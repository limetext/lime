// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/limetext/lime/backend"
	. "github.com/limetext/text"
	"testing"
)

func TestTranspose(t *testing.T) {
	ed := GetEditor()

	w := ed.NewWindow()
	defer w.Close()

	v := w.NewFile()
	defer func() {
		v.SetScratch(true)
		v.Close()
	}()

	e := v.BeginEdit()
	v.Insert(e, 0, "First&Second")
	v.EndEdit(e)

	// Test with regions at the beginning and end of the buffer,
	// and one beyond the buffer limits.
	rs1 := []Region{{0, 0}, {1, 2}, {3, 3}, {5, 5}, {12, 12}, {20, 20}}
	rs2 := []Region{{0, 0}, {3, 3}, {5, 5}, {12, 12}, {20, 20}}

	// This first should change nothing, because there's a selection with size
	// > 0.
	selectRegions(v, rs1)
	ed.CommandHandler().RunTextCommand(v, "transpose", nil)
	checkContents(t, v, 1, "First&Second")

	selectRegions(v, rs2)
	ed.CommandHandler().RunTextCommand(v, "transpose", nil)
	checkContents(t, v, 2, "Fisr&tSecond")
}

func selectRegions(v *View, rs []Region) {
	v.Sel().Clear()
	for _, r := range rs {
		v.Sel().Add(r)
	}
}

func checkContents(t *testing.T, v *View, tnum int, expected string) {
	b := v.Buffer().Substr(Region{0, v.Buffer().Size()})
	if b != expected {
		t.Error("Test %d: Expected %q but got %q", tnum, expected, b)
	}
}
