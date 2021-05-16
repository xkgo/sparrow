package sparrow

import (
	"errors"
	"reflect"
)

type BeanDefinition struct {
	Bean             interface{}   // 源
	Type             reflect.Type  // 类型
	Value            reflect.Value // 值
	Primary          bool          // 是否是主Bean
	IsPropertiesBean bool          // 是否是 配置Bean
	ChangedListen    bool          // 作为配置Bean的时候，是否需要监听配置的变化，如果明显是一些不会变更的属性，那么就没必要进行监听了
	KeyPrefix        string        // 配置前缀
	Name             string        // bean 名称
	InitFn           reflect.Value // 初始化方法，无参函数
	DestroyFn        reflect.Value // 初始化方法，无参函数
	Ready            bool          // 是否已经准备好（已经初始化，并且成功执行了初始化方法）
	Order            int64         // 排序，默认是0，越小优先级越低
}

func newBeanDefinition(bean interface{}, beanName string, keyPrefix string, propertiesBean, changedListen, primary bool) *BeanDefinition {
	bd := &BeanDefinition{
		Bean:             bean,
		Name:             beanName,
		Primary:          primary,
		KeyPrefix:        keyPrefix,
		IsPropertiesBean: propertiesBean,
		ChangedListen:    changedListen,
	}

	bd.Type = reflect.TypeOf(bean)
	bd.Value = reflect.ValueOf(bean)

	if bd.Type.Kind() != reflect.Ptr && bd.Type.Kind() != reflect.Slice {
		panic(errors.New("bean必须是指针或者Slice对象"))
	}

	bd.InitFn = bd.Value.MethodByName("Init")
	bd.DestroyFn = bd.Value.MethodByName("Destroy")

	if len(bd.Name) < 1 {
		bd.Name = bd.Type.Elem().Name()
	}

	return bd
}
