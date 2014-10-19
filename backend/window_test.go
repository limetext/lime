// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package backend

import (
	"testing"
)

func TestWindowNewFile(t *testing.T) {
	w := GetEditor().NewWindow()
	defer w.Close()

	v := w.NewFile()
	defer func() {
		v.SetScratch(true)
		v.Close()
	}()

	if len(w.Views()) != 1 {
		t.Errorf("Expected 1 view, but got %d", len(w.Views()))
	}
}

func TestWindowRemove(t *testing.T) {
	w := GetEditor().NewWindow()
	defer w.Close()

	v0 := w.NewFile()
	defer v0.Close()

	v1 := w.NewFile()
	defer v1.Close()

	v2 := w.NewFile()
	defer v2.Close()

	l := len(w.Views())

	w.remove(v1)

	if len(w.Views()) != l-1 {
		t.Errorf("Expected %d open views, but got %d", l-1, len(w.Views()))
	}
}

func TestWindowClose(t *testing.T) {
	ed := GetEditor()
	l := len(ed.Windows())
	w := ed.NewWindow()

	for _, v := range w.Views() {
		v.SetScratch(true)
		v.Close()
	}

	w.Close()

	if len(ed.Windows()) != l {
		t.Errorf("Expected window to close, but we have %d still open", len(ed.Windows()))
	}
}
