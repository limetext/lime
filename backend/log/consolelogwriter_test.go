package log_test

import (
	"github.com/limetext/lime/backend/log"
	"github.com/limetext/log4go"
	"testing"
)

func TestNewConsoleLogWriter(t *testing.T) {
	l := log.NewConsoleLogWriter()
	if l == nil {
		t.Error("NewConsoleLogWriter produced a nil")
	}
	l.Close()
}

func TestConsoleLogWriterLogWrite(t *testing.T) {
	l := log.NewConsoleLogWriter()
	logRecord := &log4go.LogRecord{}
	l.LogWrite(logRecord) // Void function. Testing for coverage.
	l.Close()
}
