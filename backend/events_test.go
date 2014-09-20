// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package backend

import (
	. "github.com/quarnster/util/text"
	"os"
	"reflect"
	"testing"
)

func TestOnSelectionModified(t *testing.T) {
	var res *RegionSet

	callCount := 0

	w := GetEditor().NewWindow()
	defer w.Close()

	v := w.NewFile()
	defer v.Close()

	OnSelectionModified.Add(func(v *View) {
		res = v.Sel()
		callCount++
	})
	edit := v.BeginEdit()
	v.Insert(edit, 0, "abcd")
	v.EndEdit(edit)
	if callCount != 1 {
		t.Fatalf("%d != 1", callCount)
	}
	if !reflect.DeepEqual(res.Regions(), []Region{{4, 4}}) {
		t.Errorf("%v", res.Regions())
	}

	edit = v.BeginEdit()
	v.Sel().Adjust(4, -1)
	v.EndEdit(edit)

	if callCount != 2 {
		t.Fatalf("%d != 2", callCount)
	}
	if !reflect.DeepEqual(res.Regions(), []Region{{3, 3}}) {
		t.Errorf("%v", res.Regions())
	}

	edit = v.BeginEdit()
	v.EndEdit(edit)

	if callCount != 2 {
		t.Fatalf("%d != 2", callCount)
	}
	if !reflect.DeepEqual(res.Regions(), []Region{{3, 3}}) {
		t.Errorf("%v", res.Regions())
	}
}

func TestOnPreSave(t *testing.T) {
	testfile := "testdata/test_event.txt"
	callCount := 0

	w := GetEditor().NewWindow()
	defer w.Close()

	v := w.NewFile()
	defer v.Close()

	OnPreSave.Add(func(v *View) {
		callCount++
	})
	edit := v.BeginEdit()
	v.Insert(edit, 0, "abcd")
	v.EndEdit(edit)
	if err := v.SaveAs(testfile); err != nil {
		t.Fatal("Could not save the view")
	}
	if callCount != 1 {
		t.Fatalf("%d != 1", callCount)
	}
	v.Buffer().SetFileName(testfile)
	if err := v.Save(); err != nil {
		t.Fatalf("Could not save the view %s", err)
	}
	if callCount != 2 {
		t.Fatalf("%d != 2", callCount)
	}
	if err := os.Remove(testfile); err != nil {
		t.Errorf("Couldn't remove test file %s", testfile)
	}
}

func TestOnPostSave(t *testing.T) {
	testfile := "testdata/test_event.txt"
	callCount := 0

	w := GetEditor().NewWindow()
	defer w.Close()

	v := w.NewFile()
	defer v.Close()

	OnPostSave.Add(func(v *View) {
		callCount++
	})
	edit := v.BeginEdit()
	v.Insert(edit, 0, "abcd")
	v.EndEdit(edit)
	if err := v.SaveAs(testfile); err != nil {
		t.Fatal("Could not save the view")
	}
	if callCount != 1 {
		t.Fatalf("%d != 1", callCount)
	}
	v.Buffer().SetFileName(testfile)
	if err := v.Save(); err != nil {
		t.Fatalf("Could not save the view: %s", err)
	}
	if callCount != 2 {
		t.Fatalf("%d != 2", callCount)
	}
	if err := os.Remove(testfile); err != nil {
		t.Errorf("Couldn't remove test file %s", testfile)
	}
}
