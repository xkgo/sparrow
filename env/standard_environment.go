package env

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/xkgo/sparrow/deploy"
	"github.com/xkgo/sparrow/logger"
	"github.com/xkgo/sparrow/util/ConvertUtils"
	"github.com/xkgo/sparrow/util/JsonUtils"
	"github.com/xkgo/sparrow/util/StringUtils"
	"reflect"
	"regexp"
	"strconv"
)

const (
	DeployInfoEnvironmentPropertySourceName         = "deployInfoEnvironment"
	DefaultApplicationEnvironmentPropertySourceName = "defaultApplicationEnvironment"
	DeployInfoSetKey                                = "deploy.set"
	DeployInfoEnvKey                                = "deploy.env"
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
	> 添加系统环境变量
	> 自动解析当前运行环境相关属性： 环境(dev,test,fat,prod), set(分组：可能以全球大区、机房等来区分部署集群等等，将这个抽象即可)，将环境相关组成 propertySource ，然后添加进去 propertySources
	> 读取默认配置文件 application.properties|yml|toml, 然后添加到 propertySources 的 命令行之后，从 propertySources 中读取 sparrow-profile-include，作为 activeProfiles
	> 获取 profileDirs 下的所有配置文件，按照profile 分组，然后按照顺序依次加载配置文件，最后按顺序添加到 propertySources 的默认配置文件之后
	> profileDirs 都加载完成后， 将 additionalPropertySources 添加到 propertySources 之后
	> 添加系统环境变量到 propertySources 最后面
	*/
	propertySources *MutablePropertySources

	/**
	配置解析器，读取配置、处理占位符
	*/
	propertyResolver PropertyResolver

	/**
	配置key变更订阅列表
	*/
	propertyChangeListeners []*PropertyChangeListener
}

/**
新建环境
*/
func New(options ...Option) *StandardEnvironment {
	env := &StandardEnvironment{
		options:                 &Options{},
		propertyChangeListeners: make([]*PropertyChangeListener, 0),
	}

	// 设置选项
	if options != nil && len(options) > 0 {
		for _, option := range options {
			option(env)
		}
	}

	env.propertySources = NewMutablePropertySources()
	env.propertySources.Subscribe(func(self *MutablePropertySources) {
		prop := &logger.Properties{}
		_ = env.doBindProperties("logger.", prop, false)

		logger.Info("property source changed, will reset logger: ", JsonUtils.ToJsonStringWithoutError(prop))

		var log logger.Logger = logger.NewZapLogger(prop)

		if env.deployInfo != nil && env.deployInfo.Env == deploy.Dev {
			logger.SetConsoleLogger(logger.NewConsoleLogger(prop))
		} else {
			logger.SetRootLogger(log)
			logger.SetConsoleLogger(nil)
		}
	})

	// 将命令行参数作为最高优先级的属性来源
	env.propertySources.AddLast(NewCommandLinePropertySource(""))
	// 添加系统环境变量
	env.propertySources.AddLast(NewSystemEnvironmentPropertySource())

	// 添加部署信息到配置来源
	// 获取部署信息
	env.deployInfo = deploy.GetInfo()
	deployProperties := env.deployInfo.Properties
	if deployProperties == nil {
		deployProperties = make(map[string]string)
	}
	deployProperties[DeployInfoEnvKey] = string(env.deployInfo.Env)
	deployProperties[DeployInfoSetKey] = env.deployInfo.Set
	env.propertySources.AddLast(NewMapPropertySource(DeployInfoEnvironmentPropertySourceName, deployProperties))

	// 计算 profileDirs
	env.profileDirs = resolveProfileDirs(env.options.profileDirs)

	if env.profileDirs != nil && len(env.profileDirs) > 0 {
		// 读取默认配置文件 application.properties|yml|toml, 然后添加到 propertySources 的 命令行之后，从 propertySources 中读取 sparrow-profile-include，作为 activeProfiles
		defaultProfileInfo := getFirstDefaultApplicationProfileInfo(env.profileDirs)
		if nil != defaultProfileInfo {
			defaultPropertySource, err := ReadLocalFileAsPropertySource(DefaultApplicationEnvironmentPropertySourceName, defaultProfileInfo.path)
			if err != nil {
				errMsg := "读取默认配置文件异常:" + defaultProfileInfo.path + ", err:" + err.Error()
				logger.Error(errMsg, err)
				panic(err)
			}
			env.propertySources.AddFirst(defaultPropertySource)
		} else {
			// 添加一个空的默认配置来源
			env.propertySources.AddFirst(NewMapPropertySource(DefaultApplicationEnvironmentPropertySourceName, make(map[string]string)))
		}

		// 读取激活的 profiles
		include, _ := env.GetProperty(SparrowProfileIncludeKey)
		if len(include) < 1 {
			// 默认激活配置
			include = env.ResolvePlaceholders(fmt.Sprintf("${%s},${%s},${%s}-${%s}", DeployInfoSetKey, DeployInfoEnvKey, DeployInfoSetKey, DeployInfoEnvKey))
		}

		if len(include) > 0 {
			env.activeProfiles = StringUtils.SplitByRegex(include, "[,，;；\\s]+")
			actives := make(map[string]bool)
			for _, profile := range env.activeProfiles {
				actives[profile] = true
			}
			activeProfileInfos := getNotDefaultProfileInfoWithExtension(env.profileDirs, "")
			if len(activeProfileInfos) > 0 {
				for _, pis := range activeProfileInfos {
					onlyOne := len(pis) == 1
					for idx, pi := range pis {
						if !actives[pi.profile] {
							continue
						}
						name := pi.profile
						if !onlyOne {
							name = name + "_" + strconv.FormatInt(int64(idx), 10)
						}
						source, err := ReadLocalFileAsPropertySource(name, pi.path)
						if err != nil {
							errMsg := "读取默认ActiveProfile文件异常, " + pi.path + ", err:" + err.Error()
							logger.Error(errMsg, err)
							panic(err)
						}
						err = env.propertySources.AddBefore(DefaultApplicationEnvironmentPropertySourceName, source)
						if err != nil {
							errMsg := "添加profile[" + name + "]到默认应用profile[" + DefaultApplicationEnvironmentPropertySourceName + "]异常, " + pi.path + ", err:" + err.Error()
							logger.Error(errMsg, err)
							panic(err)
						}
					}
				}
			}
		}
	}

	// 将 additionalPropertySources 添加到 propertySources 之后
	additionalPropertySources := env.options.additionalPropertySources
	if nil != additionalPropertySources && len(additionalPropertySources.propertySourceList) > 0 {
		additionalPropertySources.Each(func(index int, source PropertySource) (stop bool) {
			if !env.propertySources.Contains(source.GetName()) {
				env.propertySources.AddLast(source)
			}
			return false
		})
	}

	// 刷新、初始化
	env.refresh()

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

func (s *StandardEnvironment) Subscribe(keyPattern string, handler func(event *KeyChangeEvent)) {
	if s.propertyChangeListeners == nil {
		s.propertyChangeListeners = make([]*PropertyChangeListener, 0)
	}
	s.propertyChangeListeners = append(s.propertyChangeListeners, NewPropertyChangeListener(keyPattern, handler))
}

func (s *StandardEnvironment) refresh() {
	s.initPropertySourceListen()
}

func (s *StandardEnvironment) initPropertySourceListen() {
	// 执行所有配置来源的监听
	s.propertySources.Each(func(index int, source PropertySource) (stop bool) {
		source.Subscribe("*", func() func(event *KeyChangeEvent) {
			return func(event *KeyChangeEvent) {
				logger.Info("收到配置来源["+source.GetName()+"]的配置变更事件：", event)
				s.onKeyChangeEvent(source, event)
			}
		}())
		return false
	})
}

/**
Key 变更处理
*/
func (s *StandardEnvironment) onKeyChangeEvent(source PropertySource, event *KeyChangeEvent) {
	// 执行监听器
	if len(s.propertyChangeListeners) > 0 {
		for _, listener := range s.propertyChangeListeners {
			keyPattern := listener.KeyPattern
			handler := listener.Handler
			if handler == nil {
				continue
			}
			if keyPattern == "" || keyPattern == "*" || keyPattern == event.Key {
				handler(event)
				continue
			}
			regex, err := regexp.Compile(keyPattern)
			if err == nil && regex.MatchString(event.Key) {
				handler(event)
			}
		}
	}
}

func (s *StandardEnvironment) BindProperties(keyPrefix string, cfgPtr interface{}) (err error) {
	return s.doBindProperties(keyPrefix, cfgPtr, true)
}

func (s *StandardEnvironment) doBindProperties(keyPrefix string, cfgPtr interface{}, listen bool) (err error) {
	// 反射解析所有属性
	t := reflect.TypeOf(cfgPtr)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	} else {
		return errors.New("注册配置Bean异常，必须是指针类型, 当前注册类型为：[" + t.Name() + "], keyPrefix:" + keyPrefix)
	}

	v := reflect.ValueOf(cfgPtr)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		tfield := t.Field(i)
		vfield := v.Field(i)

		fieldName := tfield.Name

		// 默认是首字母小写
		configKey := keyPrefix + StringUtils.FirstLetterLower(fieldName)
		subKey := tfield.Tag.Get("ck")
		if len(subKey) > 0 {
			configKey = keyPrefix + subKey
		}

		// 初始值
		initVal := tfield.Tag.Get("def")

		// 获取配置的值
		value, exists := s.GetProperty(configKey)
		if !exists {
			value = s.ResolvePlaceholders(initVal)
		}
		// 反射进行配置回写
		s.applyBeanPropertyValue(&t, &tfield, &vfield, initVal, value, PropertyUpdate)

		if listen {
			// 注册监听器, 占位符问题，每次变更的话，都需要重新检查占位符，当占位符变化这个也要变化
			s.Subscribe(configKey, func() func(event *KeyChangeEvent) {
				return func(event *KeyChangeEvent) {
					s.applyBeanPropertyValue(&t, &tfield, &vfield, initVal, event.Nv, event.ChangeType)
				}
			}())
		}
	}
	jsonText, err := json.Marshal(cfgPtr)
	if err != nil {
		return
	}
	logger.Info("绑定配置Bean["+t.Name()+"] => ", string(jsonText))
	return
}

func (s *StandardEnvironment) applyBeanPropertyValue(beanType *reflect.Type, tfield *reflect.StructField, vfield *reflect.Value, initVal string, value string, changeType KeyChangeType) {
	if PropertyDel == changeType {
		// 删除，设置回原来的初始值
		value = initVal
	}

	var cerr error

	defer func() {
		if r := recover(); r != nil {
			logger.Error("配置转换异常：panic,Property:[" + (*beanType).Name() + "." + tfield.Name + ":" + tfield.Type.Name() + "], newVal:[" + value + "]")
		} else {
			if cerr != nil {
				logger.Error("配置转换失败,Property:[" + (*beanType).Name() + "." + tfield.Name + ":" + tfield.Type.Name() + "], newVal:[" + value + "]")
			}
		}
	}()

	typeName := tfield.Type.Name()
	switch typeName {
	case "string":
		vfield.SetString(value)
	case "bool":
		if len(value) < 1 {
			value = "false"
		}
		if val, err := ConvertUtils.ToBool(value); err == nil {
			vfield.SetBool(val)
		} else {
			cerr = err
		}
	case "int", "int8", "int16", "int32", "int64":
		if len(value) < 1 {
			value = "0"
		}
		if val, err := ConvertUtils.ToInt64(value); err == nil {
			vfield.SetInt(val)
		} else {
			cerr = err
		}
	case "uint", "uint8", "uint16", "uint32", "uint64":
		if len(value) < 1 {
			value = "0"
		}
		if val, err := ConvertUtils.ToUint64(value); err == nil {
			vfield.SetUint(val)
		} else {
			cerr = err
		}
	case "float32", "float64":
		if len(value) < 1 {
			value = "0"
		}
		if val, err := ConvertUtils.ToFloat64(value); err == nil {
			vfield.SetFloat(val)
		} else {
			cerr = err
		}
	default:
		// 其他的，使用 JSON 转换
		rval := reflect.New(tfield.Type)
		aval := rval.Interface()

		if len(value) > 0 {
			cerr = json.Unmarshal([]byte(value), aval)
			if cerr == nil {
				vfield.Set(rval.Elem())
			}
		} else {
			vfield.Set(rval.Elem())
		}
	}
}
