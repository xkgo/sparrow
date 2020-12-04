package env

/*
Interface for resolving properties against any underlying source.
*/
type PropertyResolver interface {

	/**
	 * 判断当前环境中是否配置了 key，注意，给定的 key 不允许为空字符串
	 */
	ContainsProperty(key string) bool

	/**
	获取配置项的字符串值，返回的值中，====不包含占位符====
	@param key 配置 key
	@return value 对应配置项的值
	@return exists 配置项是否存在，即 @ContainsProperty(key string) 的返回值一样意义
	*/
	GetProperty(key string) (value string, exists bool)

	/**
	获取指定配置项的值，如果对应配置项没有配置，那么返回 默认值，====不包含占位符====
	*/
	GetPropertyWithDef(key string, def string) string

	/**
	获取配置项，如果配置项没有配置的话，那么会直接 panic（通常是在项目启动的时候使用，避免项目非正常状态下启动），====不包含占位符====
	*/
	GetRequiredProperty(key string) string

	/**
	处理类似 ${...} 这种占位符， 替换对应的配置项，如果 ${...}中的配置项不存在，则不进行替换，但是也不会抛异常
	*/
	ResolvePlaceholders(text string) string

	/**
	处理类似 ${...} 这种占位符， 替换对应的配置项，如果 ${...}中的配置项不存在，则直接 panic，这是为了防止非正常启动
	*/
	ResolveRequiredPlaceholders(text string) string
}
