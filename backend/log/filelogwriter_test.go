// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package log

import (
	"github.com/limetext/log4go"
	"os"
	"testing"
)

const logfn = "some file"

func TestNewFileLogWriter(t *testing.T) {
	l := NewFileLogWriter(logfn, true)
	defer os.Remove(logfn)
	if l == nil {
		t.Error("NewFileLogWriter produced a nil")
	}
	l.Close()
}

func TestFileLogWriterLogWrite(t *testing.T) {
	l := NewFileLogWriter(logfn, true)
	defer os.Remove(logfn)
	logRecord := &log4go.LogRecord{}
	l.LogWrite(logRecord) // Void function. Testing for coverage.
	l.Close()
}
