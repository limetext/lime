// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package log_test

import (
	"github.com/limetext/lime/backend/log"
	"github.com/limetext/log4go"
	"sync"
	"testing"
	"time"
)

type testlogger func(string)

func (l testlogger) LogWrite(rec *log4go.LogRecord) {
	l(rec.Message)
}

func (l testlogger) Close() {}

func TestGlobalLog(t *testing.T) {
	var wg sync.WaitGroup
	log.Global.Close()
	log.Global.AddFilter("globaltest", log.FINEST, testlogger(func(str string) {
		if str != "Testing: hello world" {
			t.Errorf("got: %s", str)
		}
		wg.Done()
	}))
	wg.Add(1)
	log.Info("Testing: %s %s", "hello", "world")
	wg.Wait()
}

func TestLogf(t *testing.T) {
	l := log.NewLogger()

	// Log a message at each level. Because we cannot access the internals of the logger,
	// we assume that this test succeeds if it does not cause an error (although we cannot
	// actually look inside and see if the level was changed)
	for _, test_lvl := range []log.Level{log.FINEST, log.FINE, log.DEBUG, log.TRACE, log.INFO, log.WARNING, log.ERROR, log.CRITICAL, 999} {
		l.Logf(test_lvl, time.Now().String())
	}
}

func TestClose(t *testing.T) {
	l := log.NewLogger()
	l.Close()
	m := log.NewLogger()
	m.Close("something wrong")
}

func TestNewLogger(t *testing.T) {
	l := log.NewLogger()
	if l == nil {
		t.Error("Returned a nil logger")
	}
}

func TestLogLevels(t *testing.T) {
	l := log.NewLogger()

	// Again, because we cannot access the internals of log this will
	// succeed as long there is no error
	for _, test_lvl := range []log.Level{log.FINEST, log.FINE, log.DEBUG, log.TRACE, log.INFO, log.WARNING, log.ERROR, log.CRITICAL, 999} {
		// Use a random-ish string (the current time)
		l.AddFilter(time.Now().String(), test_lvl, testlogger(func(str string) {}))
	}
}

func TestLogFunctions(t *testing.T) {
	l := log.NewLogger()

	l.Finest(time.Now().String())
	l.Fine(time.Now().String())
	l.Debug(time.Now().String())
	l.Trace(time.Now().String())
	l.Warn(time.Now().String())
	l.Error(time.Now().String())
	l.Critical(time.Now().String())

}
