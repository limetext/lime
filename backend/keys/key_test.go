// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package keys

import (
	"testing"
)

func TestKeyString(t *testing.T) {
	k := Key('i')
	if k.String() != "i" {
		t.Errorf("Expected %q but got %q", "i", k.String())
	}

	k = Key(Home)
	if k.String() != "home" {
		t.Errorf("Expected %q but got %q", "home", k.String())
	}
}
