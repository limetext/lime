// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package keys

import (
	"testing"
)

func TestKeyPressFix(t *testing.T) {
	k := KeyPress{'A', false, false, false, false}
	k.fix()
	if k.Key != 'a' {
		t.Errorf("Expected the key to be %q, but it was %q", 'a', k.Key)
	}
	if !k.Shift {
		t.Error("Expected the shift modifier to be active, but it wasn't")
	}
}
