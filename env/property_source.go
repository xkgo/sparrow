package env

import (
	"fmt"
	"regexp"
)

// 属性变更类型
type KeyChangeType string

const (
	PropertyAdd    KeyChangeType = "ADD"
	PropertyUpdate KeyChangeType = "UPDATE"
	PropertyDel    KeyChangeType = "DEL"
)

// 配置变更事件
type KeyChangeEvent struct {
	Key        string        // 变更的配置 key
	Ov         string        // 旧值
	Nv         string        // 新值
	ChangeType KeyChangeType // 变更类型
}

func (e *KeyChangeEvent) String() string {
	return fmt.Sprintf("[%v][%s], old:[%s], new:[%s]", e.ChangeType, e.Key, e.Ov, e.Nv)
}

type PropertyChangeListener struct {
	KeyPattern string                      // 等值、正则匹配，如果为空字符串或者 * 那么表示所有，如果是个合法的正则，那么就按照正则匹配
	Regex      *regexp.Regexp              // 正则表达式
	Handler    func(event *KeyChangeEvent) // 处理器
}

func NewPropertyChangeListener(keyPattern string, handler func(event *KeyChangeEvent)) *PropertyChangeListener {
	var regex *regexp.Regexp
	if keyPattern != "" && keyPattern != "*" {
		regex, _ = regexp.Compile(keyPattern)
	}
	return &PropertyChangeListener{
		KeyPattern: keyPattern,
		Regex:      regex,
		Handler:    handler,
	}
}

type PropertySource interface {
	/**
	配置源名称
	*/
	GetName() string

	/**
	获取配置项的字符串值，返回的值中，包含占位符
	@param key 配置 key
	@return value 对应配置项的值
	@return exists 配置项是否存在，即 @ContainsProperty(key string) 的返回值一样意义
	*/
	GetProperty(key string) (value string, exists bool)

	/**
	获取指定配置项的值，如果对应配置项没有配置，那么返回 默认值
	*/
	GetPropertyWithDef(key string, def string) string

	/**
	遍历所有的配置项&值, consumer 处理过程中如果返回 stop=true则停止遍历
	*/
	Each(consumer func(key, value string) (stop bool))

	/**
	订阅变更, keyPattern: 等值、正则匹配，如果为空字符串或者 * 那么表示所有，如果是个合法的正则，那么就按照正则匹配
	*/
	Subscribe(keyPattern string, handler func(event *KeyChangeEvent))
}
