// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/limetext/lime/backend"
	"testing"
)

func TestNewFile(t *testing.T) {
	ed := GetEditor()
	w := ed.ActiveWindow()
	l := len(w.Views())

	ed.CommandHandler().RunWindowCommand(w, "new_file", nil)

	if len(w.Views()) != l+1 {
		t.Errorf("Expected %d views, but got %d", l+1, len(w.Views()))
	}
}
