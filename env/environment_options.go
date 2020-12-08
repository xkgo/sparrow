package env

import "github.com/xkgo/sparrow/deploy"

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

	/**
	自定义部署信息，如果设置了这个，那么直接以这个为准
	*/
	customDeployInfo *deploy.Info

	/**
	附加的命令行参数，如果原来有命令参数，再继续追加的话，原来相同key的会被覆盖
	*/
	appendCommandLine string

	/**
	追加的profiles，会放到 原来的之后
	*/
	appendProfiles []string
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

/**
自定义部署信息，如果这个设置了的话，直接使用这个部署信息，而不是用检测的，一般用于本地测试
*/
func DeployInfo(info *deploy.Info) Option {
	return func(environment *StandardEnvironment) {
		environment.options.customDeployInfo = info
	}
}

/**
追加命令行参数，如果原来有命令参数，再继续追加的话，原来相同key的会被覆盖
*/
func AppendCommandLine(commandLine string) Option {
	return func(environment *StandardEnvironment) {
		environment.options.appendCommandLine = environment.options.appendCommandLine + " " + commandLine
	}
}

/**
包含 profiles, 会放在系统计算的之后
*/
func IncludeProfiles(profiles ...string) Option {
	return func(environment *StandardEnvironment) {
		if nil == environment.options.appendProfiles {
			environment.options.appendProfiles = make([]string, 0)
		}
		environment.options.appendProfiles = append(environment.options.appendProfiles, profiles...)
	}
}
