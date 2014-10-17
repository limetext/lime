// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package backend

import (
	. "github.com/limetext/text"
	"os"
	"testing"
	"time"
)

func TestOnFileChange(t *testing.T) {
	ed := GetEditor()
	fe := ed.Frontend()
	if dfe, ok := fe.(*DummyFrontend); ok {
		// Make it trigger a reload
		dfe.SetDefaultAction(true)
	}

	filename := "testdata/test.txt"
	modified := make(chan bool)
	timer := time.NewTimer(time.Second * 2)

	window := ed.NewWindow()
	defer window.Close()

	file, err := os.Create(filename)
	if err != nil {
		t.Fatalf("Could not open file %v", err)
	}
	defer func() {
		file.Close()
		os.Remove(filename)
	}()

	if _, err := file.WriteString("foo bar baz"); err != nil {
		t.Fatalf("Could not write to fixutre file due to %v", err)
	}

	view := window.OpenFile(filename, 0)
	defer func() {
		view.SetScratch(true)
		view.Close()
	}()

	contents := view.Buffer().Substr(Region{0, 11})
	if contents != "foo bar baz" {
		t.Fatalf("Expected substring 'foo bar baz', got: %v", contents)
	}

	OnModified.Add(func(v *View) {
		if v == view {
			contents := view.Buffer().Substr(Region{0, 17})
			if contents != "" {
				if contents != "foo bar baz scuzz" {
					t.Fatalf("Expected substring 'foo bar baz scuzz', got: %v", contents)
				} else {
					modified <- true
				}
			}
		}
	})

	if _, err := file.WriteString(" scuzz"); err != nil {
		t.Fatalf("Could not write to fixture file due to %v", err)
	}

	select {
	case <-modified:
	case <-timer.C:
		t.Fatal("Timeout reached")
	}
}
