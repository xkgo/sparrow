package env

import "os"

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

	keys := os.Environ()
	for _, key := range keys {
		source.properties[key] = os.Getenv(key)
	}
	return source
}
