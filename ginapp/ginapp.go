package ginapp

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xkgo/sparrow"
	"github.com/xkgo/sparrow/env"
	"github.com/xkgo/sparrow/logger"
)

type ServerProperties struct {
	Port              int      `ck:"port" def:"8088"`
	TraceIdHeaders    []string `ck:"trace-id-headers" def:"[\"trace-id\",\"traceid\"]"` // 读取哪些header字段作为traceId
	CorsOrigins       string   `ck:"cors-origins" def:""`                               // 使用英文逗号分隔
	CorsAllowHeaders  string   `ck:"cors-allow-headers" def:""`                         // 使用英文逗号分隔
	CorsAllowMethods  string   `ck:"cors-allow-methods" def:""`                         // 使用英文逗号分隔
	CorsExposeHeaders string   `ck:"cors-expose-headers" def:""`                        // 使用英文逗号分隔
}

type GinServer struct {
	Router     *gin.Engine      `@Inject:"required=true"`
	Properties ServerProperties `@Inject:"required:true"`
}

func (g *GinServer) Run(app *sparrow.Application) (err error) {
	return g.Router.Run(fmt.Sprintf(":%v", g.Properties.Port))
}

/**
运行程序
*/
func Run(environment env.Environment, options ...sparrow.Option) (err error) {
	// 服务器配置
	sparrow.RegisterPropertiesBean(&ServerProperties{}, "", "server.", true)

	// 控制器注册
	sparrow.RegisterBean(&GinRegistry{}, "gin_GinRegistry", true)

	sparrow.RegisterBean(gin.New(), "gin_Engine", true)

	ginServer := &GinServer{}
	sparrow.RegisterBean(ginServer, "", true)

	if options == nil {
		options = make([]sparrow.Option, 0)
	}
	// 启动
	options = append(options, sparrow.WithRunner(ginServer.Run))

	err = sparrow.Run(environment, options...)
	if nil != err {
		logger.Error("程序运行异常退出: ", err)
	}
	return err
}
