package log_test

import (
	"github.com/limetext/lime/backend/log"
	"github.com/limetext/log4go"
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
