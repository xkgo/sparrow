package env

import (
	"github.com/xkgo/sparrow/logger"
	"github.com/xkgo/sparrow/util/StringUtils"
	"strings"
)

const (
	DefaultPlaceholderPrefix         = "${" // 默认占位符前缀
	DefaultPlaceholderSuffix         = "}"  // 默认占位符后缀
	DefaultPlaceholderValueSeparator = ":"  // 占位符内值分隔符
)

var (
	WellKnownSimplePrefixes = map[string]string{
		"}": "{",
		"]": "[",
		")": "(",
	}
)

/**
占位符处理Helper
*/
type PropertyPlaceholderHelper struct {
	placeholderPrefix              string // 占位符前缀，如 ${
	placeholderSuffix              string // 占位符后缀，如 }
	simplePrefix                   string // 简单前缀
	valueSeparator                 string // 值分隔符
	ignoreUnresolvablePlaceholders bool   // 是否忽略不识别的占位符，如果为 true 的话，当发现未识别的占位符的时候，直接 panic
	minValueLen                    int    // 要进行占位符处理的值最小长度
}

func (h *PropertyPlaceholderHelper) ReplacePlaceholders(value string, placeholderResolver func(key string) string) string {
	if len(value) < h.minValueLen {
		// 不需要进行处理， 加起来还没有占位符长
		return value
	}
	value, _ = h.parseStringValue(value, nil, placeholderResolver)
	return value
}

func (h *PropertyPlaceholderHelper) parseStringValue(value string, visitedPlaceholders map[string]bool, placeholderResolver func(key string) string) (val string, visited map[string]bool) {
	if visitedPlaceholders == nil {
		visitedPlaceholders = make(map[string]bool)
	}
	result := value
	startIndex := strings.Index(value, h.placeholderPrefix)
	for startIndex != -1 {
		endIndex := h.findPlaceholderEndIndex(result, startIndex)
		if endIndex != -1 {
			placeholder := result[startIndex+len(h.placeholderPrefix) : endIndex] // 注意，endIndex 是不包含进来的
			originalPlaceholder := placeholder
			if _, dup := visitedPlaceholders[originalPlaceholder]; dup {
				errMsg := "Circular placeholder reference '" + originalPlaceholder + "' in property definitions"
				logger.Fatal(errMsg)
				panic(errMsg)
			}
			visitedPlaceholders[originalPlaceholder] = true
			// 递归处理
			placeholder, visitedPlaceholders = h.parseStringValue(placeholder, visitedPlaceholders, placeholderResolver)
			propVal := placeholderResolver(placeholder)
			useDefault := false
			if len(propVal) == 0 && len(h.valueSeparator) > 0 {
				separatorIndex := strings.Index(placeholder, h.valueSeparator)
				if separatorIndex != -1 {
					actualPlaceholder := placeholder[0:separatorIndex]
					defaultValue := placeholder[separatorIndex+len(h.valueSeparator):]
					propVal = placeholderResolver(actualPlaceholder)
					if len(propVal) < 1 {
						propVal = defaultValue
						useDefault = true
					}
				}
			}
			if len(propVal) > 0 || useDefault {
				if !useDefault {
					propVal, visitedPlaceholders = h.parseStringValue(propVal, visitedPlaceholders, placeholderResolver)
				}
				// 将解析出来的值进行替换
				result, _ = StringUtils.ReplaceRange(result, propVal, startIndex, endIndex+len(h.placeholderSuffix))
				startIndex, _ = StringUtils.IndexFrom(result, h.placeholderPrefix, startIndex+len(propVal))
			} else if h.ignoreUnresolvablePlaceholders {
				// Proceed with unprocessed value.
				startIndex, _ = StringUtils.IndexFrom(result, h.placeholderPrefix, endIndex+len(h.placeholderSuffix))
			} else {
				errMsg := "Could not resolve placeholder '" + placeholder + "'" + " in value \"" + value + "\""
				logger.Fatal(errMsg)
				panic(errMsg)
			}
			delete(visitedPlaceholders, originalPlaceholder)
		} else {
			startIndex = -1
		}
	}
	return result, visitedPlaceholders
}

/**
从 startIndex 开始，查找对应的结尾符号
*/
func (h *PropertyPlaceholderHelper) findPlaceholderEndIndex(value string, startIndex int) int {
	index := startIndex + len(h.placeholderPrefix)
	withinNestedPlaceholder := 0
	for index < len(value) {
		if StringUtils.SubstringMatch(value, index, h.placeholderSuffix) {
			if withinNestedPlaceholder > 0 {
				withinNestedPlaceholder--
				index = index + len(h.placeholderSuffix)
			} else {
				return index
			}
		} else if StringUtils.SubstringMatch(value, index, h.simplePrefix) {
			withinNestedPlaceholder++
			index = index + len(h.simplePrefix)
		} else {
			index++
		}
	}
	return -1
}

func NewPropertyPlaceholderHelper(placeholderPrefix string,
	placeholderSuffix string,
	valueSeparator string,
	ignoreUnresolvablePlaceholders bool) *PropertyPlaceholderHelper {

	if len(placeholderPrefix) < 1 {
		placeholderPrefix = DefaultPlaceholderPrefix
	}
	if len(placeholderSuffix) < 1 {
		placeholderSuffix = DefaultPlaceholderSuffix
	}

	if len(valueSeparator) < 1 {
		valueSeparator = DefaultPlaceholderValueSeparator
	}

	simplePrefix := placeholderPrefix
	if simplePrefixForSuffix, ok := WellKnownSimplePrefixes[placeholderSuffix]; ok && len(simplePrefixForSuffix) > 0 {
		simplePrefix = simplePrefixForSuffix
	}

	helper := &PropertyPlaceholderHelper{
		placeholderPrefix:              placeholderPrefix,
		placeholderSuffix:              placeholderSuffix,
		simplePrefix:                   simplePrefix,
		valueSeparator:                 valueSeparator,
		ignoreUnresolvablePlaceholders: ignoreUnresolvablePlaceholders,
		minValueLen:                    len(placeholderPrefix) + len(placeholderSuffix) + 1,
	}

	return helper
}
