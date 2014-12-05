// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package log_test

import (
	"github.com/limetext/lime/backend/log"
	"testing"
)

func TestGlobalFunctions(t *testing.T) {
	log.AddFilter("add filter", log.FINE, testlogger(func(str string) {}))
	log.Finest("testing finest")
	log.Fine("testing fine")
	log.Debug("testing debug")
	log.Trace("testing trace")
	log.Warn("testing warn")
	log.Error("testing error")
	log.Critical("testing critical")
	log.Logf(log.FINE, "testing logf")
	log.Close("testing close")
}
