package logger

import (
	"context"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"path/filepath"
	"strconv"
	"time"
)

type ZapLogger struct {
	Level Level // 日志级别
	log   *zap.SugaredLogger
}

func NewZapLogger(properties *Properties) *ZapLogger {
	level := ParseLevel(properties.Level)

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

	zapLogger := zap.New(coreConfig, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel), zap.AddCallerSkip(2))
	return &ZapLogger{
		Level: level,
		log:   zapLogger.Sugar(),
	}
}

func (z *ZapLogger) IsDebugEnabled() bool {
	return z.Level >= DebugLevel
}

func (z *ZapLogger) IsInfoEnabled() bool {
	return z.Level >= InfoLevel
}

func (z *ZapLogger) IsWarnEnabled() bool {
	return z.Level >= WarnLevel
}

func (z *ZapLogger) IsErrorEnabled() bool {
	return z.Level >= ErrorLevel
}

func (z *ZapLogger) IsFatalEnabled() bool {
	return z.Level >= FatalLevel
}

func (z *ZapLogger) Debug(v ...interface{}) {
	z.log.Debug(v...)
}

func (z *ZapLogger) Debugf(template string, v ...interface{}) {
	z.log.Debugf(template, v...)
}

func (z *ZapLogger) DebugWithContext(context *context.Context, v ...interface{}) {
	z.log.Debug(v...)
}

func (z *ZapLogger) DebugfWithContext(context *context.Context, template string, v ...interface{}) {
	z.log.Debugf(template, v...)
}

func (z *ZapLogger) Info(v ...interface{}) {
	z.log.Info(v...)
}

func (z *ZapLogger) Infof(template string, v ...interface{}) {
	z.log.Infof(template, v...)
}

func (z *ZapLogger) InfoWithContext(context *context.Context, v ...interface{}) {
	z.log.Info(v...)
}

func (z *ZapLogger) InfofWithContext(context *context.Context, template string, v ...interface{}) {
	z.log.Infof(template, v...)
}

func (z *ZapLogger) Warn(v ...interface{}) {
	z.log.Warn(v...)
}

func (z *ZapLogger) Warnf(template string, v ...interface{}) {
	z.Warnf(template, v...)
}

func (z *ZapLogger) WarnWithContext(context *context.Context, v ...interface{}) {
	z.log.Warn(v...)
}

func (z *ZapLogger) WarnfWithContext(context *context.Context, template string, v ...interface{}) {
	z.log.Warnf(template, v...)
}

func (z *ZapLogger) Error(v ...interface{}) {
	z.log.Error(v...)
}

func (z *ZapLogger) Errorf(template string, v ...interface{}) {
	z.log.Errorf(template, v...)
}

func (z *ZapLogger) ErrorWithContext(context *context.Context, v ...interface{}) {
	z.log.Error(v...)
}

func (z *ZapLogger) ErrorfWithContext(context *context.Context, template string, v ...interface{}) {
	z.log.Errorf(template, v...)
}

func (z *ZapLogger) Fatal(v ...interface{}) {
	z.log.Fatal(v...)
}

func (z *ZapLogger) Fatalf(template string, v ...interface{}) {
	z.log.Fatalf(template, v...)
}

func (z *ZapLogger) FatalWithContext(context *context.Context, v ...interface{}) {
	z.log.Fatal(v...)
}

func (z *ZapLogger) FatalfWithContext(context *context.Context, template string, v ...interface{}) {
	z.log.Fatalf(template, v...)
}
