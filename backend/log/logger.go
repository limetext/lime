// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package log

import (
	"code.google.com/p/log4go"
)

type (
	Logger struct {
		logger log4go.Logger
	}
)

var (
	Global *Logger
)

func NewLogger() *Logger {
	l := &Logger{
		logger: make(log4go.Logger),
	}
	return l
}

func (l *Logger) AddFilter(name string, level Level, writer logWriter) {
	lvl := log4go.INFO
	switch level {
	case FINEST:
		lvl = log4go.FINEST
	case FINE:
		lvl = log4go.FINE
	case DEBUG:
		lvl = log4go.DEBUG
	case TRACE:
		lvl = log4go.TRACE
	case INFO:
		lvl = log4go.INFO
	case WARNING:
		lvl = log4go.WARNING
	case ERROR:
		lvl = log4go.ERROR
	case CRITICAL:
		lvl = log4go.CRITICAL
	default:
	}
	l.logger.AddFilter(name, lvl, writer)
}

func (l *Logger) LogFinest(arg0 interface{}, args ...interface{}) {
	l.logger.Finest(arg0, args)
}

func (l *Logger) LogFine(arg0 interface{}, args ...interface{}) {
	l.logger.Fine(arg0, args)
}

func (l *Logger) LogDebug(arg0 interface{}, args ...interface{}) {
	l.logger.Debug(arg0, args)
}

func (l *Logger) LogTrace(arg0 interface{}, args ...interface{}) {
	l.logger.Trace(arg0, args)
}

func (l *Logger) LogInfo(arg0 interface{}, args ...interface{}) {
	l.logger.Info(arg0, args)
}

func (l *Logger) LogWarning(arg0 interface{}, args ...interface{}) {
	l.logger.Warn(arg0, args)
}

func (l *Logger) LogError(arg0 interface{}, args ...interface{}) {
	l.logger.Error(arg0, args)
}

func (l *Logger) LogCritical(arg0 interface{}, args ...interface{}) {
	l.logger.Critical(arg0, args)
}

func (l *Logger) Logf(level Level, format string, args ...interface{}) {
	lvl := log4go.INFO
	switch level {
	case FINEST:
		lvl = log4go.FINEST
	case FINE:
		lvl = log4go.FINE
	case DEBUG:
		lvl = log4go.DEBUG
	case TRACE:
		lvl = log4go.TRACE
	case INFO:
		lvl = log4go.INFO
	case WARNING:
		lvl = log4go.WARNING
	case ERROR:
		lvl = log4go.ERROR
	case CRITICAL:
		lvl = log4go.CRITICAL
	default:
	}
	l.logger.Logf(lvl, format, args)
}

func (l *Logger) Close(args ...interface{}) {
	if len(args) > 0 {
		l.LogError(args)
	}
	l.logger.Close()
}

func init() {
	log4go.Global.Close()
	Global = &Logger{
		logger: log4go.Global,
	}
}
