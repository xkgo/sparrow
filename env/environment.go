package env

/*
环境，包含运行环境、启动项目用到的参数、配置等等
*/
type Environment interface {
	PropertyResolver

	/**
	获取激活的配置列表，返回文件的绝对路径， TODO 确定指定配置文件的方式
	*/
	GetActiveProfiles() []string

	/**
	添加一个或多个配置文件
	*/
	AddActiveProfiles(profiles ...string)

	/**
	覆盖设置配置文件
	*/
	SetActiveProfiles(profiles ...string)

	/**
	获取一个可变的属性来源对象
	*/
	GetPropertySources() *MutablePropertySources

	/**
	合并父环境信息，子环境属性优先生效，只有子环境中不存在的才会在父环境中使用，比如假设父子环境中都有相同的配置key，那么将会使用子环境的优先
	*/
	Merge(parent Environment)
}
