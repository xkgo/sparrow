package ginapp

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/xkgo/sparrow"
	"github.com/xkgo/sparrow/logger"
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
	Engine      *gin.Engine     `@Inject:"required=true"`
	Configures  []Configure     `@Inject:"required=false"`
	Controllers []GinController `@Inject:"required=false"` // 注入所有的 Controller
}

/**
初始化函数
*/
func (r *GinRegistry) Init(app *sparrow.Application) {

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
