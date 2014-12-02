package log_test

import (
	"code.google.com/p/log4go"
	"github.com/limetext/lime/backend/log"
	"testing"
)

func TestNewFileLogWriter(t *testing.T) {
	l := log.NewFileLogWriter("some file", true)
	if l == nil {
		t.Error("NewFileLogWriter produced a nil")
	}
	l.Close()
}

func TestFileLogWriterLogWrite(t *testing.T) {
	l := log.NewFileLogWriter("some file", true)
	logRecord := &log4go.LogRecord{}
	l.LogWrite(logRecord) // Void function. Testing for coverage.
	l.Close()
}
