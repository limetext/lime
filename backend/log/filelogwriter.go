// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package log

import (
	"code.google.com/p/log4go"
)

type (
	fileLogWriter struct {
		logWriter
		writer *log4go.FileLogWriter
	}
)

func NewFileLogWriter(fname string, rotate bool) *fileLogWriter {
	ret := &fileLogWriter{
		writer: log4go.NewFileLogWriter(fname, rotate),
	}
	return ret
}

// Implement LogWriter

func (l *fileLogWriter) LogWrite(rec *log4go.LogRecord) {
	l.writer.LogWrite(rec)
}

func (l *fileLogWriter) Close() {
	l.writer.Close()
}
