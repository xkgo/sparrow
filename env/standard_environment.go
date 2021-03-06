package env

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/xkgo/sparrow/deploy"
	"github.com/xkgo/sparrow/logger"
	"github.com/xkgo/sparrow/util/FileUtils"
	"github.com/xkgo/sparrow/util/JsonUtils"
	"github.com/xkgo/sparrow/util/ReflectUtils"
	"github.com/xkgo/sparrow/util/StringUtils"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
)

const (
	DeployInfoEnvironmentPropertySourceName         = "deployInfoEnvironment"
	DefaultApplicationEnvironmentPropertySourceName = "defaultApplicationEnvironment"
	DeployInfoSetKey                                = "deploy.set"
	DeployInfoEnvKey                                = "deploy.env"
)

func init() {
	logger.InitLogger(&logger.Properties{})
}

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

	/**
	Beans
	*/
	bindBeans map[reflect.Type]interface{}
}

func (s *StandardEnvironment) IsDev() bool {
	return s.deployInfo.IsDev()
}

func (s *StandardEnvironment) IsTest() bool {
	return s.deployInfo.IsTest()
}

func (s *StandardEnvironment) IsFat() bool {
	return s.deployInfo.IsFat()
}

func (s *StandardEnvironment) IsProd() bool {
	return s.deployInfo.IsProd()
}

func (s *StandardEnvironment) GetEnv() deploy.Env {
	return s.deployInfo.Env
}

func (s *StandardEnvironment) GetSet() string {
	return s.deployInfo.Set
}

/**
新建环境
*/
func New(options ...Option) *StandardEnvironment {
	env := &StandardEnvironment{
		options:                 &Options{},
		propertyChangeListeners: make([]*PropertyChangeListener, 0),
		bindBeans:               make(map[reflect.Type]interface{}),
	}

	// 设置选项
	if options != nil && len(options) > 0 {
		for _, option := range options {
			option(env)
		}
	}

	env.propertySources = NewMutablePropertySources()
	var logProp *logger.Properties = nil
	env.propertySources.Subscribe(func(self *MutablePropertySources) {
		prop := &logger.Properties{}
		_, _ = env.doBindProperties("logger.", prop, false)
		if env.deployInfo != nil && env.deployInfo.Env == deploy.Dev {
			prop.ConsoleLog = true // 开发环境下强制开启 console log
		}
		if prop.Equals(logProp) {
			return
		}
		logProp = prop
		logger.Info("property source changed, will reset logger: ", JsonUtils.ToJsonStringWithoutError(prop))

		logger.InitLogger(prop)
	})

	if env.options.customDeployInfo != nil {
		env.deployInfo = env.options.customDeployInfo
	}

	// 将命令行参数作为最高优先级的属性来源
	env.propertySources.AddLast(NewCommandLinePropertySource(env.options.appendCommandLine))
	// 添加系统环境变量
	env.propertySources.AddLast(NewSystemEnvironmentPropertySource())

	if env.deployInfo == nil {
		// 添加部署信息到配置来源
		// 获取部署信息
		env.deployInfo = deploy.GetInfo()
		deployProperties := env.deployInfo.Properties
		if deployProperties == nil {
			deployProperties = make(map[string]string)
		}
	}
	if env.deployInfo.Properties == nil {
		env.deployInfo.Properties = make(map[string]string)
	}
	env.deployInfo.Properties[DeployInfoEnvKey] = string(env.deployInfo.Env)
	env.deployInfo.Properties[DeployInfoSetKey] = env.deployInfo.Set
	env.propertySources.AddLast(NewMapPropertySource(DeployInfoEnvironmentPropertySourceName, env.deployInfo.Properties))

	// 计算 profileDirs
	env.profileDirs = resolveProfileDirs(env.options.profileDirs)

	if len(env.profileDirs) < 1 && env.deployInfo.IsDev() {
		// 开发环境并且 profileDirs 为空，从新获取
		wd, err := os.Getwd()
		if err != nil {
			logger.Fatal("获取当前工作目录失败{os.Getwd()}：", err)
			panic(err)
		}
		logger.Info("开发环境，用户指定或者默认的找不到，从os.Getwd()查找")
		parentDir := wd
		for len(parentDir) > 0 {
			tempDirs := []string{parentDir, parentDir + string(filepath.Separator) + "testdata", parentDir + string(filepath.Separator) + "config", parentDir + string(filepath.Separator) + "conf"}
			env.profileDirs = resolveProfileDirs(tempDirs)
			logger.Info("开发环境,尝试查找配置文件目录：", tempDirs)
			if len(env.profileDirs) > 0 {
				break
			}
			parentDir = FileUtils.GetParentPath(parentDir)
		}
		logger.Info("开发环境，用户指定或者默认的找不到，从os.Getwd()查找到配置文件目录：", env.profileDirs)
	}

	activeProfiles := make([]string, 0)
	include, _ := env.GetProperty(SparrowProfileIncludeKey)
	if len(include) > 0 {
		profiles := StringUtils.SplitByRegex(include, "[,，;；\\s]+")
		if len(profiles) > 0 {
			activeProfiles = profiles
		}
	}

	if env.profileDirs != nil && len(env.profileDirs) > 0 {
		// 读取默认配置文件 application.properties|yml|toml, 然后添加到 propertySources 的 命令行之后，从 propertySources 中读取 sparrow-profile-include，作为 activeProfiles
		profileList := getDefaultApplicationProfileInfos(env.profileDirs)
		if len(profileList) > 0 {
			for idx, profile := range profileList {
				profileName := DefaultApplicationEnvironmentPropertySourceName
				if idx > 0 {
					profileName = DefaultApplicationEnvironmentPropertySourceName + ":" + profile.path
				}
				propertySource, err := ReadLocalFileAsPropertySource(profileName, profile.path)
				if err != nil {
					errMsg := "读取默认配置文件异常:" + profile.path + ", err:" + err.Error()
					logger.Error(errMsg, err)
					panic(err)
				}
				env.propertySources.AddFirst(propertySource)

				include, _ := env.GetProperty(SparrowProfileIncludeKey)
				if len(include) > 0 {
					profiles := StringUtils.SplitByRegex(include, "[,，;；\\s]+")
					if len(profiles) > 0 {
						activeProfiles = append(activeProfiles, profiles...)
					}
				}
			}
		} else {
			// 添加一个空的默认配置来源
			env.propertySources.AddFirst(NewMapPropertySource(DefaultApplicationEnvironmentPropertySourceName, make(map[string]string)))
		}

		if len(env.options.appendProfiles) > 0 {
			activeProfiles = append(activeProfiles, env.options.appendProfiles...)
		}
		if len(activeProfiles) < 1 {
			// 默认激活配置
			include = env.ResolvePlaceholders(fmt.Sprintf("${%s},${%s},${%s}-${%s}", DeployInfoSetKey, DeployInfoEnvKey, DeployInfoSetKey, DeployInfoEnvKey))
			profiles := StringUtils.SplitByRegex(include, "[,，;；\\s]+")
			if len(profiles) > 0 {
				activeProfiles = append(activeProfiles, profiles...)
			}
		}

		// 读取激活的 profiles
		existsProfiles := make(map[string]bool)
		env.activeProfiles = make([]string, 0)
		for _, profile := range activeProfiles {
			if _, ok := existsProfiles[profile]; !ok {
				env.activeProfiles = append(env.activeProfiles, profile)
				existsProfiles[profile] = true
			}
		}

		if len(env.activeProfiles) > 0 {
			activeProfileInfos := getNotDefaultProfileInfoWithExtension(env.profileDirs, "")
			includedProfiles := make(map[string]bool)
			for _, profile := range env.activeProfiles {
				if include, ok := includedProfiles[profile]; ok && include {
					continue
				}
				includedProfiles[profile] = true

				if pis, e := activeProfileInfos[profile]; e && pis != nil {
					for _, pi := range pis {
						name := pi.profile + ":" + pi.path
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

func (s *StandardEnvironment) BindProperties(keyPrefix string, cfgPtr interface{}) (beanPtr interface{}, err error) {
	return s.doBindProperties(keyPrefix, cfgPtr, false)
}

func (s *StandardEnvironment) BindPropertiesListen(keyPrefix string, cfgPtr interface{}, changedListen bool) (beanPtr interface{}, err error) {
	return s.doBindProperties(keyPrefix, cfgPtr, changedListen)
}

func (s *StandardEnvironment) doBindProperties(keyPrefix string, cfgPtr interface{}, listen bool) (beanPtr interface{}, err error) {
	// 反射解析所有属性
	t := reflect.TypeOf(cfgPtr)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	} else {
		return nil, errors.New("注册配置Bean异常，必须是指针类型, 当前注册类型为：[" + t.Name() + "], keyPrefix:" + keyPrefix)
	}

	if listen {
		// 已经绑定过了
		if bean, ok := s.bindBeans[t]; ok {
			return bean, nil
		} else {
			s.bindBeans[t] = cfgPtr
		}
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
		if len(subKey) < 1 {
			subKey = tfield.Tag.Get("sk")
		}
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
		s.applyBeanPropertyValue(t, tfield, vfield, initVal, value, PropertyUpdate)

		if listen {
			// 注册监听器, 占位符问题，每次变更的话，都需要重新检查占位符，当占位符变化这个也要变化
			s.Subscribe(configKey, func() func(event *KeyChangeEvent) {
				return func(event *KeyChangeEvent) {
					s.applyBeanPropertyValue(t, tfield, vfield, initVal, event.Nv, event.ChangeType)
				}
			}())
		}
	}
	jsonText, err := json.Marshal(cfgPtr)
	if err != nil {
		return nil, err
	}
	logger.Info("绑定配置Bean["+t.Name()+"] => ", string(jsonText))
	return cfgPtr, nil
}

/**
获取属性对象
@param typeTemplate 类型模板，可以提供属性类型的指针类型，也可以直接提供 reflect.Type 类型
*/
func (s *StandardEnvironment) GetProperties(typeTemplate interface{}) (beanPtr interface{}) {
	ptype, ok := typeTemplate.(reflect.Type)
	if !ok {
		ptype = reflect.TypeOf(typeTemplate)
	}
	if ptype.Kind() == reflect.Ptr {
		ptype = ptype.Elem()
	}

	if bean, ok := s.bindBeans[ptype]; ok {
		return bean
	}
	return nil
}

func (s *StandardEnvironment) applyBeanPropertyValue(beanType reflect.Type, tfield reflect.StructField, vfield reflect.Value, initVal string, value string, changeType KeyChangeType) {
	if PropertyDel == changeType {
		// 删除，设置回原来的初始值
		value = initVal
	}

	var cerr error

	defer func() {
		if r := recover(); r != nil {
			logger.Error("配置转换异常：panic,Property:["+beanType.Name()+"."+tfield.Name+":"+tfield.Type.Name()+"], newVal:["+value+"]", r)
		} else {
			if cerr != nil {
				logger.Error("配置转换失败,Property:["+beanType.Name()+"."+tfield.Name+":"+tfield.Type.Name()+"], newVal:["+value+"]", cerr)
			}
		}
	}()

	cerr = ReflectUtils.SetFieldValueByField(tfield, vfield, value)
}
