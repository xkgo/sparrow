package logger

import "context"

type Level int8

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "Fatal"
	}
	return "UNKNOWN"
}

/*
日志初始化
*/
type Logger interface {
	IsDebugEnabled() bool
	IsInfoEnabled() bool
	IsWarnEnabled() bool
	IsErrorEnabled() bool
	IsFatalEnabled() bool

	Debug(v ...interface{})
	Debugf(template string, v ...interface{})
	DebugWithContext(context *context.Context, v ...interface{})
	DebugfWithContext(context *context.Context, template string, v ...interface{})

	Info(v ...interface{})
	Infof(template string, v ...interface{})
	InfoWithContext(context *context.Context, v ...interface{})
	InfofWithContext(context *context.Context, template string, v ...interface{})

	Warn(v ...interface{})
	Warnf(template string, v ...interface{})
	WarnWithContext(context *context.Context, v ...interface{})
	WarnfWithContext(context *context.Context, template string, v ...interface{})

	Error(v ...interface{})
	Errorf(template string, v ...interface{})
	ErrorWithContext(context *context.Context, v ...interface{})
	ErrorfWithContext(context *context.Context, template string, v ...interface{})

	Fatal(v ...interface{})
	Fatalf(template string, v ...interface{})
	FatalWithContext(context *context.Context, v ...interface{})
	FatalfWithContext(context *context.Context, template string, v ...interface{})
}

// 日志
var rootLogger Logger = &ConsoleLogger{Level: DebugLevel}

func RootLogger() *Logger {
	return &rootLogger
}

func SetRootLogger(logger Logger) {
	rootLogger = logger
}

func IsDebugEnabled() bool {
	return rootLogger.IsDebugEnabled()
}
func IsInfoEnabled() bool {
	return rootLogger.IsInfoEnabled()
}
func IsWarnEnabled() bool {
	return rootLogger.IsWarnEnabled()
}
func IsErrorEnabled() bool {
	return IsErrorEnabled()
}
func IsFatalEnabled() bool {
	return IsFatalEnabled()
}

func Debug(v ...interface{}) {
	rootLogger.Debug(v...)
}

func Debugf(template string, v ...interface{}) {
	rootLogger.Debugf(template, v...)
}

func DebugWithContext(context *context.Context, v ...interface{}) {
	rootLogger.DebugWithContext(context, v...)
}

func DebugfWithContext(context *context.Context, template string, v ...interface{}) {
	rootLogger.DebugfWithContext(context, template, v...)
}

func Info(v ...interface{}) {
	rootLogger.Info(v...)
}
func Infof(template string, v ...interface{}) {
	rootLogger.Infof(template, v...)
}

func InfoWithContext(context *context.Context, v ...interface{}) {
	rootLogger.InfoWithContext(context, v...)
}
func InfofWithContext(context *context.Context, template string, v ...interface{}) {
	rootLogger.InfofWithContext(context, template, v...)
}

func Warn(v ...interface{}) {
	rootLogger.Warn(v...)
}
func Warnf(template string, v ...interface{}) {
	rootLogger.Warnf(template, v...)
}

func WarnWithContext(context *context.Context, v ...interface{}) {
	rootLogger.WarnWithContext(context, v...)
}
func WarnfWithContext(context *context.Context, template string, v ...interface{}) {
	rootLogger.WarnfWithContext(context, template, v...)
}

func Error(v ...interface{}) {
	rootLogger.Error(v...)
}
func Errorf(template string, v ...interface{}) {
	rootLogger.Errorf(template, v...)
}

func ErrorWithContext(context *context.Context, v ...interface{}) {
	rootLogger.ErrorWithContext(context, v...)
}
func ErrorfWithContext(context *context.Context, template string, v ...interface{}) {
	rootLogger.ErrorfWithContext(context, template, v...)
}

func Fatal(v ...interface{}) {
	rootLogger.Fatal(v...)
}

func Fatalf(template string, v ...interface{}) {
	rootLogger.Fatalf(template, v...)
}

func FatalWithContext(context *context.Context, v ...interface{}) {
	rootLogger.FatalWithContext(context, v...)
}

func FatalfWithContext(context *context.Context, template string, v ...interface{}) {
	rootLogger.FatalfWithContext(context, template, v...)
}
