// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/limetext/lime/backend"
	. "github.com/quarnster/util/text"
	"io/ioutil"
	"os"
	"testing"
)

var testfile string = "../testdata/save_test.txt"

func TestSave(t *testing.T) {
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
		if data := v.Buffer().Substr(Region{0, v.Buffer().Size()}); test.expect != data {
			t.Errorf("Test %d: Expected %s, but got %s", i, test.expect, data)
		}
	}
}

func TestSaveAs(t *testing.T) {
	ed := GetEditor()
	w := ed.NewWindow()
	v := w.OpenFile(testfile, 0)
	name := "../testdata/save_as_test.txt"

	ed.CommandHandler().RunTextCommand(v, "save_as", Args{"name": name})
	if _, err := os.Stat(name); os.IsNotExist(err) {
		t.Errorf("The new test file %s wasn't created", name)
	}
}
