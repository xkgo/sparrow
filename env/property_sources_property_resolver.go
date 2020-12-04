package env

import (
	"github.com/xkgo/sparrow/logger"
)

/**
配置来源属性解析器, 实现接口：env/PropertyResolver
*/
type PropertySourcesPropertyResolver struct {
	propertySources                      PropertySources            // 配置来源
	ignoreUnresolvableNestedPlaceholders bool                       // 是否忽略无法处理的占位符，如果忽略则不处理，不忽略的话，那么遇到不能解析的占位符直接 panic
	nonStrictHelper                      *PropertyPlaceholderHelper // 当遇到未定义的配置项时，不进行替换，也不会抛出异常
	strictHelper                         *PropertyPlaceholderHelper // 当遇到未定义的配置项时，直接 panic
}

/**
创建属性解析器
@param ignoreUnresolvableNestedPlaceholders 是否忽略无法处理的占位符，如果忽略则不处理，不忽略的话，那么遇到不能解析的占位符直接 panic
*/
func NewPropertySourcesPropertyResolver(propertySources PropertySources, ignoreUnresolvableNestedPlaceholders bool) *PropertySourcesPropertyResolver {
	return &PropertySourcesPropertyResolver{
		propertySources:                      propertySources,
		ignoreUnresolvableNestedPlaceholders: ignoreUnresolvableNestedPlaceholders,
	}
}

func (p *PropertySourcesPropertyResolver) ContainsProperty(key string) bool {
	if nil == p.propertySources {
		return false
	}
	contains := false
	p.propertySources.Each(func(index int, source PropertySource) (stop bool) {
		if _, ok := source.GetProperty(key); ok {
			contains = true
			return true
		}
		return false
	})
	return contains
}

/**
获取配置项
@param resolveNestedPlaceholders 是否需要处理占位符
*/
func (p *PropertySourcesPropertyResolver) doGetProperty(key string, resolveNestedPlaceholders bool) (value string, exists bool) {
	if nil == p.propertySources {
		return "", false
	}
	p.propertySources.Each(func(index int, source PropertySource) (stop bool) {
		if logger.IsDebugEnabled() {
			logger.Debug("Searching for key '" + key + "' in PropertySource '" + source.GetName() + "'")
		}
		if val, ok := source.GetProperty(key); ok {
			exists = true
			value = val

			// 找到了key，加下日志
			if logger.IsDebugEnabled() {
				logger.Debug("Found key '" + key + "' in PropertySource '" + source.GetName() + "' with value: " + value)
			}

			// 看看是否需要替换占位符, ${...}, 长度至少是4 才能构成一个占位符
			if resolveNestedPlaceholders && len(value) > 4 {
				value = p.resolveNestedPlaceholders(value)
			}
			return true
		}
		return false
	})
	return
}

func (p *PropertySourcesPropertyResolver) GetProperty(key string) (value string, exists bool) {
	return p.doGetProperty(key, true)
}

func (p *PropertySourcesPropertyResolver) GetPropertyWithDef(key string, def string) string {
	if value, exists := p.doGetProperty(key, true); exists {
		return value
	}
	return def
}

func (p *PropertySourcesPropertyResolver) GetRequiredProperty(key string) string {
	if value, exists := p.doGetProperty(key, true); exists {
		return value
	}
	panic("Required key '" + key + "' not found")
}

func (p *PropertySourcesPropertyResolver) ResolvePlaceholders(text string) string {
	if p.nonStrictHelper == nil {
		p.nonStrictHelper = p.createPlaceholderHelper(true)
	}
	return p.doResolvePlaceholders(text, p.nonStrictHelper)
}

func (p *PropertySourcesPropertyResolver) ResolveRequiredPlaceholders(text string) string {
	if p.strictHelper == nil {
		p.strictHelper = p.createPlaceholderHelper(false)
	}
	return p.doResolvePlaceholders(text, p.strictHelper)
}

/**
处理占位符，将占位符为 ${...} 替换掉
*/
func (p *PropertySourcesPropertyResolver) resolveNestedPlaceholders(text string) string {
	if p.ignoreUnresolvableNestedPlaceholders {
		return p.ResolvePlaceholders(text)
	} else {
		return p.ResolveRequiredPlaceholders(text)
	}
}

func (p *PropertySourcesPropertyResolver) createPlaceholderHelper(ignoreUnresolvablePlaceholders bool) *PropertyPlaceholderHelper {
	return NewPropertyPlaceholderHelper(DefaultPlaceholderPrefix, DefaultPlaceholderSuffix, DefaultPlaceholderValueSeparator, ignoreUnresolvablePlaceholders)
}

func (p *PropertySourcesPropertyResolver) getPropertyAsRawString(key string) string {
	return p.GetPropertyWithDef(key, "")
}

func (p *PropertySourcesPropertyResolver) doResolvePlaceholders(text string, helper *PropertyPlaceholderHelper) string {
	return helper.ReplacePlaceholders(text, p.getPropertyAsRawString)
}
