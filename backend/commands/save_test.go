// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/limetext/lime/backend"
	"io/ioutil"
	"os"
	"testing"
)

var testfile string = "testdata/save_test.txt"

func TestSave(t *testing.T) {
	hold, err := ioutil.ReadFile(testfile)
	if err != nil {
		t.Fatalf("Couldn't read test file %s", testfile)
	}
	if err := ioutil.WriteFile(testfile, []byte("Before text"), 0644); err != nil {
		t.Fatalf("Couldn't write test file %s", testfile)
	}

	tests := []struct {
		text   string
		expect string
	}{
		{
			" ab\ncd",
			"Before text ab\ncd",
		},
		{
			"\n",
			"Before text\n",
		},
	}

	ed := GetEditor()
	w := ed.NewWindow()
	defer w.Close()

	for i, test := range tests {
		err := ioutil.WriteFile(testfile, []byte("Before text"), 0644)
		if err != nil {
			t.Fatal("Could not write to test file")
		}

		v := w.OpenFile(testfile, 0)
		e := v.BeginEdit()
		v.Insert(e, v.Buffer().Size(), test.text)
		v.EndEdit(e)

		ed.CommandHandler().RunTextCommand(v, "save", nil)
		if v.IsDirty() {
			t.Errorf("Test %d: Expected the view to be clean, but it wasn't", i)
		}

		if data, _ := ioutil.ReadFile(testfile); test.expect != string(data) {
			t.Errorf("Test %d: Expected %s, but got %s", i, test.expect, string(data))
		}

		v.Close()

		if err := ioutil.WriteFile(testfile, hold, 0644); err != nil {
			t.Fatalf("Couldn't write back test file %s", testfile)
		}
	}
}

func TestSaveAs(t *testing.T) {
	hold, err := ioutil.ReadFile(testfile)
	if err != nil {
		t.Fatalf("Couldn't read test file %s", testfile)
	}
	if err := ioutil.WriteFile(testfile, []byte(""), 0644); err != nil {
		t.Fatalf("Couldn't write test file %s", testfile)
	}

	ed := GetEditor()
	w := ed.NewWindow()
	// defer w.Close()

	v := w.OpenFile(testfile, 0)
	e := v.BeginEdit()
	v.Insert(e, 0, "Testing save_as command")
	v.BeginEdit()

	name := "testdata/save_as_test.txt"

	ed.CommandHandler().RunTextCommand(v, "save_as", Args{"name": name})

	if v.IsDirty() {
		t.Error("Expected the view to be clean, but it wasn't")
	}

	if _, err := os.Stat(name); os.IsNotExist(err) {
		t.Errorf("The new test file %s wasn't created", name)
	}
	if data, _ := ioutil.ReadFile(name); "Testing save_as command" != string(data) {
		t.Errorf("Expected %s, but got %s", "Testing save_as command", string(data))
	}

	// v.Close()

	if err := os.Remove(name); err != nil {
		t.Errorf("Couldn't remove test file %s", name)
	}
	if err := ioutil.WriteFile(testfile, hold, 0644); err != nil {
		t.Fatalf("Couldn't write back test file %s", testfile)
	}
}

func TestSaveAll(t *testing.T) {
	var err error
	holds := make(map[int][]byte)
	views := make(map[int]View)
	files := []struct {
		file   string
		expect string
	}{
		{
			"testdata/save_all_test.txt",
			"Testing save all 1",
		},
		{
			"testdata/save_another_all_test.txt",
			"Testing save all 2",
		},
	}

	ed := GetEditor()
	fe := ed.Frontend()
	if dfe, ok := fe.(*DummyFrontend); ok {
		// Make it *not* reload the file
		dfe.SetDefaultAction(false)
	}

	w := ed.NewWindow()
	// defer w.Close()

	for i, f := range files {
		holds[i], err = ioutil.ReadFile(f.file)
		if err != nil {
			t.Fatalf("Test %d: Couldn't read file %s", i, f.file)
		}
		if err := ioutil.WriteFile(f.file, []byte(""), 0644); err != nil {
			t.Fatalf("Test %d: Couldn't write test file %s", i, f.file)
		}

		v := w.OpenFile(f.file, 0)
		views[i] = *v

		e := v.BeginEdit()
		v.Insert(e, 0, f.expect)
		v.EndEdit(e)
	}

	if err := ed.CommandHandler().RunWindowCommand(w, "save_all", nil); err != nil {
		t.Errorf("failed to run save_all: %s", err)
	}

	for i, f := range files {
		if data, err := ioutil.ReadFile(f.file); err != nil {
			t.Errorf("failed to read in file: %s", err)
		} else if s := string(data); s != f.expect {
			t.Errorf("Test %d: Expected to get `%s`, but got `%s`", i, f.expect, s)
		}
	}

	for i, f := range files {
		// v := views[i]
		// v.SetScratch(true)
		// v.Close()
		ioutil.WriteFile(f.file, holds[i], 0644)
	}
}
