package log_test

import (
	"code.google.com/p/log4go"
	"fmt"
	"github.com/limetext/lime/backend/log"
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
