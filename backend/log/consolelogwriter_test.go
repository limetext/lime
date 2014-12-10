// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package log

import (
	"github.com/limetext/log4go"
	"testing"
)

func TestNewConsoleLogWriter(t *testing.T) {
	l := NewConsoleLogWriter()
	if l == nil {
		t.Error("NewConsoleLogWriter produced a nil")
	}
	l.Close()
}

func TestConsoleLogWriterLogWrite(t *testing.T) {
	l := NewConsoleLogWriter()
	logRecord := &log4go.LogRecord{}
	l.LogWrite(logRecord) // Void function. Testing for coverage.
	l.Close()
}
