// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package log

import (
	"code.google.com/p/log4go"
)

type (
	consoleLogWriter struct {
		logWriter
		writer log4go.ConsoleLogWriter
	}
)

func NewConsoleLogWriter() *consoleLogWriter {
	ret := &consoleLogWriter{
		writer: log4go.NewConsoleLogWriter(),
	}
	return ret
}

// Implement LogWriter

func (l *consoleLogWriter) LogWrite(rec *log4go.LogRecord) {
	l.writer.LogWrite(rec)
}

func (l *consoleLogWriter) Close() {
	l.writer.Close()
}
