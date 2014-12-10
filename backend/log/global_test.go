// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package log

import (
	"testing"
)

func TestGlobalFunctions(t *testing.T) {
	AddFilter("add filter", FINE, testlogger(func(str string) {}))
	Finest("testing finest")
	Fine("testing fine")
	Debug("testing debug")
	Trace("testing trace")
	Warn("testing warn")
	Error("testing error")
	Critical("testing critical")
	Logf(FINE, "testing logf")
	Close("testing close")
}
