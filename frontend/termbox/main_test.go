// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package main

import (
	"testing"
)

func TestPadLineNumber(t *testing.T) {
	var testPadData = []struct {
		line     int
		total    int
		expected []rune
	}{
		{1, 3, []rune{' ', ' ', '1', ' '}},
		{10, 3, []rune{' ', '1', '0', ' '}},
		{100, 3, []rune{'1', '0', '0', ' '}},
	}

	for _, p := range testPadData {
		runes := intToRunes(p.line)
		padded := padLineRunes(runes, p.total)

		for i, r := range p.expected {
			if r != padded[i] {
				t.Error("Expected runes to be padded")
			}
		}
	}
}
