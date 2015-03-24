// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package log

import (
	"testing"
)

/*
Since the logger only calls a function to an outside package there is no way to determine if the functions are behaving
We can only know if the functions throw an error
This should up the test coverage
*/
func TestAddFilter(t *testing.T) {
	AddFilter("add filter", FINE, testlogger(func(str string) {}))
}
func TestFinest(t *testing.T) {
	Finest("testing finest")
}
func TestFine(t *testing.T) {
	Fine("testing fine")
}
func TestDebug(t *testing.T) {
	Debug("testing debug")
}
func TestTrace(t *testing.T) {
	Trace("testing trace")
}
func TestInfo(t *testing.T) {
	Warn("testing warn")
}
func TestWarn(t *testing.T) {
	Warn("testing warn")
}
func TestError(t *testing.T) {
	Error("testing error")
}
func TestErrorf(t *testing.T) {
	Errorf("testing %s", "errorf")
}
func TestCritical(t *testing.T) {
	Critical("testing critical")
}

//TestLogf is already defined in logger_test
func TestGlobalLogf(t *testing.T) {
	Logf(FINE, "testing %s", "logf")
}

//TestClose is already defined in logger_test
func TestGlobalClose(t *testing.T) {
	Close()
}
