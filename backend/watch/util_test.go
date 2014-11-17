// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package watch

import (
	"reflect"
	"testing"
)

func TestExist(t *testing.T) {
	test := struct {
		array []string
		elms  []string
		exps  []bool
	}{
		[]string{"a", "b", "c", "d"},
		[]string{"a", "t", "A"},
		[]bool{true, false, false},
	}
	for i, exp := range test.exps {
		if exist(test.array, test.elms[i]) != exp {
			t.Errorf("Expected in %v exist result of element %s be %v, but got %v", test.array, test.elms[i], exp, exist(test.array, test.elms[i]))
		}
	}
}

func TestRemove(t *testing.T) {
	tests := []struct {
		slice  []string
		remove string
		exp    []string
	}{
		{
			[]string{"a", "b", "c"},
			"a",
			[]string{"c", "b"},
		},
		{
			[]string{"a", "b", "c"},
			"k",
			[]string{"a", "b", "c"},
		},
	}
	for i, test := range tests {
		if exp := remove(test.slice, test.remove); !reflect.DeepEqual(exp, test.exp) {
			t.Errorf("Test %d: Expected %v be equal to %v", i, exp, test.exp)
		}
	}
}
