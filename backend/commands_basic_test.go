package backend

import (
	. "lime/backend/primitives"
	"strings"
	"testing"
)

func TestBasic(t *testing.T) {
	data := `Hello world
Test
Goodbye world
`
	ed := GetEditor()
	w := ed.NewWindow()
	v := w.NewView()
	e := v.BeginEdit()
	v.Insert(e, 0, data)
	v.EndEdit(e)

	v.Sel().Clear()
	v.Sel().Add(Region{11, 11})
	v.Sel().Add(Region{16, 16})
	v.Sel().Add(Region{30, 30})
	ed.CommandHandler().RunTextCommand(v, "left_delete", nil)
	if v.Buffer().Data() != `Hello worl
Tes
Goodbye worl
` {
		t.Error(v.Buffer().Data())
	}
	ed.CommandHandler().RunTextCommand(v, "insert", Args{"characters": "a"})
	if d := v.Buffer().Data(); d != "Hello worla\nTesa\nGoodbye worla\n" {
		lines := strings.Split(v.Buffer().Data(), "\n")
		for _, l := range lines {
			t.Errorf("%d: '%s'", len(l), l)
		}
	}

	v.Settings().Set("translate_tabs_to_spaces", true)
	ed.CommandHandler().RunTextCommand(v, "insert", Args{"characters": "\t"})
	if v.Buffer().Data() != "Hello worla \nTesa    \nGoodbye worla   \n" {
		lines := strings.Split(v.Buffer().Data(), "\n")
		for _, l := range lines {
			t.Errorf("%d: '%s'", len(l), l)
		}
	}
	ed.CommandHandler().RunTextCommand(v, "left_delete", nil)
	if d := v.Buffer().Data(); d != "Hello worla\nTesa\nGoodbye worla\n" {
		lines := strings.Split(v.Buffer().Data(), "\n")
		for _, l := range lines {
			t.Errorf("%d: '%s'", len(l), l)
		}
	}

	ed.CommandHandler().RunTextCommand(v, "left_delete", nil)
	if d := v.Buffer().Data(); d != "Hello worl\nTes\nGoodbye worl\n" {
		lines := strings.Split(v.Buffer().Data(), "\n")
		for _, l := range lines {
			t.Errorf("%d: '%s'", len(l), l)
		}
	}

	ed.CommandHandler().RunTextCommand(v, "insert", Args{"characters": "\t"})
	if d := v.Buffer().Data(); d != "Hello worl  \nTes \nGoodbye worl    \n" {
		lines := strings.Split(v.Buffer().Data(), "\n")
		for _, l := range lines {
			t.Errorf("%d: '%s'", len(l), l)
		}
	}

	ed.CommandHandler().RunTextCommand(v, "left_delete", nil)
	if v.Buffer().Data() != "Hello worl\nTes\nGoodbye worl\n" {
		lines := strings.Split(v.Buffer().Data(), "\n")
		for _, l := range lines {
			t.Errorf("%d: '%s'", len(l), l)
		}
	}

}

func TestLeftDelete(t *testing.T) {
	ed := GetEditor()
	w := ed.NewWindow()
	v := w.NewView()
	e := v.BeginEdit()
	v.Insert(e, 0, "12345678")
	v.EndEdit(e)
	v.Sel().Clear()
	v.Sel().Add(Region{1, 1})
	v.Sel().Add(Region{2, 2})
	v.Sel().Add(Region{3, 3})
	v.Sel().Add(Region{4, 4})
	ed.CommandHandler().RunTextCommand(v, "left_delete", nil)
	if d := v.buffer.Data(); d != "5678" {
		t.Error(d)
	}
}
