// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/limetext/lime/backend"
	"path/filepath"
	"testing"
)

func TestNewFile(t *testing.T) {
	ed := GetEditor()

	w := ed.NewWindow()
	defer w.Close()

	l := len(w.Views())

	ed.CommandHandler().RunWindowCommand(w, "new_file", nil)

	if len(w.Views()) != l+1 {
		t.Errorf("Expected %d views, but got %d", l+1, len(w.Views()))
	}

	for _, v := range w.Views() {
		v.SetScratch(true)
		v.Close()
	}
}

func TestOpenFile(t *testing.T) {
	ed := GetEditor()

	w := ed.NewWindow()
	defer w.Close()

	l := len(w.Views())

	testPath := "open_file_test.go"
	ed.CommandHandler().RunWindowCommand(w, "open_file", Args{"path": testPath})

	if len(w.Views()) != l+1 {
		t.Errorf("Expected %d views, but got %d", l+1, len(w.Views()))
	}
	exp, err := filepath.Abs(testPath)
	if err != nil {
		exp = testPath
	}
	if w.Views()[l].Buffer().FileName() != exp {
		t.Errorf("Expected %s as FileName, but got %s", testPath, w.Views()[l].Buffer().FileName())
	}

	for _, v := range w.Views() {
		v.SetScratch(true)
		v.Close()
	}
}
