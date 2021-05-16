package logger

import (
	"context"
	"go.uber.org/zap"
)

type ZapLogger struct {
	Level Level // 日志级别
	log   *zap.SugaredLogger
}

func (z *ZapLogger) Flush() {
	_ = z.log.Sync()
}

func (z *ZapLogger) GetLevel() Level {
	return z.Level
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
