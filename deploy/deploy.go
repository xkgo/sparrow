package deploy

import (
	"github.com/xkgo/sparrow/logger"
	"github.com/xkgo/sparrow/util/JsonUtils"
	"github.com/xkgo/sparrow/util/StringUtils"
	"os"
	"strings"
)

type Env string

const (
	Dev  Env = "dev"  // 开发
	Test Env = "test" // 测试
	Fat  Env = "fat"  // 预发布
	Prod Env = "prod" // 生产
)

type Info struct {
	Env        Env               // 当前运行环境
	Set        string            // 当前部署所在部署集，如所在大区，或者说所在部署集群等标识， 默认就是空字符串
	Properties map[string]string // 当前运行环境下的属性配置信息， 可能每个部署平台都有自己特殊的一些配置信息
}

func (i *Info) String() string {
	return JsonUtils.ToJsonStringWithoutError(i)
}

func (i *Info) IsDev() bool {
	return Dev == i.Env
}

func (i *Info) IsTest() bool {
	return Test == i.Env
}

func (i *Info) IsFat() bool {
	return Fat == i.Env
}

func (i *Info) IsProd() bool {
	return Prod == i.Env
}

// 自定义环境
var customEnvs = make(map[string]Env)

/**
自定义环境
*/
func DefineEnv(envStr string, env Env) {
	if len(envStr) > 0 && len(env) > 0 {
		customEnvs[envStr] = env
	}
}

func ParseEnv(env string) Env {
	if len(customEnvs) > 0 {
		for key, val := range customEnvs {
			if StringUtils.EqualsIgnoreCase(env, key) {
				return val
			}
		}
	}
	if StringUtils.EqualsIgnoreCase(env, "test") {
		return Test
	}
	if StringUtils.EqualsIgnoreCase(env, "fat") {
		return Fat
	}
	if StringUtils.EqualsIgnoreCase(env, "prod") {
		return Prod
	}
	return Dev
}

/**
自定义部署信息，如果设置了这个，那么直接以这个为准
*/
var customDeployInfo *Info

/**
如果设置了这个，那么直接返回这个值，一般用来测试
*/
func SetCustom(info *Info) {
	customDeployInfo = info
}

/**
部署信息检测
*/
type Detect interface {
	/**
	部署平台名称
	*/
	GetName() string
	/**
	检测部署环境信息
	*/
	Detect() *Info
}

type DetectWrapper struct {
	Name    string
	Handler func() *Info
}

func (d *DetectWrapper) GetName() string {
	return d.Name
}

func (d *DetectWrapper) Detect() *Info {
	return d.Handler()
}

// 默认部署环境识别实现
var defaultDetect Detect

// 用户自定义注册的 Detect 列表
var registeredDetectList = make([]Detect, 0)

/**
清空当前注册了个DetectList
*/
func ClearDetectList() {
	registeredDetectList = make([]Detect, 0)
}

func SetDetectList(detects ...Detect) {
	if nil != detects && len(detects) > 0 {
		registeredDetectList = append(make([]Detect, 0), detects...)
	}
}

func GetDetectList() []Detect {
	return registeredDetectList
}

/**
添加到最高优先级的列表
*/
func AddFirst(name string, handler func() *Info) {
	if len(name) < 1 || nil == handler {
		return
	}
	var detect Detect = &DetectWrapper{
		Name:    name,
		Handler: handler,
	}
	list := make([]Detect, 0)
	list = append(list, detect)
	if nil != registeredDetectList && len(registeredDetectList) > 0 {
		list = append(list, registeredDetectList...)
	}
	registeredDetectList = list
}

func AddLast(name string, handler func() *Info) {
	if len(name) < 1 || nil == handler {
		return
	}
	var detect Detect = &DetectWrapper{
		Name:    name,
		Handler: handler,
	}
	if registeredDetectList == nil {
		registeredDetectList = make([]Detect, 0)
	}
	registeredDetectList = append(registeredDetectList, detect)
}

/**
获取当前部署环境信息
*/
func GetInfo() *Info {
	if nil != customDeployInfo {
		return customDeployInfo
	}
	if registeredDetectList != nil && len(registeredDetectList) > 0 {
		for _, detect := range registeredDetectList {
			info := detect.Detect()
			logger.Infof("正在使用部署平台检测器：%v, 进行j检测", detect.GetName())
			if nil != info {
				logger.Info("当前部署检测器检测成功，Detect:"+detect.GetName()+", info:", info)
				if info.Properties == nil {
					info.Properties = make(map[string]string)
				}
				return info
			}
		}
	}
	info := defaultDetect.Detect()
	logger.Info("当前部署检测器检测成功，Detect:"+defaultDetect.GetName()+", info:", info)
	if info.Properties == nil {
		info.Properties = make(map[string]string)
	}
	return info
}

/**
获取命令行参数
*/
func GetCommandLineProperties(appendCommandLine string) map[string]string {
	var args []string = os.Args
	if len(appendCommandLine) > 4 {
		args = append(args, StringUtils.SplitByRegex(appendCommandLine, "\\s+")...)
	}

	properties := make(map[string]string)

	for _, arg := range args {
		if len(arg) < 4 || !strings.HasPrefix(arg, "--") { // 至少四个字符， --k=
			continue
		}
		// 替换一次
		arg = StringUtils.Trim(strings.Replace(arg, "--", "", 1))
		index := strings.Index(arg, "=")
		if index < 1 {
			continue
		}

		key := StringUtils.Trim(arg[0:index])
		value := StringUtils.Trim(arg[index+1:])

		properties[key] = value
	}
	return properties
}
