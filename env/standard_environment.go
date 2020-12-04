package env

/**
标准环境实现, 实现接口 Environment
*/
type StandardEnvironment struct {
	activeProfiles                       []string                // 启用的配置文件列表，配置文件绝对路径
	ignoreUnresolvableNestedPlaceholders bool                    // 是否忽略无法处理的占位符，如果忽略则不处理，不忽略的话，那么遇到不能解析的占位符直接 panic
	propertySources                      *MutablePropertySources // 配置来源
	propertyResolver                     PropertyResolver        // 配置处理器
}

func (s *StandardEnvironment) InitPropertyResolver() {
	if s.propertySources == nil || s.propertyResolver == nil {
		s.propertyResolver = &PropertySourcesPropertyResolver{
			propertySources:                      s.propertySources,
			ignoreUnresolvableNestedPlaceholders: s.ignoreUnresolvableNestedPlaceholders,
		}
	}
}

func (s *StandardEnvironment) ContainsProperty(key string) bool {
	s.InitPropertyResolver()
	return s.propertyResolver.ContainsProperty(key)
}

func (s *StandardEnvironment) GetProperty(key string) (value string, exists bool) {
	s.InitPropertyResolver()
	return s.propertyResolver.GetProperty(key)
}

func (s *StandardEnvironment) GetPropertyWithDef(key string, def string) string {
	s.InitPropertyResolver()
	return s.propertyResolver.GetPropertyWithDef(key, def)
}

func (s *StandardEnvironment) GetRequiredProperty(key string) string {
	s.InitPropertyResolver()
	return s.propertyResolver.GetRequiredProperty(key)
}

func (s *StandardEnvironment) ResolvePlaceholders(text string) string {
	s.InitPropertyResolver()
	return s.propertyResolver.ResolvePlaceholders(text)
}

func (s *StandardEnvironment) ResolveRequiredPlaceholders(text string) string {
	s.InitPropertyResolver()
	return s.propertyResolver.ResolveRequiredPlaceholders(text)
}

func (s *StandardEnvironment) GetActiveProfiles() []string {
	return s.activeProfiles
}

/**
添加的时候， 重新初始化所有的PropertySource，有些 profiles 路径可能是有占位符，需要先解决占位符
*/
func (s *StandardEnvironment) AddActiveProfiles(profiles ...string) {
	// TODO 先解决配置文件占位符问题

	// 如果没有占位符的情况下，还是找不到文件的，直接 panic

	// 文件存在，进行初始化 PropertySource

	// 初始化完成后，反复进行初始化，直到只剩下哪些还有占位符的

}

func (s *StandardEnvironment) SetActiveProfiles(profiles ...string) {
	// 先把原来的 PropertySource 全部删除
	if s.activeProfiles == nil || len(s.activeProfiles) < 1 {
		s.AddActiveProfiles(profiles...)
		return
	}
	if s.propertySources == nil || len(s.propertySources.propertySourceList) < 1 {
		s.activeProfiles = make([]string, 0)
		s.AddActiveProfiles(profiles...)
		return
	}

	// 删除原来的
	for _, profile := range s.activeProfiles {
		s.propertySources.Remove(profile)
	}
	// 重置
	s.activeProfiles = make([]string, 0)

	s.AddActiveProfiles(profiles...)
}

func (s *StandardEnvironment) GetPropertySources() *MutablePropertySources {
	if nil == s.propertySources {
		s.propertySources = &MutablePropertySources{
			propertySourceList: make([]PropertySource, 0),
		}
	}
	return s.propertySources
}

func (s *StandardEnvironment) Merge(parent Environment) {
	if parent == nil {
		return
	}

	parentSources := parent.GetPropertySources()
	if parentSources != nil {
		if s.propertySources == nil {
			s.GetPropertySources()
		}
		parentSources.Each(func(index int, source PropertySource) (stop bool) {
			if !s.propertySources.Contains(source.GetName()) {
				s.propertySources.AddLast(source)
			}
			return false
		})
	}
	// 添加激活的配置文件
	parentActiveProfiles := parent.GetActiveProfiles()
	if len(parentActiveProfiles) > 0 {
		s.AddActiveProfiles(parentActiveProfiles...)
	}
}
