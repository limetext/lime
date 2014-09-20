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

	if d, err := ioutil.ReadFile("testdata/Default.sublime-keymap"); err == nil {
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
	}

	if d, err := ioutil.ReadFile("testdata/Default.sublime-keymap"); err == nil {
		var bindings KeyBindings
		loaders.LoadJSON(d, &bindings)

		var qc func(key string, operator util.Op, operand interface{}, match_all bool) bool

		for i, test := range tests {
			qc = func(key string, operator util.Op, operand interface{}, match_all bool) bool {
				return key == test.ck
			}
			b := bindings.Filter(test.kp)
			if a := b.Action(qc); a.Context[0].Key != test.ck {
				t.Errorf("Test %d: Expected %s, but got %s", i, test.ck, a.Context[0].Key)
			}
		}
	}
}
