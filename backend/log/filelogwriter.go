// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package log

import (
	"code.google.com/p/log4go"
)

type (
	FileLogWriter struct {
		logWriter
		writer *log4go.FileLogWriter
	}
)

func NewFileLogWriter(fname string, rotate bool) *FileLogWriter {
	ret := &FileLogWriter{
		writer: log4go.NewFileLogWriter(fname, rotate),
	}
	return ret
}

// Implement logWriter

func (l *FileLogWriter) LogWrite(rec *log4go.LogRecord) {
	l.writer.LogWrite(rec)
}

func (l *FileLogWriter) Close() {
	l.writer.Close()
}
