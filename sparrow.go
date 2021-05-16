package sparrow

import (
	"github.com/xkgo/sparrow/env"
	"github.com/xkgo/sparrow/logger"
	"reflect"
)

/**
注册Bean, 允许多个别名
@param beanPtr 要注册的bean指针对象
@param beanName bean的名称
@param primary 是否是主bean
*/
func RegisterBean(beanPtr interface{}, beanName string, primary bool) {
	getApp().RegisterBean(beanPtr, beanName, primary)
}

/**
注册 PropertiesBean, 允许多个别名
@param beanPtr 要注册的bean指针对象
@param beanName bean的名称
@param keyPrefix key前缀，会直接拼接，如果注意有必要的话需要加 .
@param primary 是否是主bean
*/
func RegisterPropertiesBean(beanPtr interface{}, beanName string, keyPrefix string, primary bool) {
	getApp().RegisterPropertiesBean(beanPtr, beanName, keyPrefix, primary)
}

func RegisterPropertiesBeanListen(beanPtr interface{}, beanName string, keyPrefix string, changedListen, primary bool) {
	getApp().RegisterPropertiesBeanListen(beanPtr, beanName, keyPrefix, changedListen, primary)
}

/**
获取指定名称的Bean，不存在则返回 nil
*/
func GetBeanByName(beanName string) (beanPtr interface{}) {
	return getApp().GetBeanByName(beanName)
}

/**
获取指定类型的 beans， key 为beanName， value 为对应的对象指针
*/
func GetBeansOfType(beanType reflect.Type) (beansPtr map[string]interface{}) {
	return getApp().GetBeansOfType(beanType)
}

func GetBeanByType(beanTemplate interface{}) (beanPtr interface{}, err error) {
	return getApp().GetBeanByType(beanTemplate)
}

/**
获取所有的 BeanNames
*/
func GetBeanNames() []string {
	return getApp().GetBeanNames()
}

var app *Application

func getApp() *Application {
	if app == nil {
		app = NewApplication()
	}
	return app
}

func AppendBeforeInitHandler(handler Runner) {
	getApp().AppendBeforeInitHandler(handler)
}

/**
运行程序
*/
func Run(environment env.Environment, options ...Option) (err error) {
	err = getApp().Run(environment, options...)
	if nil != err {
		logger.Error("程序运行异常退出: ", err)
	}
	return err
}
