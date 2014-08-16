// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package backend

import (
	"testing"
)

func TestNewFile(t *testing.T) {
	ed := GetEditor()
	w := ed.NewWindow()
	w.NewFile()

	if len(w.Views()) != 1 {
		t.Errorf("Expected 1 view, but got %d", len(w.Views()))
	}
}

func TestClose(t *testing.T) {
	ed := GetEditor()
	l := len(ed.Windows())
	w := ed.NewWindow()

	w.Close()

	if len(ed.Windows()) != l {
		t.Errorf("Expected window to close, but we have %d still open", len(ed.Windows()))
	}
}
