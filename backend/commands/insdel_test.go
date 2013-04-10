package commands

import (
	. "lime/backend"
	. "lime/backend/primitives"
	"reflect"
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
	v := w.NewFile()
	e := v.BeginEdit()
	v.Insert(e, 0, data)
	v.EndEdit(e)

	v.Sel().Clear()
	v.Sel().Add(Region{11, 11})
	v.Sel().Add(Region{16, 16})
	v.Sel().Add(Region{30, 30})
	ed.CommandHandler().RunTextCommand(v, "left_delete", nil)
	if v.Buffer().String() != `Hello worl
Tes
Goodbye worl
` {
		t.Error(v.Buffer().String())
	}
	ed.CommandHandler().RunTextCommand(v, "insert", Args{"characters": "a"})
	if d := v.Buffer().String(); d != "Hello worla\nTesa\nGoodbye worla\n" {
		lines := strings.Split(v.Buffer().String(), "\n")
		for _, l := range lines {
			t.Errorf("%d: '%s'", len(l), l)
		}
	}

	v.Settings().Set("translate_tabs_to_spaces", true)
	ed.CommandHandler().RunTextCommand(v, "insert", Args{"characters": "\t"})
	if v.Buffer().String() != "Hello worla \nTesa    \nGoodbye worla   \n" {
		lines := strings.Split(v.Buffer().String(), "\n")
		for _, l := range lines {
			t.Errorf("%d: '%s'", len(l), l)
		}
	}
	ed.CommandHandler().RunTextCommand(v, "left_delete", nil)
	if d := v.Buffer().String(); d != "Hello worla\nTesa\nGoodbye worla\n" {
		lines := strings.Split(v.Buffer().String(), "\n")
		for _, l := range lines {
			t.Errorf("%d: '%s'", len(l), l)
		}
	}

	ed.CommandHandler().RunTextCommand(v, "left_delete", nil)
	if d := v.Buffer().String(); d != "Hello worl\nTes\nGoodbye worl\n" {
		lines := strings.Split(v.Buffer().String(), "\n")
		for _, l := range lines {
			t.Errorf("%d: '%s'", len(l), l)
		}
	}

	ed.CommandHandler().RunTextCommand(v, "insert", Args{"characters": "\t"})
	if d := v.Buffer().String(); d != "Hello worl  \nTes \nGoodbye worl    \n" {
		lines := strings.Split(v.Buffer().String(), "\n")
		for _, l := range lines {
			t.Errorf("%d: '%s'", len(l), l)
		}
	}

	ed.CommandHandler().RunTextCommand(v, "left_delete", nil)
	if v.Buffer().String() != "Hello worl\nTes\nGoodbye worl\n" {
		lines := strings.Split(v.Buffer().String(), "\n")
		for _, l := range lines {
			t.Errorf("%d: '%s'", len(l), l)
		}
	}

	v.Buffer().Erase(0, len(v.Buffer().String()))
	v.Buffer().Insert(0, "€þıœəßðĸʒ×ŋµåäö𝄞")
	orig := "€þıœəßðĸʒ×ŋµåäö𝄞"
	if d := v.Buffer().String(); d != orig {
		t.Errorf("%s\n\t%v\n\t%v", d, []byte(d), []byte(orig))
	} else {
		t.Logf("ref %s\n\t%v\n\t%v", d, []byte(d), []byte(orig))
	}
	v.Sel().Clear()
	v.Sel().Add(Region{3, 3})
	v.Sel().Add(Region{6, 6})
	v.Sel().Add(Region{9, 9})
	ed.CommandHandler().RunTextCommand(v, "left_delete", nil)
	exp := "€þœəðĸ×ŋµåäö𝄞"
	if d := v.Buffer().String(); d != exp {
		t.Errorf("%s\n\t%v\n\t%v", d, []byte(d), []byte(exp))
	}
}

func TestLeftDelete(t *testing.T) {
	ed := GetEditor()
	w := ed.NewWindow()
	v := w.NewFile()

	type Test struct {
		in, out []Region
		text    string
		ins     string
	}

	tests := []Test{
		{[]Region{{1, 1}, {2, 2}, {3, 3}, {4, 4}}, []Region{{0, 0}}, "5678", "12345678"},
		{[]Region{{1, 1}, {3, 3}, {5, 5}, {7, 7}}, []Region{{0, 0}, {1, 1}, {2, 2}, {3, 3}}, "2468", "12345678"},
		{[]Region{{1, 3}}, []Region{{1, 1}}, "145678", "12345678"},
		{[]Region{{3, 1}}, []Region{{1, 1}}, "145678", "12345678"},
		{[]Region{{100, 5}}, []Region{{93, 5}}, "abc\nd", "abc\ndef\nghi\n"}, // Yes, this is indeed what ST3 does too.
	}

	for i, test := range tests {
		v.Sel().Clear()
		e := v.BeginEdit()
		v.Insert(e, 0, test.ins)
		v.EndEdit(e)

		v.Sel().Clear()
		for _, r := range test.in {
			v.Sel().Add(r)
		}
		var s2 RegionSet
		for _, r := range test.out {
			s2.Add(r)
		}

		ed.CommandHandler().RunTextCommand(v, "left_delete", nil)
		if d := v.Buffer().String(); d != test.text {
			t.Errorf("Test %02d: Expected %s, but got %s", i, test.text, d)
		} else if !reflect.DeepEqual(*v.Sel(), s2) {
			t.Errorf("Test %02d: Expected %v, but have %v", i, s2, v.Sel())
		}
		ed.CommandHandler().RunTextCommand(v, "undo", nil)
		ed.CommandHandler().RunTextCommand(v, "undo", nil)
	}
}

func TestInsert(t *testing.T) {
	ed := GetEditor()
	ch := ed.CommandHandler()
	w := ed.NewWindow()
	v := w.NewFile()
	e := v.BeginEdit()
	v.Insert(e, 0, "Hello World!\nTest123123\nAbrakadabra\n")
	v.EndEdit(e)

	type Test struct {
		in   []Region
		data string
		expd string
		expr []Region
	}

	tests := []Test{
		{
			[]Region{{1, 1}, {3, 3}, {6, 6}},
			"a",
			"Haelalo aWorld!\nTest123123\nAbrakadabra\n",
			[]Region{{2, 2}, {5, 5}, {9, 9}},
		},
		{
			[]Region{{1, 1}, {3, 3}, {6, 9}},
			"a",
			"Haelalo ald!\nTest123123\nAbrakadabra\n",
			[]Region{{2, 2}, {5, 5}, {9, 9}},
		},
		{
			[]Region{{1, 1}, {3, 3}, {6, 9}},
			"€þıœəßðĸʒ×ŋµåäö𝄞",
			"H€þıœəßðĸʒ×ŋµåäö𝄞el€þıœəßðĸʒ×ŋµåäö𝄞lo €þıœəßðĸʒ×ŋµåäö𝄞ld!\nTest123123\nAbrakadabra\n",
			[]Region{{17, 17}, {35, 35}, {54, 54}},
		},
	}
	for i, test := range tests {
		v.Sel().Clear()
		for _, r := range test.in {
			v.Sel().Add(r)
		}
		ed.CommandHandler().RunTextCommand(v, "insert", Args{"characters": test.data})
		if d := v.Buffer().String(); d != test.expd {
			t.Errorf("Insert test %d failed: %s", i, d)
		}
		if sr := v.Sel().Regions(); !reflect.DeepEqual(sr, test.expr) {
			t.Errorf("Insert test %d failed: %v", i, sr)
		}
		ch.RunTextCommand(v, "undo", nil)
	}
}
