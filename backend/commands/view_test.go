// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/limetext/lime/backend"
	"testing"
)

func TestNextView(t *testing.T) {
	ed := GetEditor()
	w := ed.NewWindow()

	v0 := w.NewFile()
	v1 := w.NewFile()
	v2 := w.NewFile()
	v3 := w.NewFile()

	w.SetActiveView(v1)

	ed.CommandHandler().RunWindowCommand(w, "next_view", nil)

	av := w.ActiveView()
	if av != v2 {
		t.Error("Expected to get v2, but didn't")
	}

	w.SetActiveView(v3)

	ed.CommandHandler().RunWindowCommand(w, "next_view", nil)

	av = w.ActiveView()
	if av != v0 {
		t.Error("Expected to get v0, but didn't")
	}
}

func TestPrevView(t *testing.T) {
	ed := GetEditor()
	w := ed.NewWindow()

	v0 := w.NewFile()
	v1 := w.NewFile()
	v2 := w.NewFile()
	v3 := w.NewFile()

	w.SetActiveView(v2)

	ed.CommandHandler().RunWindowCommand(w, "prev_view", nil)

	av := w.ActiveView()
	if av != v1 {
		t.Error("Expected to get v1, but didn't")
	}

	w.SetActiveView(v0)

	ed.CommandHandler().RunWindowCommand(w, "prev_view", nil)

	av = w.ActiveView()
	if av != v3 {
		t.Error("Expected to get v3, but didn't")
	}
}
