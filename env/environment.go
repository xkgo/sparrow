package env

import "github.com/xkgo/sparrow/deploy"

/*
环境，包含运行环境、启动项目用到的参数、配置等等
*/
type Environment interface {
	PropertyResolver

	/**
	获取激活的配置列表，返回文件的绝对路径
	*/
	GetActiveProfiles() []string

	/**
	获取一个可变的属性来源对象
	*/
	GetPropertySources() *MutablePropertySources

	/**
	合并父环境信息，子环境属性优先生效，只有子环境中不存在的才会在父环境中使用，比如假设父子环境中都有相同的配置key，那么将会使用子环境的优先
	*/
	Merge(parent Environment)

	/**
	订阅变更, keyPattern: 等值、正则匹配，如果为空字符串或者 * 那么表示所有，如果是个合法的正则，那么就按照正则匹配
	*/
	Subscribe(keyPattern string, handler func(event *KeyChangeEvent))

	/**
	绑定配置项到某个模型对象，注意传进来的必须是指针类型, keyPrefix key前缀，会直接和配置struct的属性直接拼接，如果有.的话要注意了
	@param name 名称，唯一
	@param cfgPtr 配置指针
	@param changedListen 是否需要进行监听
	*/
	BindProperties(keyPrefix string, cfgPtr interface{}) (beanPtr interface{}, err error)

	BindPropertiesListen(keyPrefix string, cfgPtr interface{}, changedListen bool) (beanPtr interface{}, err error)

	/**
	获取属性对象
	@param typeTemplate 类型模板，可以提供属性类型的指针类型，也可以直接提供 reflect.Type 类型
	*/
	GetProperties(typeTemplate interface{}) (beanPtr interface{})

	IsDev() bool
	IsTest() bool
	IsFat() bool
	IsProd() bool
	/**
	当前环境
	*/
	GetEnv() deploy.Env
	/**
	部署集合
	*/
	GetSet() string
}
