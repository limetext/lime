// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package backend

import (
	"github.com/quarnster/util/text"
	"testing"
)

func TestUndoStackIndex(t *testing.T) {
	us := UndoStack{}

	w := GetEditor().NewWindow()
	defer w.Close()

	v := w.NewFile()
	defer v.Close()

	e := v.BeginEdit()
	v.Erase(e, text.Region{})
	v.EndEdit(e)
	us.Add(e)

	e = v.BeginEdit()
	v.Erase(e, text.Region{})
	v.EndEdit(e)
	us.Add(e)

	if i := us.index(-1, false); i != 0 {
		t.Errorf("Expected 0, but it got %d", i)
	}
}

func TestUndoStackAdd(t *testing.T) {
	us := UndoStack{}

	w := GetEditor().NewWindow()
	defer w.Close()

	v := w.NewFile()
	defer v.Close()

	e := v.BeginEdit()
	v.Erase(e, text.Region{})
	v.EndEdit(e)
	us.Add(e)

	e = v.BeginEdit()
	v.Erase(e, text.Region{})
	v.EndEdit(e)
	us.Add(e)

	if us.Position() != 2 {
		t.Errorf("Expected the UndoStack positon to be 2, but it was %d", us.Position())
	}
	if len(us.actions) != 2 {
		t.Errorf("Expected the UndoStack to only contain 2 things, but it had %d", len(us.actions))
	}

	us.Undo(false)
	us.Undo(false)

	e = v.BeginEdit()
	v.Erase(e, text.Region{})
	v.EndEdit(e)
	us.Add(e)

	if us.Position() != 1 {
		t.Errorf("Expected the UndoStack positon to be 1, but it was %d", us.Position())
	}
	if len(us.actions) != 1 {
		t.Errorf("Expected the UndoStack to only contain 1 thing, but it had %d", len(us.actions))
	}
}

func TestUndoStackUndo(t *testing.T) {
	us := UndoStack{}
	us.Undo(false)

	if us.Position() != 0 {
		t.Errorf("Expected the UndoStack positon to be 0, but it was %d", us.Position())
	}

	w := GetEditor().NewWindow()
	defer w.Close()

	v := w.NewFile()
	defer v.Close()

	e := v.BeginEdit()
	v.Erase(e, text.Region{})
	v.EndEdit(e)
	us.Add(e)

	e = v.BeginEdit()
	v.Erase(e, text.Region{})
	v.EndEdit(e)
	us.Add(e)

	us.Undo(false)

	if us.Position() != 1 {
		t.Errorf("Expected the UndoStack positon to be 1, but it was %d", us.Position())
	}
	if len(us.actions) != 2 {
		t.Errorf("Expected the UndoStack to only contain 2 things, but it had %d", len(us.actions))
	}
}

func TestUndoStackRedo(t *testing.T) {
	us := UndoStack{}
	us.Redo(false)

	if us.Position() != 0 {
		t.Errorf("Expected the UndoStack positon to be 0, but it was %d", us.Position())
	}

	w := GetEditor().NewWindow()
	defer w.Close()

	v := w.NewFile()
	defer v.Close()

	e := v.BeginEdit()
	v.Erase(e, text.Region{})
	v.EndEdit(e)
	us.Add(e)

	e = v.BeginEdit()
	v.Erase(e, text.Region{})
	v.EndEdit(e)
	us.Add(e)

	us.Undo(false)
	us.Redo(false)

	if us.Position() != 2 {
		t.Errorf("Expected the UndoStack positon to be 1, but it was %d", us.Position())
	}
	if len(us.actions) != 2 {
		t.Errorf("Expected the UndoStack to only contain 2 things, but it had %d", len(us.actions))
	}
}

func TestUndoStackGlueFrom(t *testing.T) {
	us := UndoStack{}

	w := GetEditor().NewWindow()
	defer w.Close()

	v := w.NewFile()
	defer v.Close()

	e := v.BeginEdit()
	v.Erase(e, text.Region{})
	v.EndEdit(e)
	us.Add(e)

	e = v.BeginEdit()
	v.Erase(e, text.Region{})
	v.EndEdit(e)
	us.Add(e)

	e = v.BeginEdit()
	v.Erase(e, text.Region{})
	v.EndEdit(e)
	us.Add(e)

	e = v.BeginEdit()
	v.Erase(e, text.Region{})
	v.EndEdit(e)
	us.Add(e)

	us.GlueFrom(1)

	if us.Position() != 2 {
		t.Errorf("Expected the UndoStack positon to be 2, but it was %d", us.Position())
	}
	if len(us.actions) != 2 {
		t.Errorf("Expected the UndoStack to only contain 2 things, but it had %d", len(us.actions))
	}

	us.GlueFrom(10)

	if us.Position() != 2 {
		t.Errorf("Expected the UndoStack positon to be 2, but it was %d", us.Position())
	}
	if len(us.actions) != 2 {
		t.Errorf("Expected the UndoStack to only contain 2 things, but it had %d", len(us.actions))
	}
}
