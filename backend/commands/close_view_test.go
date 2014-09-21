// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/limetext/lime/backend"
	"testing"
)

func TestCloseView(t *testing.T) {
	ed := GetEditor()

	w := ed.NewWindow()
	defer w.Close()

	l := len(w.Views())
	testPath := "open_file_test.go"
	ed.CommandHandler().RunWindowCommand(w, "open_file", Args{"path": testPath})
	ed.CommandHandler().RunWindowCommand(w, "close_view", nil)

	if len(w.Views()) != l {
		t.Errorf("Expected %d view, but got %d", l, len(w.Views()))
	}
}
