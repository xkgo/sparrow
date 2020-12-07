package env

// 选项
type Option func(environment *StandardEnvironment)

/**
配置选项
*/
type Options struct {
	// profile 所在目录列表
	profileDirs []string

	/**
	附加配置来源，默认会添加到环境变量之前, 一般是要接入额外的配置心
	*/
	additionalPropertySources *MutablePropertySources

	/**
	是否忽略无法处理的占位符，如果忽略则不处理，不忽略的话，那么遇到不能解析的占位符直接 panic
	*/
	ignoreUnresolvableNestedPlaceholders bool
}

/**
配置文件扫描路径，会按照顺序依次搜索配置文件，并且搜索的顺序就是生效的顺序
*/
func ConfigDirs(dirs ...string) Option {
	return func(environment *StandardEnvironment) {
		// 直接进行覆盖
		environment.options.profileDirs = dirs
	}
}

/**
添加额外的配置来源
*/
func AdditionalPropertySources(additionalPropertySources *MutablePropertySources) Option {
	return func(environment *StandardEnvironment) {
		environment.options.additionalPropertySources = additionalPropertySources
	}
}

/**
是否忽略无法处理的占位符，如果忽略则不处理，不忽略的话，那么遇到不能解析的占位符直接 panic
*/
func IgnoreUnresolvableNestedPlaceholders(ignore bool) Option {
	return func(environment *StandardEnvironment) {
		environment.ignoreUnresolvableNestedPlaceholders = ignore
	}
}
