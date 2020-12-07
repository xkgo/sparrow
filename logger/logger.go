package logger

import (
	"context"
	"github.com/xkgo/sparrow/util/StringUtils"
)

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

func ParseLevel(level string) Level {
	if StringUtils.EqualsIgnoreCase("DEBUG", level) {
		return DebugLevel
	}
	if StringUtils.EqualsIgnoreCase("INFO", level) {
		return InfoLevel
	}
	if StringUtils.EqualsIgnoreCase("WARN", level) {
		return WarnLevel
	}
	if StringUtils.EqualsIgnoreCase("ERROR", level) {
		return ErrorLevel
	}
	if StringUtils.EqualsIgnoreCase("FATAL", level) {
		return FatalLevel
	}
	return DebugLevel
}

/**
日志配置
*/
type Properties struct {
	Level      string `ck:"level" def:"DEBUG"`                         // 日志级别: DEBUG, INFO, WARN, ERROR, FATAL， 默认是： DEBUG
	Dir        string `ck:"dir" def:"./logs"`                          // 日志存放目录, 默认是 ./logs
	Filename   string `ck:"filename" def:"app.log"`                    // 文件名，含后缀, 默认：app.log
	TimeFormat string `ck:"time-format" def:"2006-01-02 15:04:05.000"` // 时间格式，默认是 2006-01-02 15:04:05.000
	MaxSize    int    `ck:"max-size" def:"500"`                        // 单个配置文件大小最大限制，单位：M，默认是 500 M
	MaxBackups int    `ck:"max-backups" def:"30"`                      // 最多保留多少个日志文件，默认 30
	MaxAge     int    `ck:"max-age" def:"30"`                          // 日志文件存活时间，单位：天，默认是30天
	Compress   bool   `ck:"compress" def:"false"`                      // 是否需要自动gzip进行压缩，默认：false
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
var consoleLogger Logger = &ConsoleLogger{Level: DebugLevel}

func RootLogger() *Logger {
	return &rootLogger
}

func SetRootLogger(logger Logger) {
	rootLogger = logger
}

func SetConsoleLogger(logger Logger) {
	consoleLogger = logger
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
	if nil != consoleLogger {
		consoleLogger.Debug(v...)
	}
}

func Debugf(template string, v ...interface{}) {
	rootLogger.Debugf(template, v...)
	if nil != consoleLogger {
		consoleLogger.Debugf(template, v...)
	}
}

func DebugWithContext(context *context.Context, v ...interface{}) {
	rootLogger.DebugWithContext(context, v...)
	if nil != consoleLogger {
		consoleLogger.DebugWithContext(context, v...)
	}
}

func DebugfWithContext(context *context.Context, template string, v ...interface{}) {
	rootLogger.DebugfWithContext(context, template, v...)
	if consoleLogger != nil {
		consoleLogger.DebugfWithContext(context, template, v...)
	}
}

func Info(v ...interface{}) {
	rootLogger.Info(v...)
	if nil != consoleLogger {
		consoleLogger.Info(v...)
	}
}
func Infof(template string, v ...interface{}) {
	rootLogger.Infof(template, v...)
	if nil != consoleLogger {
		consoleLogger.Infof(template, v...)
	}
}

func InfoWithContext(context *context.Context, v ...interface{}) {
	rootLogger.InfoWithContext(context, v...)
	if nil != consoleLogger {
		consoleLogger.InfoWithContext(context, v...)
	}
}
func InfofWithContext(context *context.Context, template string, v ...interface{}) {
	rootLogger.InfofWithContext(context, template, v...)
	if nil != consoleLogger {
		consoleLogger.InfofWithContext(context, template, v...)
	}
}

func Warn(v ...interface{}) {
	rootLogger.Warn(v...)
	if nil != consoleLogger {
		consoleLogger.Warn(v...)
	}
}
func Warnf(template string, v ...interface{}) {
	rootLogger.Warnf(template, v...)
	if nil != consoleLogger {
		consoleLogger.Warnf(template, v...)
	}
}

func WarnWithContext(context *context.Context, v ...interface{}) {
	rootLogger.WarnWithContext(context, v...)
	if nil != consoleLogger {
		consoleLogger.WarnWithContext(context, v...)
	}
}
func WarnfWithContext(context *context.Context, template string, v ...interface{}) {
	rootLogger.WarnfWithContext(context, template, v...)
	if nil != consoleLogger {
		consoleLogger.WarnfWithContext(context, template, v...)
	}
}

func Error(v ...interface{}) {
	rootLogger.Error(v...)
	if nil != consoleLogger {
		consoleLogger.Error(v...)
	}
}
func Errorf(template string, v ...interface{}) {
	rootLogger.Errorf(template, v...)
	if nil != consoleLogger {
		consoleLogger.Errorf(template, v...)
	}
}

func ErrorWithContext(context *context.Context, v ...interface{}) {
	rootLogger.ErrorWithContext(context, v...)
	if nil != consoleLogger {
		consoleLogger.ErrorWithContext(context, v...)
	}
}
func ErrorfWithContext(context *context.Context, template string, v ...interface{}) {
	rootLogger.ErrorfWithContext(context, template, v...)
	if nil != consoleLogger {
		consoleLogger.ErrorfWithContext(context, template, v...)
	}
}

func Fatal(v ...interface{}) {
	rootLogger.Fatal(v...)
	if nil != consoleLogger {
		consoleLogger.Fatal(v...)
	}
}

func Fatalf(template string, v ...interface{}) {
	rootLogger.Fatalf(template, v...)
	if nil != consoleLogger {
		consoleLogger.Fatalf(template, v...)
	}
}

func FatalWithContext(context *context.Context, v ...interface{}) {
	rootLogger.FatalWithContext(context, v...)
	if nil != consoleLogger {
		consoleLogger.FatalWithContext(context, v...)
	}
}

func FatalfWithContext(context *context.Context, template string, v ...interface{}) {
	rootLogger.FatalfWithContext(context, template, v...)
	if nil != consoleLogger {
		consoleLogger.FatalfWithContext(context, template, v...)
	}
}
