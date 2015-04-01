// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package keys

import (
	"testing"
)

func TestKeyPressIndex(t *testing.T) {
	tests := []struct {
		kp  KeyPress
		exp int
	}{
		{
			KeyPress{Key: 'a', Shift: false, Super: false, Alt: false, Ctrl: false},
			int('a'),
		},
		{
			KeyPress{Key: 'a', Shift: true, Super: false, Alt: false, Ctrl: false},
			int('a') + shift,
		},
		{
			KeyPress{Key: 'a', Shift: true, Super: true, Alt: false, Ctrl: false},
			int('a') + shift + super,
		},
		{
			KeyPress{Key: 'a', Shift: true, Super: true, Alt: true, Ctrl: false},
			int('a') + shift + super + alt,
		},
		{
			KeyPress{Key: 'a', Shift: true, Super: true, Alt: true, Ctrl: true},
			int('a') + shift + super + alt + ctrl,
		},
	}

	for i, test := range tests {
		if test.kp.Index() != test.exp {
			t.Errorf("Test %d: Expected %d, but got %d", i, test.exp, test.kp.Index())
		}
	}
}

func TestKeyPressIsCharacter(t *testing.T) {
	tests := []struct {
		kp  KeyPress
		exp bool
	}{
		{
			KeyPress{Key: 'a', Shift: false, Super: false, Alt: false, Ctrl: false},
			true,
		},
		{
			KeyPress{Key: 'a', Shift: true, Super: false, Alt: false, Ctrl: false},
			true,
		},
		{
			KeyPress{Key: 'a', Shift: false, Super: true, Alt: false, Ctrl: false},
			false,
		},
		{
			KeyPress{Key: 'a', Shift: false, Super: false, Alt: false, Ctrl: true},
			false,
		},
		// {
		// 	KeyPress{Key: F1, Shift: false, Super: false, Alt: false, Ctrl: false},
		// 	false,
		// },
	}

	for i, test := range tests {
		if test.kp.IsCharacter() != test.exp {
			t.Errorf("Test %d: Expected %v, but got %v", i, test.exp, test.kp.IsCharacter())
		}
	}
}

func TestKeyPressFix(t *testing.T) {
	k := KeyPress{"A", 'A', false, false, false, false}
	k.fix()
	if k.Key != 'a' {
		t.Errorf("Expected the key to be %q, but it was %q", 'a', k.Key)
	}
	if !k.Shift {
		t.Error("Expected the shift modifier to be active, but it wasn't")
	}
}

func TestKeyPressUnmarshalJSON(t *testing.T) {
	var k KeyPress
	d := `"super+ctrl+alt+shift+f1+λλλ"`
	err := k.UnmarshalJSON([]byte(d))
	if err != nil {
		t.Error(err)
	}
}

func TestKeyPressString(t *testing.T) {
	k1 := KeyPress{"a", 'a', true, true, false, false}
	if k1.String() != "super+shift+a" {
		t.Errorf("Expected %q, but got %q", "super+shift+a", k1.String())
	}

	k2 := KeyPress{"b", 'b', true, false, true, true}
	if k2.String() != "ctrl+alt+shift+b" {
		t.Errorf("Expected %q, but got %q", "ctrl+alt+shift+b", k2.String())
	}
}
