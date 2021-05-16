package logger

import (
	"context"
	"fmt"
	"testing"
)

func TestSetTraceIdGenerator(t *testing.T) {

	properties := &Properties{
		Level:            "DEBUG",
		Dir:              "./logs",
		Filename:         "app.log",
		TimeFormat:       "2006-01-02 15:04:05.000",
		MaxSize:          100,
		MaxBackups:       30,
		MaxAge:           30,
		Compress:         false,
		ConsoleLog:       true,
		CallerSkipOffset: 0,
	}

	fmt.Println(properties)

	InitLogger(properties)

	SetTraceIdGenerator(func(ctx *context.Context) string {
		return "xxxxxxxxxx"
	})

	Debug("你好")
	Info("你好")
	Warn("你好")
	Error("你好")
	Fatal("你好")
}
