// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package keys

import (
	"io/ioutil"
	"testing"

	"github.com/limetext/lime/backend/loaders"
	"github.com/limetext/lime/backend/util"
)

func TestLoadKeyBindingsFromJSON(t *testing.T) {
	tests := []string{
		"testdata/Default.sublime-keymap",
	}
	for i, fn := range tests {
		if d, err := ioutil.ReadFile(fn); err != nil {
			t.Errorf("Test %d: Couldn't load file %s: %s", i, fn, err)
		} else {
			var bindings KeyBindings
			if err := loaders.LoadJSON(d, &bindings); err != nil {
				t.Errorf("Test %d: Error on LoadJSON: %s", i, err)
			}
		}
	}
}

func TestUnmarshalError(t *testing.T) {
	var bindings KeyBindings
	if err := bindings.UnmarshalJSON([]byte(``)); err == nil {
		t.Errorf("Expected error on loading empty string")
	}
}

func TestDropLessEqualKeys(t *testing.T) {
	fn := "testdata/Default.sublime-keymap"
	d, err := ioutil.ReadFile(fn)
	if err != nil {
		t.Fatalf("Couldn't read %s: %s", fn, err)
	}

	var bd KeyBindings
	if err = loaders.LoadJSON(d, &bd); err != nil {
		t.Fatalf("Error loading json: %s", err)
	}
	bd.DropLessEqualKeys(1)
	if cmd := bd.Bindings[0].Command; cmd != "test2" {
		t.Errorf("Expected Command %s, but got %s", "test2", cmd)
	}
}

func TestSetParent(t *testing.T) {
	fn := "testdata/Default.sublime-keymap"
	fnp := "testdata/test.sublime-keymap"
	var (
		bd KeyBindings
		p  HasKeyBindings
	)

	d, err := ioutil.ReadFile(fn)
	if err != nil {
		t.Fatalf("Couldn't read %s: %s", fn, err)
	}
	if err = loaders.LoadJSON(d, &bd); err != nil {
		t.Fatalf("Error loading json: %s", err)
	}
	d, err = ioutil.ReadFile(fnp)
	if err != nil {
		t.Fatalf("Couldn't read %s: %s", fn, err)
	}
	if err = loaders.LoadJSON(d, p.KeyBindings()); err != nil {
		t.Fatalf("Error loading json: %s", err)
	}

	p.KeyBindings().seqIndex = 10
	bd.SetParent(&p)
	if bd.seqIndex != p.KeyBindings().seqIndex {
		t.Fatalf("Expected parent and child seqIndex be equal %d != %d", p.KeyBindings().seqIndex, bd.seqIndex)
	}

	ret := bd.Filter(KeyPress{Key: 'd', Ctrl: true})
	if ret.Len() != 1 {
		t.Fatalf("Expected ret keyBindings len %d, but got %d", 1, ret.Len())
	}
	if ret.parent.KeyBindings().Len() != 1 {
		t.Fatalf("Expected ret parent keyBindings len %d, but got %d", 1, ret.parent.KeyBindings().Len())
	}
	if cmd := ret.Bindings[0].Command; cmd != "test4" {
		t.Errorf("Expected Command %s, but got %s", "test4", cmd)
	}
	if cmd := ret.parent.KeyBindings().Bindings[0].Command; cmd != "t1" {
		t.Errorf("Expected Command %s, but got %s", "t1", cmd)
	}
}

func TestParent(t *testing.T) {
	fn := "testdata/Default.sublime-keymap"
	fnp := "testdata/test.sublime-keymap"
	var (
		bd KeyBindings
		p  HasKeyBindings
	)

	d, err := ioutil.ReadFile(fn)
	if err != nil {
		t.Fatalf("Couldn't read %s: %s", fn, err)
	}
	if err = loaders.LoadJSON(d, &bd); err != nil {
		t.Fatalf("Error loading json: %s", err)
	}

	d, err = ioutil.ReadFile(fnp)
	if err != nil {
		t.Fatalf("Couldn't read %s: %s", fn, err)
	}
	if err = loaders.LoadJSON(d, p.KeyBindings()); err != nil {
		t.Fatalf("Error loading json: %s", err)
	}

	bd.SetParent(&p)

	if cmd := bd.Parent().KeyBindings().Bindings[0].Command; cmd != "t2" {
		t.Errorf("Expected Command %s, but got %s", "t2", cmd)
	}
}

func TestKeyBindingsFilter(t *testing.T) {
	tests := []struct {
		kp    KeyPress
		count int
	}{
		{
			KeyPress{Key: 'i', Ctrl: true},
			2,
		},
		{
			KeyPress{Key: 'i'},
			1,
		},
	}

	if d, err := ioutil.ReadFile("testdata/Default.sublime-keymap"); err != nil {
		t.Fatal(err)
	} else {
		var bindings KeyBindings
		loaders.LoadJSON(d, &bindings)

		for i, test := range tests {
			if b := bindings.Filter(test.kp); b.Len() != test.count {
				t.Errorf("Test %d: Expected %d bindings, but got %d", i, test.count, b.Len())
			}
		}
	}
}

func TestKeyBindingsAction(t *testing.T) {
	tests := []struct {
		kp     KeyPress
		retNil bool
		ck     string
	}{
		{
			KeyPress{Key: 'i'},
			false,
			"test3",
		},
		{
			KeyPress{Key: 'p'},
			false,
			"t2",
		},
		{
			KeyPress{Key: 'i', Ctrl: true},
			true,
			"",
		},
		{
			KeyPress{Key: 'c'},
			false,
			"t5",
		},
	}

	if d, err := ioutil.ReadFile("testdata/Default.sublime-keymap"); err != nil {
		t.Fatal(err)
	} else {
		var (
			bindings KeyBindings
			p        HasKeyBindings
		)
		loaders.LoadJSON(d, &bindings)

		if d, err = ioutil.ReadFile("testdata/test.sublime-keymap"); err != nil {
			t.Fatal(err)
		}
		loaders.LoadJSON(d, p.KeyBindings())
		bindings.SetParent(&p)

		for i, test := range tests {
			qc := func(key string, operator util.Op, operand interface{}, match_all bool) bool {
				return key == test.ck
			}
			b := bindings.Filter(test.kp)
			if a := b.Action(qc); test.retNil {
				if a != nil {
					t.Errorf("Test %d: Expected action to be nil but got %v", i, a)
				}
			} else if a.Context[0].Key != test.ck {
				t.Errorf("Test %d: Expected %s, but got %s", i, test.ck, a.Context[0].Key)
			}
		}
	}
}

func TestSeqIndex(t *testing.T) {
	var bd KeyBindings
	bd.seqIndex = 3
	if bd.SeqIndex() != 3 {
		t.Errorf("Expected SeqIndex %d, but got %d", 3, bd.SeqIndex())
	}
}

func TestKeyBindingsString(t *testing.T) {
	fn := "testdata/test.sublime-keymap"
	var bd KeyBindings

	d, err := ioutil.ReadFile(fn)
	if err != nil {
		t.Fatalf("Couldn't read %s: %s", fn, err)
	}
	if err = loaders.LoadJSON(d, &bd); err != nil {
		t.Fatalf("Error loading json: %s", err)
	}

	expected :=
		`&{Keys:[p] Command:t2 Args:map[] Context:[{rawKeyContext:{Key:t2 Operator:0 Operand:true MatchAll:false}}] priority:1}
&{Keys:[ctrl+d ctrl+k] Command:t1 Args:map[] Context:[{rawKeyContext:{Key:t1 Operator:0 Operand:true MatchAll:false}}] priority:0}
`
	if bd.String() != expected {
		t.Errorf("Expected String %s, but got %s", expected, bd.String())
	}
}
