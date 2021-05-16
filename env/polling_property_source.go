package env

import (
	"errors"
	"github.com/xkgo/sparrow/logger"
	"github.com/xkgo/sparrow/util/GoUtils"
	"regexp"
	"sync"
	"time"
)

type PropertyReader interface {
	/**
	读取所有的配置项, key-value 结构
	*/
	ReadAll() (kvs map[string]string, err error)
}

type PropertyReaderWrapper struct {
	/**
	读取配置的方法
	*/
	Reader func() (kvs map[string]string, err error)
}

func NewPropertyReader(reader func() (kvs map[string]string, err error)) *PropertyReaderWrapper {
	return &PropertyReaderWrapper{Reader: reader}
}

func (p *PropertyReaderWrapper) ReadAll() (kvs map[string]string, err error) {
	return p.Reader()
}

/**
基于轮询实现的配置来源
*/
type PollingPropertySource struct {
	Name            string            // 名称
	PropertyReader  PropertyReader    // 配置读取实现
	PollingInterval int64             // 轮询间隔，单位：秒
	kvs             map[string]string // 内存配置项， key->value
	scheduleOnce    sync.Once
	/**
	配置key变更订阅列表
	*/
	propertyChangeListeners []*PropertyChangeListener
}

/*
创建轮询配置源, 参数不正确的话直接抛出panic
@param refreshInterval 刷新时间间隔，单位：秒， 小于0表示不进行轮询
*/
func NewPollingPropertySource(name string, refreshInterval int64, reader PropertyReader) (source *PollingPropertySource, err error) {
	if reader == nil {
		err = errors.New("ConfigReader is required")
		panic(err)
		return
	}

	source = &PollingPropertySource{
		Name:            name,
		PollingInterval: refreshInterval,
		PropertyReader:  reader,
	}

	source.Init()

	return
}

func (p *PollingPropertySource) Init() {
	p.scheduleReload()
}

func (p *PollingPropertySource) scheduleReload() {

	_ = p.Reload()

	p.scheduleOnce.Do(func() {
		if p.PollingInterval < 1 {
			return
		}

		// 调度刷新
		GoUtils.RunGoroutine(func() {
			logger.Info("调度刷新配置，刷新间隔：[", p.PollingInterval, "]秒")
			for {
				_ = p.Reload()
				time.Sleep(time.Duration(p.PollingInterval) * time.Second)
			}
		}, func(r interface{}) {
			logger.Info("调度刷新配置异常：", r)
		})
	})
}

/*
重新加载配置
*/
func (p *PollingPropertySource) Reload() (err error) {

	defer func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("本次配置Reload失败， error:", r)
			}
		}()
	}()

	nkvs, err := p.PropertyReader.ReadAll()
	if err != nil {
		return
	}
	if nkvs == nil {
		nkvs = make(map[string]string)
	}
	if p.kvs == nil {
		p.kvs = make(map[string]string)
	}

	okvs := p.kvs

	// 新的配置
	p.kvs = nkvs

	// 比较计算哪些属性发生变更，变化了的调用变更监听器
	if len(p.propertyChangeListeners) < 1 || okvs == nil {
		// 首次加载
		return
	}

	// 判断是否有更新或者删除
	for key, ov := range okvs {
		nv, exists := nkvs[key]
		if exists && nv != ov {
			// 更新了
			p.onKeyChangeEvent(&KeyChangeEvent{
				Key:        key,
				Ov:         ov,
				Nv:         nv,
				ChangeType: PropertyUpdate,
			})
		} else if !exists {
			// 删除
			p.onKeyChangeEvent(&KeyChangeEvent{
				Key:        key,
				Ov:         ov,
				Nv:         "",
				ChangeType: PropertyDel,
			})
		}
	}

	for key, nv := range nkvs {
		if _, exists := okvs[key]; !exists {
			// 添加
			p.onKeyChangeEvent(&KeyChangeEvent{
				Key:        key,
				Ov:         "",
				Nv:         nv,
				ChangeType: PropertyAdd,
			})
		}
	}
	return nil
}

/**
Key 变更处理
*/
func (p *PollingPropertySource) onKeyChangeEvent(event *KeyChangeEvent) {
	logger.Info("["+p.Name+"]配置发生了变更：key:["+event.Key+"], ov:["+event.Ov+"], nv:["+event.Nv+"], changeType:[", event.ChangeType+"]")
	// 执行监听器
	if len(p.propertyChangeListeners) > 0 {
		for _, listener := range p.propertyChangeListeners {
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
				GoUtils.Run(func() {
					handler(event)
				}, func(r interface{}) {
					logger.Warn("配置源["+p.Name+"]执行配置变更[", listener, "]发生panic： ", r)
				})
			}
		}
	}
}

func (p *PollingPropertySource) GetName() string {
	return p.Name
}

func (p *PollingPropertySource) GetProperty(key string) (value string, exists bool) {
	value, exists = p.kvs[key]
	return
}

func (p *PollingPropertySource) GetPropertyWithDef(key string, def string) string {
	if value, exists := p.kvs[key]; exists && len(value) > 0 {
		return value
	}
	return def
}

func (p *PollingPropertySource) Each(consumer func(key string, value string) (stop bool)) {
	for k, v := range p.kvs {
		if consumer(k, v) {
			return
		}
	}
}

func (p *PollingPropertySource) Subscribe(keyPattern string, handler func(event *KeyChangeEvent)) {
	if p.propertyChangeListeners == nil {
		p.propertyChangeListeners = make([]*PropertyChangeListener, 0)
	}
	p.propertyChangeListeners = append(p.propertyChangeListeners, NewPropertyChangeListener(keyPattern, handler))
}
