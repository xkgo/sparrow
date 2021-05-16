package annotations

import (
	"errors"
	"github.com/xkgo/sparrow/util/StringUtils"
	"reflect"
	"strings"
)

/*
@Inject 注入注解
*/
type InjectAnn struct {
	Name     string `json:"name"`
	Required bool   `json:"required"`
}

func ScanInner(tag reflect.StructTag, annName, split string, consumer func(key, val string)) (exists bool, err error) {
	if !strings.Contains(string(tag), annName) {
		return false, nil
	}
	injectCfg := tag.Get(annName)
	if len(injectCfg) < 1 {
		// 没有填写
		return true, nil
	}
	if len(split) < 1 {
		split = "="
	}

	injectCfg = strings.TrimSpace(injectCfg)
	kvs := StringUtils.SplitByRegex(injectCfg, "[,，；;]+")
	if len(kvs) < 1 {
		return false, errors.New(annName + "注解配置格式错误")
	}

	for _, kv := range kvs {
		arr := strings.Split(kv, split)
		if len(arr) == 2 {
			key := StringUtils.Trim(arr[0])
			value := StringUtils.Trim(arr[1])
			consumer(key, value)
		}
	}
	return true, nil
}

/*
@Inject 格式： name=value,required=true
*/
func FindInject(tag reflect.StructTag) (inject *InjectAnn, err error) {
	inject = &InjectAnn{}
	exists, err := ScanInner(tag, "@Inject", "=", func(key, val string) {
		switch key {
		case "name":
			inject.Name = val
		case "required":
			inject.Required = val == "true" || val == ""
		}
	})
	if err == nil && exists {
		return inject, nil
	}
	return nil, err
}

/**
@Value 注解
*/
type ValueAnn struct {
	Key      string // 配置属性名称
	Required bool   // 是否必须要求，如果是必须且不存在的话，那么直接抛出异常
}

/*
@Value 格式： ${...:${...}}
*/
func FindValue(tag reflect.StructTag) (value *ValueAnn, err error) {
	value = &ValueAnn{}
	exists, err := ScanInner(tag, "@Value", "=", func(key, val string) {
		switch key {
		case "key":
			value.Key = val
		case "required":
			value.Required = val == "true" || val == ""
		}
	})
	if err == nil && exists {
		return value, nil
	}
	return nil, err
}
