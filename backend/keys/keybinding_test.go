// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package keys

import (
	"github.com/limetext/lime/backend/loaders"
	"github.com/limetext/lime/backend/util"
	"io/ioutil"
	"testing"
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
	var bd, p KeyBindings

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
	if err = loaders.LoadJSON(d, &p); err != nil {
		t.Fatalf("Error loading json: %s", err)
	}

	p.keyOff = 10
	bd.SetParent(&p)
	if bd.keyOff != p.keyOff {
		t.Fatalf("Expected parent and child keyOff be equal %d != %d", p.keyOff, bd.keyOff)
	}

	ret := bd.Filter(KeyPress{Key: 'd', Ctrl: true})
	if ret.Len() != 1 {
		t.Fatalf("Expected ret keyBindings len %d, but got %d", 1, ret.Len())
	}
	if ret.parent.Len() != 1 {
		t.Fatalf("Expected ret parent keyBindings len %d, but got %d", 1, ret.parent.Len())
	}
	if cmd := ret.Bindings[0].Command; cmd != "test4" {
		t.Errorf("Expected Command %s, but got %s", "test4", cmd)
	}
	if cmd := ret.parent.Bindings[0].Command; cmd != "t1" {
		t.Errorf("Expected Command %s, but got %s", "t1", cmd)
	}
}

func TestParent(t *testing.T) {
	fn := "testdata/Default.sublime-keymap"
	fnp := "testdata/test.sublime-keymap"
	var bd, p KeyBindings

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
	if err = loaders.LoadJSON(d, &p); err != nil {
		t.Fatalf("Error loading json: %s", err)
	}

	bd.SetParent(&p)

	if cmd := bd.Parent().Bindings[0].Command; cmd != "t2" {
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
		kp KeyPress
		ck string
	}{
		{
			KeyPress{Key: 'i'},
			"test3",
		},
		{
			KeyPress{Key: 'p'},
			"t2",
		},
	}

	if d, err := ioutil.ReadFile("testdata/Default.sublime-keymap"); err != nil {
		t.Fatal(err)
	} else {
		var bindings, p KeyBindings
		loaders.LoadJSON(d, &bindings)

		if d, err = ioutil.ReadFile("testdata/test.sublime-keymap"); err != nil {
			t.Fatal(err)
		}
		loaders.LoadJSON(d, &p)
		bindings.SetParent(&p)

		for i, test := range tests {
			qc := func(key string, operator util.Op, operand interface{}, match_all bool) bool {
				return key == test.ck
			}
			b := bindings.Filter(test.kp)
			if a := b.Action(qc); a.Context[0].Key != test.ck {
				t.Errorf("Test %d: Expected %s, but got %s", i, test.ck, a.Context[0].Key)
			}
		}
	}
}
