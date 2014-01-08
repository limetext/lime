// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package backend

import (
	"io/ioutil"
	"lime/backend/loaders"
	"testing"
)

func TestLoadKeyBindings(t *testing.T) {
	tests := []string{
		"loaders/json/testdata/Default (OSX).sublime-keymap",
	}
	for i, fn := range tests {
		if d, err := ioutil.ReadFile(fn); err != nil {
			if i == 0 {
				t.Errorf("Couldn't load file %s: %s", fn, err)
			}
		} else {
			var bindings KeyBindings
			if err := loaders.LoadJSON(d, &bindings); err != nil {
				t.Error(err)
			}
		}
	}
}

func TestKeyFilter(t *testing.T) {
	fn := "loaders/json/testdata/Default (OSX).sublime-keymap"
	if d, err := ioutil.ReadFile(fn); err != nil {
		t.Errorf("Couldn't load file %s: %s", fn, err)
	} else {
		var bindings KeyBindings
		if err := loaders.LoadJSON(d, &bindings); err != nil {
			t.Error(err)
		}

		if b2 := bindings.Filter(KeyPress{Key: 'j', Ctrl: true}); b2.Len() != 3 {
			t.Errorf("Not of the expected length: %d, %s", 3, b2)
		} else if b3 := b2.Filter(KeyPress{Key: 's'}); b3.Len() != 1 {
			t.Errorf("Not of the expected length: %d, %s", 1, b3)
		}
	}
}

func TestKeyFilter2(t *testing.T) {
	ed := GetEditor()
	w := ed.NewWindow()
	v := w.NewFile()
	enable := "test1"
	OnQueryContext.Add(func(v *View, key string, operator Op, operand interface{}, match_all bool) QueryContextReturn {
		if key == enable {
			return True
		}
		return Unknown
	})
	fn := "testdata/Default.sublime-keymap"
	if d, err := ioutil.ReadFile(fn); err != nil {
		t.Errorf("Couldn't load file %s: %s", fn, err)
	} else {
		var bindings KeyBindings
		if err := loaders.LoadJSON(d, &bindings); err != nil {
			t.Error(err)
		}
		b2 := bindings.Filter(KeyPress{Key: 'i'})
		a := b2.Action(v)
		if a.Context[0].Key != enable {
			t.Error(b2, a)
		}
	}
}

func TestVintageous(t *testing.T) {
	fn := "testdata/Vintageous.sublime-keymap"
	ed := GetEditor()
	w := ed.NewWindow()
	v := w.NewFile()
	v.Settings().Set("command_mode", true)

	OnQueryContext.Add(func(v *View, key string, op Op, operand interface{}, match_all bool) QueryContextReturn {
		if key == "vi_has_action" {
			return True
		}
		return Unknown
	})

	if d, err := ioutil.ReadFile(fn); err != nil {
		t.Errorf("Couldn't load file %s: %s", fn, err)
	} else {
		var bindings KeyBindings
		if err := loaders.LoadJSON(d, &bindings); err != nil {
			t.Error(err)
		}

		b2 := bindings.Filter(KeyPress{Key: 'g'})
		if a := b2.Action(v); a == nil || a.Command != "set_action" {
			t.Error(a)
		}
		b2 = b2.Filter(KeyPress{Key: 'g'})
		if a := b2.Action(v); a == nil || a.Command != "set_motion" {
			t.Error(a)
		}
	}
}
