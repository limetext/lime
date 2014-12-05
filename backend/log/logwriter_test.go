// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package log_test

import (
	"fmt"
	"github.com/limetext/lime/backend/log"
	"github.com/limetext/log4go"
	"testing"
)

func TestNewLogWriter(t *testing.T) {
	l := log.NewLogWriter(func(str string) {})
	if l == nil {
		t.Error("NewLogWriter produced a nil")
	}
	l.Close()
}

func TestLogWrite(t *testing.T) {
	l := log.NewLogWriter(func(str string) { fmt.Print(str) })
	logRecord := &log4go.LogRecord{}
	l.LogWrite(logRecord) // Void function. Testing for coverage.
	l.Close()
}
