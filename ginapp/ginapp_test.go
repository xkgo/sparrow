package ginapp

import (
	"github.com/gin-gonic/gin"
	"github.com/xkgo/sparrow"
	"github.com/xkgo/sparrow/env"
	"github.com/xkgo/sparrow/logger"
	"net/http"
	"testing"
)

type HelloController struct {
}

func (h *HelloController) RequestMappings() (mappings []RequestMapping) {
	return []RequestMapping{
		NewRequestMapping("GET", "/hello/sayHello", h.sayHello),
	}
}

func (h *HelloController) sayHello(c *gin.Context) {
	logger.Infof("SayHello请求：", c.Query("name"))
	c.JSON(http.StatusOK, "Hello, "+c.Query("name"))
}

type Middleware struct {
}

func (m *Middleware) Config(router *gin.Engine, app *sparrow.Application) {
	router.Use(func(context *gin.Context) {

		logger.Info(context.Request.Method + ": " + context.Request.RequestURI)

		context.Next()
	})
}

func TestGinrunner(t *testing.T) {
	sparrow.RegisterBean(&HelloController{}, "", true)
	sparrow.RegisterBean(&Middleware{}, "", true)

	err := Run(env.New())
	if err != nil {
		logger.Error("系统退出：", err)
	}
}
