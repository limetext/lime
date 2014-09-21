// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/limetext/lime/backend"
	. "github.com/quarnster/util/text"
	"testing"
)

func TestUndoRedoCommands(t *testing.T) {
	ed := GetEditor()
	ch := ed.CommandHandler()
	w := ed.NewWindow()
	defer w.Close()

	v := w.NewFile()
	defer func() {
		v.SetScratch(true)
		v.Close()
	}()

	edit := v.BeginEdit()
	v.Insert(edit, 0, "abcd")
	v.EndEdit(edit)
	v.Sel().Clear()
	r := []Region{
		{0, 0},
		{1, 1},
		{2, 2},
		{3, 3},
	}
	for _, r2 := range r {
		v.Sel().Add(r2)
	}

	edit = v.BeginEdit()
	for _, ins := range "1234" {
		for i := 0; i < v.Sel().Len(); i++ {
			v.Insert(edit, v.Sel().Get(i).Begin(), string(ins))
		}
	}
	v.EndEdit(edit)

	if v.Buffer().Substr(Region{0, v.Buffer().Size()}) != "1234a1234b1234c1234d" {
		t.Error(v.Buffer().Substr(Region{0, v.Buffer().Size()}))
	}
	ch.RunTextCommand(v, "undo", nil)
	if v.Buffer().Substr(Region{0, v.Buffer().Size()}) != "abcd" {
		t.Error("expected 'abcd', but got: ", v.Buffer().Substr(Region{0, v.Buffer().Size()}))
	}
	ch.RunTextCommand(v, "redo", nil)
	if v.Buffer().Substr(Region{0, v.Buffer().Size()}) != "1234a1234b1234c1234d" {
		t.Error("expected '1234a1234b1234c1234d', but got: ", v.Buffer().Substr(Region{0, v.Buffer().Size()}))
	}

	v.Sel().Clear()
	r = []Region{
		{0, 0},
		{5, 5},
		{10, 10},
		{15, 15},
	}
	for _, r2 := range r {
		v.Sel().Add(r2)
	}

	edit = v.BeginEdit()
	for _, ins := range []string{"hello ", "world"} {
		for i := 0; i < v.Sel().Len(); i++ {
			v.Insert(edit, v.Sel().Get(i).Begin(), ins)
		}
	}
	v.EndEdit(edit)

	if v.Buffer().Substr(Region{0, v.Buffer().Size()}) != "hello world1234ahello world1234bhello world1234chello world1234d" {
		t.Error(v.Buffer().Substr(Region{0, v.Buffer().Size()}))
	}
	ch.RunTextCommand(v, "undo", nil)

	if v.Buffer().Substr(Region{0, v.Buffer().Size()}) != "1234a1234b1234c1234d" {
		t.Error("expected '1234a1234b1234c1234d', but got: ", v.Buffer().Substr(Region{0, v.Buffer().Size()}))
	}
	ch.RunTextCommand(v, "undo", nil)
	if v.Buffer().Substr(Region{0, v.Buffer().Size()}) != "abcd" {
		t.Error("expected 'abcd', but got: ", v.Buffer().Substr(Region{0, v.Buffer().Size()}))
	}
	ch.RunTextCommand(v, "undo", nil)
	if v.Buffer().Substr(Region{0, v.Buffer().Size()}) != "" {
		t.Error("expected '', but got: ", v.Buffer().Substr(Region{0, v.Buffer().Size()}))
	}
	v.UndoStack().Redo(true)
	if v.Buffer().Substr(Region{0, v.Buffer().Size()}) != "abcd" {
		t.Error("expected 'abcd', but got: ", v.Buffer().Substr(Region{0, v.Buffer().Size()}))
	}

	v.UndoStack().Redo(true)
	if v.Buffer().Substr(Region{0, v.Buffer().Size()}) != "1234a1234b1234c1234d" {
		t.Error("expected '1234a1234b1234c1234d', but got: ", v.Buffer().Substr(Region{0, v.Buffer().Size()}))
	}

	v.UndoStack().Redo(true)
	if v.Buffer().Substr(Region{0, v.Buffer().Size()}) != "hello world1234ahello world1234bhello world1234chello world1234d" {
		t.Error(v.Buffer().Substr(Region{0, v.Buffer().Size()}))
	}
}
