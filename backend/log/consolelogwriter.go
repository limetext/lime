// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package log

import (
	"code.google.com/p/log4go"
)

type (
	ConsoleLogWriter struct {
		logWriter
		writer log4go.ConsoleLogWriter
	}
)

func NewConsoleLogWriter() *ConsoleLogWriter {
	ret := &ConsoleLogWriter{
		writer: log4go.NewConsoleLogWriter(),
	}
	return ret
}

// Implement logWriter

func (l *ConsoleLogWriter) LogWrite(rec *log4go.LogRecord) {
	l.writer.LogWrite(rec)
}

func (l *ConsoleLogWriter) Close() {
	l.writer.Close()
}
