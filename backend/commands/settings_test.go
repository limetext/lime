// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/limetext/lime/backend"
	"testing"
)

func TestToggleSetting(t *testing.T) {
	ed := GetEditor()
	w := ed.NewWindow()
	defer w.Close()

	v := w.NewFile()
	defer func() {
		v.SetScratch(true)
		v.Close()
	}()

	v.Settings().Set("duck", true)
	v.Settings().Set("witch", false)
	v.Settings().Set("favorite_color", "blue")

	// non-existant -> true
	ed.CommandHandler().RunTextCommand(v, "toggle_setting", Args{"setting": "rabbit"})
	if val, ok := v.Settings().Get("rabbit").(bool); !ok || !val {
		t.Errorf("Toggling an non-existant setting should make it true")
	}

	// non-bool -> true
	ed.CommandHandler().RunTextCommand(v, "toggle_setting", Args{"setting": "favorite_color"})
	if val, ok := v.Settings().Get("favorite_color").(bool); !ok || !val {
		t.Errorf("Toggling an non-bool setting should make it true")
	}

	// bool: true -> false
	ed.CommandHandler().RunTextCommand(v, "toggle_setting", Args{"setting": "duck"})
	if val, ok := v.Settings().Get("duck").(bool); !ok || val {
		t.Errorf("Setting should be toggled from true to false")
	}

	// bool: false -> true
	ed.CommandHandler().RunTextCommand(v, "toggle_setting", Args{"setting": "witch"})
	if val, ok := v.Settings().Get("witch").(bool); !ok || !val {
		t.Errorf("Setting should be toggled from false to true")
	}
}

func TestSetSetting(t *testing.T) {
	ed := GetEditor()
	w := ed.NewWindow()
	defer w.Close()

	v := w.NewFile()
	defer func() {
		v.SetScratch(true)
		v.Close()
	}()

	v.Settings().Set("favorite_color", "blue")

	exp := "red"
	ed.CommandHandler().RunTextCommand(v, "set_setting", Args{"setting": "favorite_color", "value": exp})
	val := v.Settings().Get("favorite_color")
	if s, ok := val.(string); !ok || s != exp {
		t.Errorf("Expecting setting value to be %#v, was %#v", exp, val)
	}
}
