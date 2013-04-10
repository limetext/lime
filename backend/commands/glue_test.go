package commands

import (
	. "lime/backend"
	"testing"
)

func TestGlueCmds(t *testing.T) {
	ed := GetEditor()
	ch := ed.CommandHandler()
	w := ed.NewWindow()
	v := w.NewFile()
	v.SetScratch(true)
	e := v.BeginEdit()
	v.Insert(e, 0, "Hello World!\nTest123123\nAbrakadabra\n")
	v.EndEdit(e)
	v.SetScratch(false)
	ch.RunTextCommand(v, "mark_undo_groups_for_gluing", nil)
	ch.RunTextCommand(v, "insert", Args{"characters": "a"})
	ch.RunTextCommand(v, "insert", Args{"characters": "b"})
	ch.RunTextCommand(v, "insert", Args{"characters": "c"})
	ch.RunTextCommand(v, "glue_marked_undo_groups", nil)
	if v.UndoStack().Position() != 1 {
		t.Error(v.UndoStack().Position())
	} else if d := v.Buffer().String(); d != "Hello World!\nTest123123\nAbrakadabra\nabc" {
		t.Error(d)
	}
	ch.RunTextCommand(v, "undo", nil)
	if d := v.Buffer().String(); d != "Hello World!\nTest123123\nAbrakadabra\n" {
		t.Error(d)
	}
	ch.RunTextCommand(v, "redo", nil)
	if d := v.Buffer().String(); d != "Hello World!\nTest123123\nAbrakadabra\nabc" {
		t.Error(d)
	}
	if v.UndoStack().Position() != 1 {
		t.Error(v.UndoStack().Position())
	} else if d := v.Buffer().String(); d != "Hello World!\nTest123123\nAbrakadabra\nabc" {
		t.Error(d)
	}
	ch.RunTextCommand(v, "undo", nil)
	if d := v.Buffer().String(); d != "Hello World!\nTest123123\nAbrakadabra\n" {
		t.Error(d)
	}

	ch.RunTextCommand(v, "maybe_mark_undo_groups_for_gluing", nil)
	ch.RunTextCommand(v, "insert", Args{"characters": "a"})
	ch.RunTextCommand(v, "maybe_mark_undo_groups_for_gluing", nil)
	ch.RunTextCommand(v, "insert", Args{"characters": "b"})
	ch.RunTextCommand(v, "maybe_mark_undo_groups_for_gluing", nil)
	ch.RunTextCommand(v, "insert", Args{"characters": "c"})
	ch.RunTextCommand(v, "maybe_mark_undo_groups_for_gluing", nil)
	ch.RunTextCommand(v, "glue_marked_undo_groups", nil)
	if v.UndoStack().Position() != 1 {
		t.Error(v.UndoStack().Position())
	} else if d := v.Buffer().String(); d != "Hello World!\nTest123123\nAbrakadabra\nabc" {
		t.Error(d)
	}
	ch.RunTextCommand(v, "undo", nil)
	if d := v.Buffer().String(); d != "Hello World!\nTest123123\nAbrakadabra\n" {
		t.Error(d)
	}
	ch.RunTextCommand(v, "redo", nil)
	if d := v.Buffer().String(); d != "Hello World!\nTest123123\nAbrakadabra\nabc" {
		t.Error(d)
	}
	if v.UndoStack().Position() != 1 {
		t.Error(v.UndoStack().Position())
	} else if d := v.Buffer().String(); d != "Hello World!\nTest123123\nAbrakadabra\nabc" {
		t.Error(d)
	}

	ch.RunTextCommand(v, "mark_undo_groups_for_gluing", nil)
	ch.RunTextCommand(v, "move", Args{"forward": false, "extend": true, "by": "lines"})
	ch.RunTextCommand(v, "move", Args{"forward": false, "extend": true, "by": "lines"})
	ch.RunTextCommand(v, "move", Args{"forward": false, "extend": true, "by": "lines"})
	ch.RunTextCommand(v, "left_delete", nil)
	ch.RunTextCommand(v, "insert", Args{"characters": "a"})
	ch.RunTextCommand(v, "insert", Args{"characters": "b"})
	ch.RunTextCommand(v, "insert", Args{"characters": "c"})
	ch.RunTextCommand(v, "glue_marked_undo_groups", nil)
	if v.UndoStack().Position() != 2 {
		t.Error(v.UndoStack().Position())
	} else if d := v.Buffer().String(); d != "Helabc" {
		t.Error(d)
	}

	ch.RunTextCommand(v, "undo", nil)
	if v.UndoStack().Position() != 1 {
		t.Error(v.UndoStack().Position())
	} else if d := v.Buffer().String(); d != "Hello World!\nTest123123\nAbrakadabra\nabc" {
		t.Error(d)
	}
}
