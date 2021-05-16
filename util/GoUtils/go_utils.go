package GoUtils

import (
	"context"
	"github.com/xkgo/sparrow/logger"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

/**
使用 Goroutine 执行，并且自动进行panic异常cache
*/
func RunGoroutine(handler func(), panicHandler func(r interface{}), afterHandlers ...func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("通过 Goroutine 执行任务Panic：", r)
				if nil != panicHandler {
					func() {
						defer func() {
							if r := recover(); r != nil {
								logger.Errorf("执行PanicHandler异常，err:%v", r)
							}
						}()
						panicHandler(r)
					}()
				}
			}
			if nil != afterHandlers && len(afterHandlers) > 0 {
				for _, afterHandler := range afterHandlers {
					func() {
						defer func() {
							if r := recover(); r != nil {
								logger.Errorf("执行AfterHandler异常，err:%v", r)
							}
						}()
						afterHandler()
					}()
				}
			}
		}()

		handler()
	}()
}

/**
执行任务并且自动执行 panicHandler（当抛出 panic的时候才会调用）
*/
func Run(handler func(), panicHandler func(r interface{}), afterHandlers ...func()) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("程序执行任务 Panic：", r)
			if nil != panicHandler {
				func() {
					defer func() {
						if r := recover(); r != nil {
							logger.Errorf("执行PanicHandler异常，err:%v", r)
						}
					}()
					panicHandler(r)
				}()
			}
		}
		if nil != afterHandlers && len(afterHandlers) > 0 {
			for _, afterHandler := range afterHandlers {
				func() {
					defer func() {
						if r := recover(); r != nil {
							logger.Errorf("执行AfterHandler异常，err:%v", r)
						}
					}()
					afterHandler()
				}()
			}
		}
	}()

	handler()
}

/**
获取当前 Goroutine ID
*/
func GetGoroutineId() int {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("获取Goroutine ID 发生 panic recover:panic info:", r)
		}
	}()

	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		return -1
	}
	return id
}

// Goroutine ID --> context
var goidCtxMap = sync.Map{}

/**
将 context 绑定到当前 Goroutine
*/
func BindContext(ctx *context.Context) {
	defer func() {
		if r := recover(); r != nil {
			logger.Warn("绑定Goid和Context异常(panic): ", r)
		}
	}()
	gid := GetGoroutineId()
	goidCtxMap.Store(gid, ctx)
}

/**
获取当前Goroutine 所context
*/
func GetContext() *context.Context {
	gid := GetGoroutineId()
	if val, ok := goidCtxMap.Load(gid); ok {
		if val == nil {
			return nil
		}
		if ctx, y := val.(*context.Context); y {
			return ctx
		}
	}
	return nil
}

/**
将当前 goroutine 所绑定的 context 进行解绑
*/
func UnbindContext() {
	defer func() {
		if r := recover(); r != nil {
			logger.Warn("解绑定Goid和Context异常(panic): ", r)
		}
	}()
	gid := GetGoroutineId()
	goidCtxMap.Delete(gid)
}
