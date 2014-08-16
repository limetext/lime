// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package backend

import (
	"testing"
)

func TestPascalCaseToSnakeCase(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{
			"TestString",
			"test_string",
		},
		{
			"Teststring",
			"teststring",
		},
	}

	for i, test := range tests {
		out := PascalCaseToSnakeCase(test.in)

		if out != test.out {
			t.Errorf("Test %d: Expected %s, but got %s", i, test.out, out)
		}
	}
}
