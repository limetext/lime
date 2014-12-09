// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package log

import (
	"github.com/limetext/log4go"
)

var (
	Global *Logger
)

func init() {
	log4go.Global.Close()
	Global = &Logger{
		logger: log4go.Global,
	}
}

func AddFilter(name string, level Level, writer LogWriter) {
	Global.AddFilter(name, level, writer)
}

func Finest(arg0 interface{}, args ...interface{}) {
	Global.Finest(arg0, args...)
}

func Fine(arg0 interface{}, args ...interface{}) {
	Global.Fine(arg0, args...)
}

func Debug(arg0 interface{}, args ...interface{}) {
	Global.Debug(arg0, args...)
}

func Trace(arg0 interface{}, args ...interface{}) {
	Global.Trace(arg0, args...)
}

func Info(arg0 interface{}, args ...interface{}) {
	Global.Info(arg0, args...)
}

func Warn(arg0 interface{}, args ...interface{}) {
	Global.Warn(arg0, args...)
}

func Error(arg0 interface{}, args ...interface{}) {
	Global.Error(arg0, args...)
}

func Errorf(format string, args ...interface{}) {
	Global.Errorf(format, args...)
}

func Critical(arg0 interface{}, args ...interface{}) {
	Global.Critical(arg0, args...)
}

func Logf(level Level, format string, args ...interface{}) {
	Global.Logf(level, format, args...)
}

func Close(args ...interface{}) {
	Global.Close(args...)
}
