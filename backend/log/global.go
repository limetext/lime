package log

import (
	"code.google.com/p/log4go"
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

func LogFinest(arg0 interface{}, args ...interface{}) {
	Global.LogFinest(arg0, args)
}

func LogFine(arg0 interface{}, args ...interface{}) {
	Global.LogFine(arg0, args)
}

func LogDebug(arg0 interface{}, args ...interface{}) {
	Global.LogDebug(arg0, args)
}

func LogTrace(arg0 interface{}, args ...interface{}) {
	Global.LogTrace(arg0, args)
}

func LogInfo(arg0 interface{}, args ...interface{}) {
	Global.LogInfo(arg0, args)
}

func LogWarning(arg0 interface{}, args ...interface{}) {
	Global.LogWarning(arg0, args)
}

func LogError(arg0 interface{}, args ...interface{}) {
	Global.LogError(arg0, args)
}

func LogCritical(arg0 interface{}, args ...interface{}) {
	Global.LogCritical(arg0, args)
}

func Logf(level Level, format string, args ...interface{}) {
	Global.Logf(level, format, args)
}

func Close(args ...interface{}) {
	Global.Close(args)
}
