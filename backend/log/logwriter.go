// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package log

import (
	"code.google.com/p/log4go"
	. "github.com/limetext/lime/backend/util"
	"sync"
)

type (
	LogWriter interface {
		log4go.LogWriter
	}

	logWriter struct {
		LogWriter
		log     chan string
		handler func(string)
		lock    sync.Mutex
	}
)

func NewLogWriter(h func(string)) *logWriter {
	ret := &logWriter{
		log:     make(chan string, 100),
		handler: h,
	}
	go ret.handle()
	return ret
}

func (l *logWriter) handle() {
	for fl := range l.log {
		l.handler(fl)
	}
}

// Implement LogWriter

func (l *logWriter) LogWrite(rec *log4go.LogRecord) {
	p := Prof.Enter("log")
	defer p.Exit()
	l.lock.Lock()
	defer l.lock.Unlock()
	fl := log4go.FormatLogRecord(log4go.FORMAT_DEFAULT, rec)
	l.log <- fl
}

func (l *logWriter) Close() {
	l.lock.Lock()
	defer l.lock.Unlock()
	close(l.log)
}
