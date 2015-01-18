// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package keys

import "testing"

func TestKeyContextUnmarshalError(t *testing.T) {
	var context KeyContext
	if err := context.UnmarshalJSON([]byte(``)); err == nil {
		t.Errorf("Expected error on loading empty string")
	}
}
