// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/limetext/lime/backend"
	"io/ioutil"
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

func TestOpenDir(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	file1, err := ioutil.TempFile(dir, "")
	if err != nil {
		t.Fatal(err)
	}
	file2, err := ioutil.TempFile(dir, "")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		path string
		exp  []string
	}{
		{
			dir,
			[]string{file1.Name(), file2.Name()},
		},
	}

	ed := GetEditor()
	for i, test := range tests {
		w := ed.NewWindow()
		err := ed.CommandHandler().RunWindowCommand(w, "open_dir",
			Args{"path": test.path})
		if err != nil {
			t.Fatal(err)
		}
		vs := w.Views()
		for _, exp := range test.exp {
			found := false
			for _, v := range vs {
				if v.Buffer().FileName() == exp {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Test %d: Expected to find %s in views", i, exp)
			}
		}
		w.Close()
	}

	w := ed.NewWindow()
	defer w.Close()

	if err := ed.CommandHandler().RunWindowCommand(w, "open_dir",
		Args{"path": file1.Name()}); err == nil {
		t.Error("Expected error on running open_dir command on file")
	}
}
