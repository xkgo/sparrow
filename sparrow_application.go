package sparrow

import (
	"errors"
	"github.com/xkgo/sparrow/annotations"
	"github.com/xkgo/sparrow/env"
	"github.com/xkgo/sparrow/logger"
	"github.com/xkgo/sparrow/util/GoUtils"
	"github.com/xkgo/sparrow/util/ReflectUtils"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

const (
	PropertyKeyApplicationName = "sparrow.application.name"
)

type application interface {
	/**
	  注册Bean, 允许多个别名
	  @param beanPtr 要注册的bean指针对象
	  @param beanName bean的名称
	  @param primary 是否是主bean
	*/
	RegisterBean(beanPtr interface{}, beanName string, primary bool)

	/**
	  注册 PropertiesBean, 允许多个别名
	  @param beanPtr 要注册的bean指针对象
	  @param beanName bean的名称
	  @param keyPrefix key前缀，会直接拼接，如果注意有必要的话需要加 .
	  @param primary 是否是主bean
	*/
	RegisterPropertiesBean(beanPtr interface{}, beanName string, keyPrefix string, primary bool)

	/**
	  获取指定名称的Bean，不存在则返回 nil
	*/
	GetBeanByName(beanName string) (beanPtr interface{})

	/**
	  获取指定类型的 beans， key 为beanName， value 为对应的对象指针
	*/
	GetBeansOfType(beanType reflect.Type) (beansPtr map[string]interface{})

	/**
	获取指定bean 模板的类型，如果有多个且没有 primary 的那么会返回error
	*/
	GetBeanByType(beanTemplate interface{}) (beanPtr interface{}, err error)

	/**
	  获取所有的 BeanNames
	*/
	GetBeanNames() []string

	AppendBeforeInitHandler(handler Runner)

	/**
	  运行程序
	*/
	Run(environment env.Environment, options ...Option) (err error)
}

type Application struct {
	Environment        env.Environment            // 环境
	Name               string                     // APP 名称
	runner             Runner                     // 环境、ioc初始化完成后执行的操作
	destroyer          Runner                     // 程序结束，执行销毁
	container          map[string]*BeanDefinition // 容器，key 为beanName
	beforeInitHandlers []Runner                   // 初始化之前执行
}

/**
创建 Application
*/
func NewApplication() *Application {
	return &Application{
		container:          make(map[string]*BeanDefinition),
		beforeInitHandlers: make([]Runner, 0),
	}
}

func (a *Application) doRegisterBean(beanPtr interface{}, beanName string, keyPrefix string, propertiesBean, changedListen, primary bool) {

	bd := newBeanDefinition(beanPtr, beanName, keyPrefix, propertiesBean, changedListen, primary)

	if bd, ok := a.container[bd.Name]; ok {
		err := errors.New("BeanName(" + bd.Name + ")已经注册过了：" + bd.Type.String() + ", 请勿重复注册")
		logger.Error(err)
		panic(err)
	}
	a.container[bd.Name] = bd
}

func (a *Application) RegisterBean(beanPtr interface{}, beanName string, primary bool) {
	a.doRegisterBean(beanPtr, beanName, "", false, false, primary)
}

func (a *Application) RegisterPropertiesBean(beanPtr interface{}, beanName string, keyPrefix string, primary bool) {
	a.doRegisterBean(beanPtr, beanName, keyPrefix, true, false, primary)
}

func (a *Application) RegisterPropertiesBeanListen(beanPtr interface{}, beanName string, keyPrefix string, changedListen, primary bool) {
	a.doRegisterBean(beanPtr, beanName, keyPrefix, true, changedListen, primary)
}

func (a *Application) GetBeanByName(beanName string) (beanPtr interface{}) {
	if bd, ok := a.container[beanName]; ok {
		return bd.Bean
	}
	return nil
}

func (a *Application) GetBeansOfType(beanType reflect.Type) (beansPtr map[string]interface{}) {
	beansPtr = make(map[string]interface{})
	if len(a.container) < 1 {
		return
	}

	for beanName, bd := range a.container {
		if bd.Type == beanType {
			beansPtr[beanName] = bd.Bean
		}
	}
	return
}

func (a *Application) GetBeanByType(beanTemplate interface{}) (beanPtr interface{}, err error) {
	beanType, ok := beanTemplate.(reflect.Type)
	if !ok {
		beanType = reflect.TypeOf(beanTemplate)
	} else {
		if val, ok := beanTemplate.(reflect.Value); ok {
			beanType = val.Type()
		}
	}
	if beanType.Kind() != reflect.Ptr {
		beanType = reflect.PtrTo(beanType)
	}
	bd, err := a.getPrimaryBeanDefinitionOfType(beanType)
	if err != nil {
		return nil, err
	}
	bean, err := ReflectUtils.ConvertTo(bd.Bean, beanType)
	if err == nil {
		return bean.Interface(), nil
	}
	return nil, err
}

/**
获取 primary bean definition
*/
func (a *Application) getPrimaryBeanDefinitionOfType(beanType reflect.Type) (bd *BeanDefinition, err error) {
	bds := a.getBeanDefinitionsOfType(beanType)
	if len(bds) == 1 {
		for _, item := range bds {
			return item, nil
		}
	}
	if len(bds) > 1 {
		for _, item := range bds {
			if item.Primary {
				return item, nil
			}
		}
	}
	return nil, errors.New("无法找到类型（" + beanType.String() + "）的 Primary BeanDefinition，期望 1, 实际:" + strconv.FormatInt(int64(len(bds)), 10))
}

func (a *Application) getBeanDefinitionsOfType(beanType reflect.Type) (bds map[string]*BeanDefinition) {
	bds = make(map[string]*BeanDefinition)
	if len(a.container) < 1 {
		return
	}

	nonePtrFieldType := beanType
	if nonePtrFieldType.Kind() == reflect.Ptr {
		nonePtrFieldType = nonePtrFieldType.Elem()
	}
	for beanName, bd := range a.container {
		if ReflectUtils.CanConvertTo(bd.Type, beanType) {
			bds[beanName] = bd
		}
	}
	return
}

func (a *Application) GetBeanNames() []string {
	beanNames := make([]string, 0)

	for beanName, _ := range a.container {
		beanNames = append(beanNames, beanName)
	}

	return beanNames
}

func (a *Application) AppendBeforeInitHandler(handler Runner) {
	if nil == a.beforeInitHandlers {
		a.beforeInitHandlers = make([]Runner, 0)
	}
	a.beforeInitHandlers = append(a.beforeInitHandlers, handler)
}

func (a *Application) Run(environment env.Environment, options ...Option) (err error) {

	defer func() {
		logger.Flush()
	}()

	a.Environment = environment
	if a.container == nil {
		a.container = make(map[string]*BeanDefinition)
	}
	if name, ok := environment.GetProperty(PropertyKeyApplicationName); ok && len(name) > 0 {
		app.Name = name
	} else {
		// 直接通过命令行参数计算
		app.Name = filepath.Base(os.Args[0])
	}

	if len(options) > 0 {
		for _, option := range options {
			option(app)
		}
	}

	if len(a.beforeInitHandlers) > 0 {
		for _, handler := range a.beforeInitHandlers {
			err := handler(a)
			if err != nil {
				return err
			}
		}
	}

	// 初始化，ioc容器初始化，处理依赖注入
	err = a.init()

	if err != nil {
		return
	}

	if nil != a.runner {
		err = a.runner(a)
		if err != nil {
			return
		}
	}

	// 自动执行注册过来的Bean的销毁方法

	// 程序结束
	if nil != a.destroyer {
		err = a.destroyer(a)
		if err != nil {
			return
		}
	}

	return
}

/**
应用程序初始化，
1. 处理 Bean以来关系注入（Bean注入@Autowire()，属性注入@Value()）
2. 执行函数初始化方法
*/
func (a *Application) init() (err error) {
	// 自动注入处理
	err = a.autoInjectProcess()
	if nil != err {
		return
	}

	return
}

func (a *Application) autoInjectProcess() (err error) {
	if len(a.container) < 1 {
		return
	}
	for _, bd := range a.container {
		err = a.wireBean(bd, nil)
		if nil != err {
			return
		}
	}
	return
}

/**
执行注入
*/
func (a *Application) wireBean(bd *BeanDefinition, dependencies []*BeanDefinition) (err error) {
	if bd.Ready {
		// 已经初始化过了
		return
	}

	// 属性绑定
	if bd.IsPropertiesBean {
		bd.Bean, err = a.Environment.BindPropertiesListen(bd.KeyPrefix, bd.Bean, bd.ChangedListen)
		if err != nil {
			return err
		}
		bd.Ready = true
		// 计算排序
		bd.Order, _ = ReflectUtils.GetRetInt64(bd.Bean, "GetOrder")
		return nil
	}

	if nil == dependencies {
		dependencies = make([]*BeanDefinition, 0)
	} else {
		// 检查是否有循环注入问题（检查当前元素在依赖链中是否存在，存在的话说明具有循环依赖）
		depChain := make([]string, 0)
		cycleDep := false
		for _, dbd := range dependencies {
			if dbd.Type == bd.Type {
				cycleDep = true
			}
			depChain = append(depChain, "|--("+dbd.Type.String()+")["+dbd.Name+"]\n|")
		}

		if cycleDep {
			depChain = append(depChain, "|--("+bd.Type.String()+")["+bd.Name+"]")
			return errors.New("存在循环依赖问题：\n\n" + strings.Join(depChain, "\n") + "\n")
		}
	}

	dependencies = append(dependencies, bd)

	// 执行自动注入
	bt := reflect.TypeOf(bd.Bean)
	if bt.Kind() == reflect.Ptr {
		bt = bt.Elem()
	}

	bv := reflect.ValueOf(bd.Bean)
	if bv.Kind() == reflect.Ptr {
		bv = bv.Elem()
	}

	// 遍历属性
	for i := 0; i < bt.NumField(); i++ {
		tf := bt.Field(i) // 属性类型

		inject, err := annotations.FindInject(tf.Tag)
		if err != nil {
			logger.Fatal("非法的@Inject 注解: ", err)
			return err
		}
		valueAnn, err := annotations.FindValue(tf.Tag)
		if err != nil {
			logger.Fatal("非法的@Value 注解: ", err)
			return err
		}

		if inject != nil && valueAnn != nil {
			logger.Fatal("不允许同时设置 @Inject 和 @Value 注解")
			return err
		}

		if inject != nil {
			err = a.wireBeanFieldByInjectAnnotation(bt.Field(i), bv.Field(i), inject, dependencies)
			if nil != err {
				return err
			}
		}
		if valueAnn != nil {
			err = a.wireBeanFieldByValueAnnotation(valueAnn, bt.Field(i), bv.Field(i))
			if nil != err {
				return err
			}
		}
	}

	err = a.doInitOrDestroy(bd, true)
	if nil != err {
		return
	}

	// 标记状态
	bd.Ready = true
	// 计算排序
	bd.Order, _ = ReflectUtils.GetRetInt64(bd.Bean, "GetOrder")
	return
}

func (a *Application) doInitOrDestroy(bd *BeanDefinition, init bool) (err error) {
	fn := bd.InitFn
	if !init {
		fn = bd.DestroyFn
	}
	// 执行初始化 或者 destroy 方法
	if fn.IsValid() {
		args := make([]reflect.Value, 0)
		for i := 0; i < fn.Type().NumIn(); i++ {
			arg := fn.Type().In(i)
			if arg.Kind() == reflect.Interface {
				if reflect.TypeOf(a.Environment).Implements(arg) {
					args = append(args, reflect.ValueOf(a.Environment))
				}
				if reflect.TypeOf(a).Implements(arg) {
					args = append(args, reflect.ValueOf(a))
				}
			} else {
				if arg == reflect.TypeOf(a) {
					args = append(args, reflect.ValueOf(a))
				} else {
					err = errors.New("初始化参数必须是 en.Environment 或者 *Application")
					return
				}
			}
		}
		fn.Call(args)
	}
	return
}

/**
按照 @Inject 注解进行注入
*/
func (a *Application) wireBeanFieldByInjectAnnotation(tf reflect.StructField, vf reflect.Value, inject *annotations.InjectAnn, dependencies []*BeanDefinition) (err error) {

	if tf.Type.Kind() == reflect.Slice {
		et := tf.Type.Elem() // 元素类型
		bds := a.getBeanDefinitionsOfType(et)
		if len(bds) < 1 {
			if inject.Required {
				err = errors.New("Ref Beans(ByBeanType) not found: " + inject.Name)
				logger.Fatal(err)
				return err
			}
			return // 没有强制要求
		}

		beanList := reflect.New(tf.Type).Elem()
		bdList := make([]*BeanDefinition, 0)
		for _, bd := range bds {
			if !bd.Ready { // 先处理以来问题
				err = a.wireBean(bd, dependencies)
				if err != nil {
					return err
				}
			}
			bdList = append(bdList, bd)
		}

		// 排序
		sort.Slice(bdList, func(i, j int) bool {
			o1 := bdList[i].Order
			o2 := bdList[j].Order
			return o2 < o1
		})

		for _, bd := range bdList {
			element, err := ReflectUtils.ConvertTo(bd.Bean, et)
			if err != nil || !element.IsValid() {
				return err
			}
			beanList = reflect.Append(beanList, element)
		}
		err = ReflectUtils.SetFieldValueByField(tf, vf, beanList)
		return err
	}

	// 非 Slice 类型
	var refBd *BeanDefinition
	if len(inject.Name) > 0 {
		refBd = a.container[inject.Name]
	} else {
		refBd, err = a.getPrimaryBeanDefinitionOfType(tf.Type)
		if err != nil {
			refBd = nil
		}
	}

	if nil == refBd && inject.Required {
		err = errors.New("===========Ref Bean(ByBeanName) not found: " + inject.Name)
		logger.Fatal(err.Error(), err)
		return err
	}

	if nil == refBd { // 没有强制要求需要注入
		return
	}

	if !refBd.Ready {
		err = a.wireBean(refBd, dependencies)
		if err != nil {
			return err
		}
	}

	return ReflectUtils.SetFieldValueByField(tf, vf, refBd.Bean)
}

func (a *Application) wireBeanFieldByValueAnnotation(valueAnn *annotations.ValueAnn, fieldType reflect.StructField, fieldValue reflect.Value) (err error) {
	GoUtils.Run(func() {
		var value string
		if valueAnn.Required {
			value = a.Environment.ResolveRequiredPlaceholders(valueAnn.Key)
		} else {
			value = a.Environment.ResolvePlaceholders(valueAnn.Key)
		}
		err = ReflectUtils.SetFieldValueByField(fieldType, fieldValue, value)
	}, func(r interface{}) {
		logger.Error("@Value 配置属性["+valueAnn.Key+"] 无法解析:", r)
	})
	return
}
