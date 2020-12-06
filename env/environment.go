package env

/*
环境，包含运行环境、启动项目用到的参数、配置等等
*/
type Environment interface {
	PropertyResolver

	/**
	获取激活的配置列表，返回文件的绝对路径， TODO 确定指定配置文件的方式
	*/
	GetActiveProfiles() []string

	/**
	获取一个可变的属性来源对象
	*/
	GetPropertySources() *MutablePropertySources

	/**
	合并父环境信息，子环境属性优先生效，只有子环境中不存在的才会在父环境中使用，比如假设父子环境中都有相同的配置key，那么将会使用子环境的优先
	*/
	Merge(parent Environment)

	/**
	订阅 Key 变更, 名字唯一
	@param consumer 给这个监听器命名，一般标识是谁在监听
	@param keyPattern 正则，只要匹配这个pattern 的配置项发生了变更，那么就发布事件
	@return queue 如果支持，就会返回一个 channel，当配置变更的时候会发送事件过去
	如果订阅异常或者重复订阅的，那么会直接 panic
	*/
	SubscribeKeyChange(consumer, keyPattern string) (queue chan *KeyChangeEvent)

	/**
	取消订阅keyPattern变更
	*/
	UnsubscribeKeyChange(consumer, keyPattern string)
}
