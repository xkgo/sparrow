package env

// 属性变更类型
type KeyChangeType string

const (
	PropertyAdd    KeyChangeType = "ADD"
	PropertyUpdate KeyChangeType = "UPDATE"
	PropertyDel    KeyChangeType = "DEL"
)

// 配置变更事件
type KeyChangeEvent struct {
	key        string        // 变更的配置 key
	ov         string        // 旧值
	nv         string        // 新值
	changeType KeyChangeType // 变更类型
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
	订阅 Key 变更, 名字唯一
	@param name 给这个监听器命名，一般标识是谁在监听
	@return queue 如果支持，就会返回一个 channel，当配置变更的时候会发送事件过去
	@return support 表示是否支持订阅，一般如果配置来源不会变更的话，就不存在说订阅的问题了
	如果订阅异常，那么会直接 panic
	*/
	SubscribeKeyChange(name string) (queue chan *KeyChangeEvent, support bool)

	/**
	取消订阅key变更
	*/
	UnsubscribeKeyChange(name string)
}
