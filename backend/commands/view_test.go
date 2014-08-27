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
	w := ed.ActiveWindow()
	i := w.ActiveViewIndex()

	e := 0

	if i != 7 {
		t.Errorf("Expected to start with view index 7, but got %d", i)
	}

	ed.CommandHandler().RunWindowCommand(w, "next_view", nil)
	if w.ActiveViewIndex() != e {
		t.Errorf("Expected view index %d, but got %d", e, w.ActiveViewIndex())
	}

	e = 1

	ed.CommandHandler().RunWindowCommand(w, "next_view", nil)
	if w.ActiveViewIndex() != e {
		t.Errorf("Expected view index %d, but got %d", e, w.ActiveViewIndex())
	}
}

func TestPrevView(t *testing.T) {
	ed := GetEditor()
	w := ed.ActiveWindow()
	i := w.ActiveViewIndex()

	w.SetActiveView(w.Views()[1])

	e := 0

	if i != 1 {
		t.Errorf("Expected to start with view index 1, but got %d", i)
	}

	ed.CommandHandler().RunWindowCommand(w, "prev_view", nil)
	if w.ActiveViewIndex() != e {
		t.Errorf("Expected view index %d, but got %d", e, w.ActiveViewIndex())
	}

	e = 7

	ed.CommandHandler().RunWindowCommand(w, "prev_view", nil)
	if w.ActiveViewIndex() != e {
		t.Errorf("Expected view index %d, but got %d", e, w.ActiveViewIndex())
	}
}
