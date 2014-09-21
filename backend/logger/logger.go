// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package logger

import (
	"code.google.com/p/log4go"
	"fmt"
	. "github.com/limetext/lime/backend/util"
	"sync"
)

type (
	Logger struct {
		log     chan string
		handler func(s string)
		lock    sync.Mutex
	}
)

func NewLogger(h func(s string)) *Logger {
	ret := &Logger{
		log:     make(chan string, 100),
		handler: h,
	}
	go ret.handle()
	return ret
}

func (l *Logger) handle() {
	for fl := range l.log {
		l.handler(fl)
	}
}

func (l *Logger) LogWrite(rec *log4go.LogRecord) {
	p := Prof.Enter("log")
	defer p.Exit()
	l.lock.Lock()
	defer l.lock.Unlock()
	fl := log4go.FormatLogRecord(log4go.FORMAT_DEFAULT, rec)
	l.log <- fl
}

func (l *Logger) Close() {
	l.lock.Lock()
	defer l.lock.Unlock()
	fmt.Println("Closing...")
	close(l.log)
}
