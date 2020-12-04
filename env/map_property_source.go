package env

import "github.com/xkgo/sparrow/logger"

/**
基于 Map 实现的 env/PropertySource 接口
*/
type MapPropertySource struct {
	name       string            // 给这个命个名
	properties map[string]string // 配置map
}

func NewMapPropertySource(name string, properties map[string]string) *MapPropertySource {
	return &MapPropertySource{
		name:       name,
		properties: properties,
	}
}

func (m *MapPropertySource) GetName() string {
	return m.name
}

func (m *MapPropertySource) GetProperty(key string) (value string, exists bool) {
	if nil == m.properties {
		m.properties = make(map[string]string)
		return "", false
	}
	if len(m.properties) < 1 {
		return "", false
	}
	value, exists = m.properties[key]
	return
}

func (m *MapPropertySource) GetPropertyWithDef(key string, def string) string {
	if value, exists := m.GetProperty(key); exists {
		return value
	}
	return def
}

func (m *MapPropertySource) Each(consumer func(key string, value string) (stop bool)) {
	if consumer == nil || nil == m.properties || len(m.properties) < 1 {
		return
	}
	for k, v := range m.properties {
		if consumer(k, v) {
			return
		}
	}
}

func (m *MapPropertySource) SubscribeKeyChange(name string) (queue chan *KeyChangeEvent, support bool) {
	logger.Error("Invalid key change event subscribe for MapPropertySource from subscriber:" + name)
	return nil, false
}

func (m *MapPropertySource) UnsubscribeKeyChange(name string) {
	// map 不可变，不需要处理
}
