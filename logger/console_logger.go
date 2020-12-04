package logger

import (
	"context"
	"fmt"
	"time"
)

type ConsoleLogger struct {
	Level Level // 日志级别
}

func (c *ConsoleLogger) IsDebugEnabled() bool {
	return c.Level >= DebugLevel
}

func (c *ConsoleLogger) IsInfoEnabled() bool {
	return c.Level >= InfoLevel
}

func (c *ConsoleLogger) IsWarnEnabled() bool {
	return c.Level >= WarnLevel
}

func (c *ConsoleLogger) IsErrorEnabled() bool {
	return c.Level >= ErrorLevel
}

func (c *ConsoleLogger) IsFatalEnabled() bool {
	return c.Level >= FatalLevel
}

func (c *ConsoleLogger) log(context *context.Context, level Level, template string, fmtArgs ...interface{}) {
	if level < c.Level {
		return
	}

	// Format with Sprint, Sprintf, or neither.
	msg := template
	if msg == "" && len(fmtArgs) > 0 {
		msg = fmt.Sprint(fmtArgs...)
	} else if msg != "" && len(fmtArgs) > 0 {
		msg = fmt.Sprintf(template, fmtArgs...)
	}
	timeLabel := time.Now().Format("2006-01-02 15:04:05.000")
	fmt.Println(timeLabel, level, msg)
}

func (c *ConsoleLogger) Debug(v ...interface{}) {
	c.log(nil, DebugLevel, "", v...)
}

func (c *ConsoleLogger) Debugf(template string, v ...interface{}) {
	c.log(nil, DebugLevel, template, v...)
}

func (c *ConsoleLogger) DebugWithContext(context *context.Context, v ...interface{}) {
	c.log(context, DebugLevel, "", v...)
}

func (c *ConsoleLogger) DebugfWithContext(context *context.Context, template string, v ...interface{}) {
	c.log(context, DebugLevel, template, v...)
}

func (c *ConsoleLogger) Info(v ...interface{}) {
	c.log(nil, InfoLevel, "", v...)
}

func (c *ConsoleLogger) Infof(template string, v ...interface{}) {
	c.log(nil, InfoLevel, template, v...)
}

func (c *ConsoleLogger) InfoWithContext(context *context.Context, v ...interface{}) {
	c.log(context, InfoLevel, "", v...)
}

func (c *ConsoleLogger) InfofWithContext(context *context.Context, template string, v ...interface{}) {
	c.log(context, InfoLevel, template, v...)
}

func (c *ConsoleLogger) Warn(v ...interface{}) {
	c.log(nil, WarnLevel, "", v...)
}

func (c *ConsoleLogger) Warnf(template string, v ...interface{}) {
	c.log(nil, WarnLevel, template, v...)
}

func (c *ConsoleLogger) WarnWithContext(context *context.Context, v ...interface{}) {
	c.log(context, WarnLevel, "", v...)
}

func (c *ConsoleLogger) WarnfWithContext(context *context.Context, template string, v ...interface{}) {
	c.log(context, WarnLevel, template, v...)
}

func (c *ConsoleLogger) Error(v ...interface{}) {
	c.log(nil, ErrorLevel, "", v...)
}

func (c *ConsoleLogger) Errorf(template string, v ...interface{}) {
	c.log(nil, ErrorLevel, template, v...)
}

func (c *ConsoleLogger) ErrorWithContext(context *context.Context, v ...interface{}) {
	c.log(context, ErrorLevel, "", v...)
}

func (c *ConsoleLogger) ErrorfWithContext(context *context.Context, template string, v ...interface{}) {
	c.log(context, ErrorLevel, template, v...)
}

func (c *ConsoleLogger) Fatal(v ...interface{}) {
	c.log(nil, FatalLevel, "", v...)
}

func (c *ConsoleLogger) Fatalf(template string, v ...interface{}) {
	c.log(nil, FatalLevel, template, v...)
}

func (c *ConsoleLogger) FatalWithContext(context *context.Context, v ...interface{}) {
	c.log(context, FatalLevel, "", v...)
}

func (c *ConsoleLogger) FatalfWithContext(context *context.Context, template string, v ...interface{}) {
	c.log(context, FatalLevel, template, v...)
}
