package env

import (
	"github.com/xkgo/sparrow/deploy"
	"github.com/xkgo/sparrow/logger"
)

/**
标准环境实现, 实现接口 Environment
*/
type StandardEnvironment struct {
	// 选项
	options *Options

	/**
	当前部署环境信息，部署环境&配置信息
	*/
	deployInfo *deploy.Info

	/*
		activeProfiles 配置搜索文件夹列表，来源规则如下：
		1. 如果是程序自己设置了（相当于是程序自己定义了路径），那么直接使用程序自己设定的
		2. 上一步没有获取到，解析命令行参数，当存在 --sparrow-profile-dirs=...... 的时候，那么直接以 --sparrow-profile-dirs 指定的为准，
		   如果--sparrow-profile-dirs设置了，但是为空字符串，那么默认: ./,./config,./conf
		3. 上一步没有获取到，那么计算系统环境变量中，是否定义了 sparrow-profile-dirs， 如果定义了就以此为准，如果定义的是空字符串，那么就是默认：./,./config,./conf
		4. 上一步没有，那么检查：./,./config,./conf，搜索是否存在（理论上是本框架支持的文件格式） application*.properties, application*.yml, application*.toml，存在则：./,./config,./conf
		5. 上一步没有，向上获取到一个目录，然后重复上一步，直到找到符合上一步的为止
		6. 如果始终找不到，那么说明不需要配置文件，系统一样是可以运行的
	*/
	profileDirs []string

	/**
	激活的 profile
	*/
	activeProfiles []string

	/**
	是否忽略无法处理的占位符，如果忽略则不处理，不忽略的话，那么遇到不能解析的占位符直接 panic
	*/
	ignoreUnresolvableNestedPlaceholders bool

	/**
	配置来源，默认是从 profileDirs 中进行检索实例化的，当然，也是可以定义外部配置中心的，比如 apollo、nacos、consul、zookeeper等等
	默认配置处理逻辑：
	> 将命令行参数 作为优先级最高的 propertySource， --------- 之后每次向 propertySources 添加元素，都要重新进行日志配置，这样子才能每次应用最新配置
	> 自动解析当前运行环境相关属性： 环境(dev,test,fat,prod), set(分组：可能以全球大区、机房等来区分部署集群等等，将这个抽象即可)，将环境相关组成 propertySource ，然后添加进去 propertySources
	> 关于运行环境属性解析，应该支持扩展，允许注入 不同部署平台检查的插件
	> 遍历每一个 profileDirs
	> profileDir 下，查找 application.properties|yml|toml, 以该配置文件后缀作为要使用的配置文件类型（只会加载和这个配置相同后缀的配置文件，因此一个APP使用一种后缀的配置文件即可）
	> 检查 sparrow.profile.include 配置项，然后进行有序追加
	> 遍历 activeProfiles, 检查 configDir 下是否有 activeProfile 的，有的话进行初始化，然后追加到 propertySources
	> profileDirs 都加载完成后， 将 additionalPropertySources 添加到 propertySources 之后
	> 添加系统环境变量到 propertySources 最后面
	*/
	propertySources *MutablePropertySources

	/**
	配置解析器，读取配置、处理占位符
	*/
	propertyResolver PropertyResolver

	/**
	配置key变更订阅列表,map 结构为：consumer->keyPattern->channel
	*/
	keyChangeSubscribers map[string]map[string]chan *KeyChangeEvent

	/**
	配置key变更队列，每一个PropertySource 都会有一个
	*/
	keyChangeQueues []chan *KeyChangeEvent
}

/**
新建环境
*/
func New(options ...Option) *StandardEnvironment {
	env := &StandardEnvironment{
		options: &Options{},
	}

	// 设置选项
	if options != nil && len(options) > 0 {
		for _, option := range options {
			option(env)
		}
	}

	// 获取部署信息
	env.deployInfo = deploy.GetInfo()

	// 计算 profileDirs
	env.profileDirs = resolveProfileDirs(env.options.profileDirs)

	// 计算 profiles

	return env
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
		if s.activeProfiles == nil {
			s.activeProfiles = parentActiveProfiles
		} else {
			s.activeProfiles = append(s.activeProfiles, parentActiveProfiles...)
		}
	}
}

func (s *StandardEnvironment) SubscribeKeyChange(consumer, keyPattern string) (queue chan *KeyChangeEvent) {
	if s.keyChangeSubscribers == nil {
		s.keyChangeSubscribers = make(map[string]map[string]chan *KeyChangeEvent)
	}
	keyPatternSubscribers, ok := s.keyChangeSubscribers[consumer]
	if ok {
		if subscriber, ex := keyPatternSubscribers[keyPattern]; ex && subscriber != nil {
			// 重复订阅的，那么会直接 panic
			errMsg := "Duplicate subscribe from consumer [" + consumer + "],keyPattern:[" + keyPattern + "]"
			logger.Fatal(errMsg)
			panic(errMsg)
		} else {
			queue = make(chan *KeyChangeEvent)
			keyPatternSubscribers[keyPattern] = queue
			return queue
		}
	} else {
		queue = make(chan *KeyChangeEvent)
		s.keyChangeSubscribers[consumer] = map[string]chan *KeyChangeEvent{
			keyPattern: queue,
		}
		return queue
	}
}

func (s *StandardEnvironment) UnsubscribeKeyChange(consumer, keyPattern string) {
	if s.keyChangeSubscribers == nil {
		return
	}
	keyPatternSubscribers, ok := s.keyChangeSubscribers[consumer]
	if ok {
		if subscriber, ex := keyPatternSubscribers[keyPattern]; ex && subscriber != nil {
			logger.Info("Unsubscribe from consumer ["+consumer+"],keyPattern:["+keyPattern+"], subscriber:", subscriber)
			close(subscriber)
		}
	}
}
