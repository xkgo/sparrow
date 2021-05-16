package PlaceholderUtils

import (
	"github.com/xkgo/sparrow/env"
	"sync"
)

// Goroutine ID --> context
var helpers = sync.Map{}

func init() {
	helpers.Store("${_:_}", env.NewPropertyPlaceholderHelper("${", "}", ":", true))
}

func getHelper(prefix, suffix, valueSeparator string) *env.PropertyPlaceholderHelper {
	if len(prefix) < 1 {
		prefix = "${"
	}
	if len(suffix) < 1 {
		suffix = "}"
	}
	if len(valueSeparator) < 1 {
		valueSeparator = ":"
	}

	key := prefix + "_" + valueSeparator + "_" + suffix
	if helper, ok := helpers.Load(key); ok {
		return helper.(*env.PropertyPlaceholderHelper)
	}

	helper := env.NewPropertyPlaceholderHelper(prefix, suffix, valueSeparator, true)
	helpers.Store(key, helper)
	return helper
}

/**
替换占位符工具类
@param text	string 要处理占位符的字符串
@param params	map[string]string 占位符的参数
@param placeholderPrefix string 占位符前缀，默认是 "${"
@param placeholderSuffix string 占位符后缀，默认是 "}"
@param valueSeparator string 占位符内值切割字符串，默认是 ":", 比如 ${name:${def-name:arvin}}
*/
func ResolvePlaceholdersExt(
	text string,
	params map[string]string,
	placeholderPrefix string,
	placeholderSuffix string,
	valueSeparator string) (value string) {

	placeholderResolver := func(key string) string {
		if params == nil {
			return ""
		}
		return params[key]
	}
	return getHelper(placeholderPrefix, placeholderSuffix, valueSeparator).ReplacePlaceholders(text, placeholderResolver)
}

/**
替换占位符工具类
@param text	string 要处理占位符的字符串
@param params	map[string]string 占位符的参数
*/
func ResolvePlaceholders(text string, params map[string]string) (value string) {
	return ResolvePlaceholdersExt(text, params, "${", "}", ":")
}
