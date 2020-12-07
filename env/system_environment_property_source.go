package env

import (
	"github.com/xkgo/sparrow/util/StringUtils"
	"os"
)

const (
	/** 系统环境变量 PropertySource GetName */
	SystemEnvironmentPropertySourceName = "systemEnvironment"
)

/**
系统环境变量 属性来源
*/
type SystemEnvironmentPropertySource struct {
	MapPropertySource
}

func NewSystemEnvironmentPropertySource() *SystemEnvironmentPropertySource {
	source := &SystemEnvironmentPropertySource{}
	source.name = SystemEnvironmentPropertySourceName
	source.properties = make(map[string]string)

	envs := os.Environ()
	for _, kv := range envs {
		kvs := StringUtils.SplitByRegex(kv, "\\s*=\\s*")
		if len(kvs) == 1 {
			source.properties[StringUtils.Trim(kvs[0])] = ""
		} else if len(kvs) == 2 {
			source.properties[StringUtils.Trim(kvs[0])] = StringUtils.Trim(kvs[1])
		}
	}
	return source
}
