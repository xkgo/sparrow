package GoUtils

import (
	"context"
	"github.com/xkgo/sparrow/logger"
	"github.com/xkgo/sparrow/util"
	"github.com/xkgo/sparrow/util/RandomUtils"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

func init() {
	logger.SetTraceIdGenerator(func(ctx *context.Context) string {
		return GetTraceIdWithContext(ctx)
	})
}

// TraceId 获取，如果context是其他类型的话，允许自定义
type TraceIdGetter func(ctx *context.Context) string

var traceIdGetter TraceIdGetter

func RegisterTraceIdGetter(getter TraceIdGetter) {
	if nil != traceIdGetter {
		traceIdGetter = func(ctx *context.Context) string {
			traceId := getter(ctx)
			if len(traceId) > 0 {
				return traceId
			}
			return traceIdGetter(ctx)
		}
	}
	traceIdGetter = getter
}

/**
使用 Goroutine 执行，并且自动进行panic异常cache
*/
func RunGoroutine(handler func(), panicHandler func(r interface{}), afterHandlers ...func()) {
	ctx, traceId := BindContext(GetContext())
	go func() {
		BindContextWithTraceId(ctx, traceId)
		defer func() {
			UnbindContext()
		}()
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

type invokeContext struct {
	traceId   string
	dataPairs sync.Map
}

/**
将 context 绑定到当前 Goroutine
@return bindCtx 实际绑定的 context
@return bindTraceId 实际绑定的 traceId
@return dataPairs 数据键值对
*/
func BindContext(ctx *context.Context, dataPairs ...*util.Pair) (bindCtx *context.Context, bindTraceId string) {
	return BindContextWithTraceId(ctx, "", dataPairs...)
}

/**
将 context,traceId 绑定到当前 Goroutine
@return bindCtx 实际绑定的 context
@return bindTraceId 实际绑定的 traceId
@return dataPairs 数据键值对
*/
func BindContextWithTraceId(ctx *context.Context, traceId string, dataPairs ...*util.Pair) (bindCtx *context.Context, bindTraceId string) {
	defer func() {
		if r := recover(); r != nil {
			logger.Warn("绑定Goid和Context异常(panic): ", r)
		}
	}()
	var iContext *invokeContext

	if ctx == nil {
		bctx := context.Background()
		iContext = &invokeContext{dataPairs: sync.Map{}}
		nctx := context.WithValue(bctx, "iContext", iContext)
		ctx = &nctx
	} else {
		if invoke, ok := (*ctx).Value("iContext").(*invokeContext); ok {
			iContext = invoke
		} else {
			iContext = &invokeContext{dataPairs: sync.Map{}}
			nctx := context.WithValue(*ctx, "iContext", iContext)
			ctx = &nctx
		}
	}

	if len(traceId) < 1 {
		traceId = iContext.traceId
		if len(traceId) < 1 {
			// 县尝试从 context 中获取，看看是否已经绑定过了
			// 尝试从 context 中解析出来，没有的话就直接创建一个新的
			if traceIdGetter != nil {
				traceId = traceIdGetter(ctx)
			}
			if len(traceId) < 1 {
				traceId = RandomUtils.RandomLetterAndNumberString(10)
			}
		}
	}
	if len(traceId) > 0 {
		iContext.traceId = traceId
	}

	// 数据
	if dataPairs != nil && len(dataPairs) > 0 {
		for _, pair := range dataPairs {
			iContext.dataPairs.Store(pair.K, pair.V)
		}
	}

	gid := GetGoroutineId()
	goidCtxMap.Store(gid, ctx)
	return ctx, iContext.traceId
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
获取当前协程 绑定的 traceId
*/
func GetTraceId() string {
	return GetTraceIdWithContext(nil)
}

/**
获取当前协程 绑定的 traceId
*/
func GetTraceIdWithContext(ctx *context.Context) string {
	if nil == ctx {
		ctx = GetContext()
	}
	if nil == ctx {
		_, traceId := BindContextWithTraceId(nil, "")

		return traceId
	}
	if invoke, ok := (*ctx).Value("iContext").(*invokeContext); ok {
		return invoke.traceId
	}
	return ""
}

/**
获取当前协程 绑定的某个 key 的值
*/
func GetData(key interface{}) (value interface{}, exists bool) {
	ctx := GetContext()
	if nil == ctx {
		return nil, false
	}
	if invoke, ok := (*ctx).Value("iContext").(*invokeContext); ok {
		return invoke.dataPairs.Load(key)
	}
	return nil, false
}

/**
给当前协程绑定数据
*/
func SetData(key interface{}, value interface{}) {
	ctx := GetContext()
	if nil == ctx {
		_, _ = BindContext(nil, util.NewPair(key, value))
		return
	}
	if invoke, ok := (*ctx).Value("iContext").(*invokeContext); ok {
		invoke.dataPairs.Load(key)
	}
	return
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
