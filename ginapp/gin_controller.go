package ginapp

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xkgo/sparrow"
	"github.com/xkgo/sparrow/logger"
	"github.com/xkgo/sparrow/util/ConvertUtils"
	"github.com/xkgo/sparrow/util/GoUtils"
	"github.com/xkgo/sparrow/util/StringUtils"
	"net/http"
	"net/url"
	"strings"
)

// 请求映射
type RequestMapping struct {
	Methods      []string        // 请求方法，支持多个
	Paths        []string        // 请求路径，允许多个
	HandleMethod gin.HandlerFunc // 请求处理方法
}

func (mapping *RequestMapping) Validate() error {
	if len(mapping.Paths) < 1 {
		return errors.New("RequestMapping Paths not found")
	}
	if len(mapping.Methods) < 1 {
		return errors.New("RequestMapping Methods not found")
	}

	if len(mapping.Methods) < 1 {
		return errors.New("RequestMapping Methods not found")
	}

	if mapping.HandleMethod == nil {
		return errors.New("RequestMapping HandleMethod not found")
	}

	return nil
}

func NewRequestMapping(methods, paths string, handleMethod gin.HandlerFunc) RequestMapping {
	rpaths := strings.Split(paths, ",")
	rmethods := strings.Split(strings.ToUpper(methods), ",")
	return RequestMapping{
		Paths:        rpaths,
		Methods:      rmethods,
		HandleMethod: handleMethod,
	}
}

type GinController interface {
	/**
	获取请求映射列表
	*/
	RequestMappings() (mappings []RequestMapping)
}

type Configure interface {
	Config(router *gin.Engine, app *sparrow.Application)
}

type GinRegistry struct {
	Properties  *ServerProperties `@Inject:"required:true"`
	Engine      *gin.Engine       `@Inject:"required=true"`
	Configures  []Configure       `@Inject:"required=false"`
	Controllers []GinController   `@Inject:"required=false"` // 注入所有的 Controller
}

/**
初始化函数
*/
func (r *GinRegistry) Init(app *sparrow.Application) {

	// TraceId and cors
	r.prepareTraceIdAndCors()

	// 初始化配置
	if len(r.Configures) > 0 {
		for _, configure := range r.Configures {
			configure.Config(r.Engine, app)
		}
	}

	logger.Info("准备自动根据Controller注册请求映射，共有Controllers: ", len(r.Controllers))
	if len(r.Controllers) > 0 {
		router := r.Engine
		for _, controller := range r.Controllers {
			mappings := controller.RequestMappings()
			if mappings == nil || len(mappings) < 1 {
				return
			}
			for _, mapping := range mappings {
				if err := mapping.Validate(); err != nil {
					panic(err)
				}
				logger.Debugf("注册请求, Methods:%v, Paths:%v, handler:%v", mapping.Methods, mapping.Paths, mapping.HandleMethod)
				for _, path := range mapping.Paths {
					for _, method := range mapping.Methods {
						router.Handle(method, path, mapping.HandleMethod)
					}
				}
			}
		}
	}
}

/**
TraceId 和 CORS
*/
func (r *GinRegistry) prepareTraceIdAndCors() {
	r.Engine.Use(func(context *gin.Context) {

		defer func() {
			GoUtils.UnbindContext()
		}()

		// trace
		r.prepareTraceId(context)

		// 跨域处理
		if r.prepareCors(context) {
			return
		}

		context.Next()
	})
}

/**
TraceId
*/
func (r *GinRegistry) prepareTraceId(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			_ = fmt.Errorf("prepare trace id panic recover: trace panic info:%v", r)
		}
	}()
	var ctx context.Context = c
	traceIdHeaders := r.Properties.TraceIdHeaders
	traceId := ""
	if nil != traceIdHeaders && len(traceIdHeaders) > 0 {
		for _, traceIdHeader := range traceIdHeaders {
			traceId = c.GetHeader(traceIdHeader)
			if len(traceId) > 0 {
				break
			}
		}
	}
	// 绑定 TraceId
	GoUtils.BindContextWithTraceId(&ctx, traceId)
}

/**
跨域
*/
func (r *GinRegistry) prepareCors(c *gin.Context) (stop bool) {
	defer func() {
		if r := recover(); r != nil {
			_ = fmt.Errorf("prepare cors panic recover: trace panic info:%v", r)
		}
	}()

	// 跨域, 计算 origin
	origin := r.Properties.CorsOrigins
	if len(r.Properties.CorsOrigins) > 0 {
		if r.Properties.CorsOrigins == "*" {
			// 允许所有
			requestOrigin := c.GetHeader("Origin")
			if len(requestOrigin) < 1 {
				referer := c.Request.Referer()
				if len(referer) > 0 {
					refererUrl, e := url.Parse(referer)
					if nil == e {
						requestOrigin = refererUrl.Scheme + "://" + refererUrl.Host
					}
				}
			}
			origin = requestOrigin
		}
	}
	c.Writer.Header().Add("Access-Control-Allow-Origin", origin)

	// 计算 allowHeaders
	allowHeaders := r.Properties.CorsAllowHeaders
	if allowHeaders == "*" {
		// 允许所有的 headers，直接计算请求头中的 header，然后加到这里来
		headers := make([]string, 0)
		for k, _ := range c.Request.Header {
			headers = append(headers, k)
		}
		allowHeaders = strings.Join(headers, ",")
	}
	c.Writer.Header().Add("Access-Control-Allow-Headers", allowHeaders)

	// 计算允许的方法
	allowMethods := r.Properties.CorsAllowMethods
	if allowMethods == "*" {
		// 允许所有的 headers，直接计算请求头中的 header，然后加到这里来
		allowMethods = "POST, GET, OPTIONS, DELETE, PATCH"
	}
	c.Writer.Header().Add("Access-Control-Allow-Methods", allowMethods)

	// 允许 Expose 的请求头
	c.Writer.Header().Add("Access-Control-Expose-Headers", r.mergeHeaders("Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type, X-AuthType, x-requested-with, Access-Token", r.Properties.CorsExposeHeaders))

	c.Writer.Header().Add("Access-Control-Allow-Credentials", "true")

	//放行所有OPTIONS方法
	if c.Request.Method == http.MethodOptions {
		logger.Infof("请求来源, referer:%s, origin:%s", c.Request.Referer(), origin)
		c.Writer.WriteHeader(http.StatusNoContent)
		return true
	}

	return false
}

func (r *GinRegistry) mergeHeaders(headers1, headers2 string) string {
	if len(headers1) < 1 && len(headers2) < 1 {
		return ""
	}
	if len(headers1) > 0 && len(headers2) < 1 {
		return headers1
	}
	if len(headers1) < 1 && len(headers2) > 0 {
		return headers2
	}
	arr1 := ConvertUtils.ToStringBoolMap(StringUtils.SplitByRegex(headers1, ",\\s+"), true)
	arr2 := ConvertUtils.ToStringBoolMap(StringUtils.SplitByRegex(headers2, ",\\s+"), true)

	finalArr := make([]string, 0)
	for v1, _ := range arr1 {
		finalArr = append(finalArr, v1)
	}
	for v2, _ := range arr2 {
		if _, ok := arr1[v2]; !ok {
			finalArr = append(finalArr, v2)
		}
	}
	return strings.Join(finalArr, ", ")
}
