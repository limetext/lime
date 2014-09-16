// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/limetext/lime/backend"
	"testing"
)

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
	if w.Views()[l].Buffer().FileName() != testPath {
		t.Errorf("Expected %s as FileName, but got %s", testPath, w.Views()[l].Buffer().FileName())
	}
}
