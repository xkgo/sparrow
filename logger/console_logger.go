package logger

import (
	"context"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type ConsoleLogger struct {
	Level            Level // 日志级别
	CallerSkipOffset int   // 输出日志时候，计算输入日志的日志所在文件和行数偏移，一般给应用进行二次封装使用，正负数都可以
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
	now := time.Now()
	zone, offset := now.Zone()
	timeLabel := now.Format("2006-01-02 15:04:05.000")
	timeLabel = timeLabel + ":" + zone + ":" + strconv.FormatInt(int64(offset), 10)

	_, file, line, ok := runtime.Caller(4 + c.CallerSkipOffset)
	if ok {

		idx := strings.LastIndexByte(file, '/')
		if idx != -1 {
			// Find the penultimate separator.
			idx = strings.LastIndexByte(file[:idx], '/')
			if idx != -1 {
				file = file[idx+1:]
			}
		}

		msg = file + ":" + strconv.FormatInt(int64(line), 10) + "\t" + msg
	}
	msg = timeLabel + "\t" + level.String() + "\t" + msg

	fmt.Println(msg)

	// 如果是 ERROR，FATAL 则输出堆栈信息
	if level >= ErrorLevel {
		buff := make([]byte, 1<<10)
		// 堆栈信息
		runtime.Stack(buff, true)
		fmt.Printf("%v", string(buff))
	}

}

func (c *ConsoleLogger) Flush() {}

func (c *ConsoleLogger) GetLevel() Level {
	return c.Level
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
