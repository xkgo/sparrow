package logger

import (
	"context"
	"fmt"
	"github.com/natefinch/lumberjack"
	"github.com/xkgo/sparrow/util/StringUtils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"path/filepath"
	"strconv"
	"time"
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
	Level            string `ck:"level" def:"DEBUG"`                         // 日志级别: DEBUG, INFO, WARN, ERROR, FATAL， 默认是： DEBUG
	Dir              string `ck:"dir" def:"./logs"`                          // 日志存放目录, 默认是 ./logs
	Filename         string `ck:"filename" def:"app.log"`                    // 文件名，含后缀, 默认：app.log
	TimeFormat       string `ck:"time-format" def:"2006-01-02 15:04:05.000"` // 时间格式，默认是 2006-01-02 15:04:05.000
	MaxSize          int    `ck:"max-size" def:"500"`                        // 单个配置文件大小最大限制，单位：M，默认是 500 M
	MaxBackups       int    `ck:"max-backups" def:"30"`                      // 最多保留多少个日志文件，默认 30
	MaxAge           int    `ck:"max-age" def:"30"`                          // 日志文件存活时间，单位：天，默认是30天
	Compress         bool   `ck:"compress" def:"false"`                      // 是否需要自动gzip进行压缩，默认：false
	ConsoleLog       bool   `ck:"console-log" def:"false"`                   // 是否需要输出控制台日志，默认是 false
	CallerSkipOffset int    `ck:"caller-skip-offset" def:"0"`                // 输出日志时候，计算输入日志的日志所在文件和行数偏移，一般给应用进行二次封装使用，正负数都可以
}

func (p *Properties) Equals(properties *Properties) bool {
	if nil == properties {
		return false
	}
	if p.Level != properties.Level {
		return false
	}
	if p.Dir != properties.Dir {
		return false
	}
	if p.Filename != properties.Filename {
		return false
	}
	if p.TimeFormat != properties.TimeFormat {
		return false
	}
	if p.MaxSize != properties.MaxSize {
		return false
	}
	if p.MaxBackups != properties.MaxBackups {
		return false
	}
	if p.MaxAge != properties.MaxAge {
		return false
	}
	if p.Compress != properties.Compress {
		return false
	}
	if p.ConsoleLog != properties.ConsoleLog {
		return false
	}
	return true
}

/**
处理默认配置信息
*/
func ResolveAndApplyDefaultProperties(properties *Properties) {
	if len(properties.Level) < 1 {
		properties.Level = DebugLevel.String()
	}
	if len(properties.Dir) < 1 {
		properties.Dir = "./logs"
	}
	if len(properties.Filename) < 1 {
		properties.Filename = "app.log"
	}
	if len(properties.TimeFormat) < 1 {
		properties.TimeFormat = "2006-01-02 15:04:05.000"
	}
	if properties.MaxSize < 1 {
		properties.MaxSize = 500
	}
	if properties.MaxBackups < 1 {
		properties.MaxBackups = 30
	}
	if properties.MaxAge < 1 {
		properties.MaxAge = 30
	}
}

/*
日志初始化
*/
type Logger interface {
	Flush()
	GetLevel() Level
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

func InitLogger(properties *Properties) {
	ResolveAndApplyDefaultProperties(properties)

	level := ParseLevel(properties.Level)
	if properties.ConsoleLog {
		consoleLogger = &ConsoleLogger{Level: level, CallerSkipOffset: properties.CallerSkipOffset}
	} else {
		consoleLogger = nil
	}

	zapLevel := zapcore.DebugLevel
	switch level {
	case DebugLevel:
		zapLevel = zapcore.DebugLevel
	case InfoLevel:
		zapLevel = zapcore.InfoLevel
	case WarnLevel:
		zapLevel = zapcore.WarnLevel
	case ErrorLevel:
		zapLevel = zapcore.ErrorLevel
	case FatalLevel:
		zapLevel = zapcore.FatalLevel
	default:
		zapLevel = zapcore.DebugLevel
	}

	// 构造新的
	var writerSyncer zapcore.WriteSyncer
	// 输出到文件中去
	lumberJackLogger := &lumberjack.Logger{
		Filename:   properties.Dir + string(filepath.Separator) + properties.Filename,
		MaxSize:    properties.MaxSize,
		MaxAge:     properties.MaxAge,
		MaxBackups: properties.MaxBackups,
		Compress:   properties.Compress,
	}
	writerSyncer = zapcore.AddSync(lumberJackLogger)

	var encoder zapcore.Encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		str := t.Format(properties.TimeFormat)
		zone, offset := t.Zone()
		enc.AppendString(str + ":" + zone + ":" + strconv.FormatInt(int64(offset), 10))
	}

	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder = zapcore.NewConsoleEncoder(encoderConfig)

	var coreConfig = zapcore.NewCore(encoder, writerSyncer, zapLevel)

	zapLogger := zap.New(coreConfig, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel), zap.AddCallerSkip(3+properties.CallerSkipOffset))

	SetRootLogger(&ZapLogger{
		Level: level,
		log:   zapLogger.Sugar(),
	})
}

/**
记录TraceId，提供扩展方法，支持从context 中获取 context
*/

/**
TraceId 生成器, 允许用户自定义
*/
type TraceIdGenerator func(ctx *context.Context) string

// 日志
var rootLogger Logger = &ConsoleLogger{Level: DebugLevel}
var consoleLogger Logger

var traceIdGenerator TraceIdGenerator

// 写入日志之后的处理逻辑
var afterLogHandler func(ctx *context.Context, logText string, level Level)

func SetTraceIdGenerator(generator TraceIdGenerator) {
	traceIdGenerator = generator
}

func SetAfterLogHandler(handler func(ctx *context.Context, logText string, level Level)) {
	afterLogHandler = func(ctx *context.Context, logText string, level Level) {
		defer func() {
			if r := recover(); r != nil {
				_ = fmt.Errorf("执行AfterLogHandler 异常: %v", r)
			}
		}()
		handler(ctx, logText, level)
	}
}

func RootLogger() *Logger {
	return &rootLogger
}

func SetRootLogger(logger Logger) {
	rootLogger = logger
}

func Flush() {
	rootLogger.Flush()
}

func GetLevel() Level {
	return rootLogger.GetLevel()
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

func log(context *context.Context, level Level, template string, fmtArgs ...interface{}) {
	if level < rootLogger.GetLevel() {
		return
	}

	// Format with Sprint, Sprintf, or neither.
	msg := template
	if msg == "" && len(fmtArgs) > 0 {
		msg = fmt.Sprint(fmtArgs...)
	} else if msg != "" && len(fmtArgs) > 0 {
		msg = fmt.Sprintf(template, fmtArgs...)
	}
	if traceIdGenerator != nil {
		traceId := traceIdGenerator(context)
		if len(traceId) > 0 {
			msg = traceIdGenerator(context) + " " + msg
		}
	}
	defer func() {
		if nil != afterLogHandler {
			afterLogHandler(context, msg, level)
		}
	}()
	switch level {
	case DebugLevel:
		rootLogger.Debug(msg)
		if nil != consoleLogger {
			consoleLogger.Debug(msg)
		}
	case InfoLevel:
		rootLogger.Info(msg)
		if nil != consoleLogger {
			consoleLogger.Info(msg)
		}
	case WarnLevel:
		rootLogger.Warn(msg)
		if nil != consoleLogger {
			consoleLogger.Warn(msg)
		}
	case ErrorLevel:
		rootLogger.Error(msg)
		if nil != consoleLogger {
			consoleLogger.Error(msg)
		}
	case FatalLevel:
		rootLogger.Fatal(msg)
		if nil != consoleLogger {
			consoleLogger.Fatal(msg)
		}
	default:
		rootLogger.Debug(msg)
		if nil != consoleLogger {
			consoleLogger.Debug(msg)
		}
	}
}

func Debug(v ...interface{}) {
	log(nil, DebugLevel, "", v...)
}

func Debugf(template string, v ...interface{}) {
	log(nil, DebugLevel, template, v...)
}

func DebugWithContext(context *context.Context, v ...interface{}) {
	log(context, DebugLevel, "", v...)
}

func DebugfWithContext(context *context.Context, template string, v ...interface{}) {
	log(context, DebugLevel, template, v...)
}

func Info(v ...interface{}) {
	log(nil, InfoLevel, "", v...)
}

func Infof(template string, v ...interface{}) {
	log(nil, InfoLevel, template, v...)
}

func InfoWithContext(context *context.Context, v ...interface{}) {
	log(context, InfoLevel, "", v...)
}

func InfofWithContext(context *context.Context, template string, v ...interface{}) {
	log(context, InfoLevel, template, v...)
}

func Warn(v ...interface{}) {
	log(nil, WarnLevel, "", v...)
}

func Warnf(template string, v ...interface{}) {
	log(nil, WarnLevel, template, v...)
}

func WarnWithContext(context *context.Context, v ...interface{}) {
	log(context, WarnLevel, "", v...)
}

func WarnfWithContext(context *context.Context, template string, v ...interface{}) {
	log(context, WarnLevel, template, v...)
}

func Error(v ...interface{}) {
	log(nil, ErrorLevel, "", v...)
}

func Errorf(template string, v ...interface{}) {
	log(nil, ErrorLevel, template, v...)
}

func ErrorWithContext(context *context.Context, v ...interface{}) {
	log(context, ErrorLevel, "", v...)
}

func ErrorfWithContext(context *context.Context, template string, v ...interface{}) {
	log(context, ErrorLevel, template, v...)
}

func Fatal(v ...interface{}) {
	log(nil, FatalLevel, "", v...)
}

func Fatalf(template string, v ...interface{}) {
	log(nil, FatalLevel, template, v...)
}

func FatalWithContext(context *context.Context, v ...interface{}) {
	log(context, FatalLevel, "", v...)
}

func FatalfWithContext(context *context.Context, template string, v ...interface{}) {
	log(context, FatalLevel, template, v...)
}
