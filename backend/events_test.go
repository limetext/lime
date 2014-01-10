// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package backend

import (
	. "github.com/quarnster/util/text"
	"reflect"
	"testing"
)

func TestOnSelectionModified(t *testing.T) {
	var (
		w         Window
		v         = w.NewFile()
		res       *RegionSet
		callCount = 0
	)
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
