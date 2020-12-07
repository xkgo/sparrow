package env

import (
	"fmt"
	"github.com/magiconair/properties/assert"
	"github.com/xkgo/sparrow/logger"
	"testing"
	"time"
)

type FixedPublishEventMapPropertySource struct {
	MapPropertySource
}

func NewFixedPublishEventMapPropertySource(name string, properties map[string]string) *FixedPublishEventMapPropertySource {
	source := &FixedPublishEventMapPropertySource{
		MapPropertySource: MapPropertySource{
			name:       name,
			properties: properties,
		},
	}

	source.init()

	return source
}

func (m *FixedPublishEventMapPropertySource) Subscribe(keyPattern string, handler func(event *KeyChangeEvent)) {
	go func() {
		looptimes := 0
		for {
			if looptimes >= 3 {
				break
			}
			event := &KeyChangeEvent{
				Key:        "my.var",
				Ov:         "",
				Nv:         "aaaaa",
				ChangeType: PropertyUpdate,
			}
			fmt.Println("发布事件")
			handler(event)
			looptimes++
			time.Sleep(time.Duration(2) * time.Second)
		}
	}()
}

func (m *FixedPublishEventMapPropertySource) init() {
}

func TestStandardEnvironment_New(t *testing.T) {

	additionalPropertySources := NewMutablePropertySources(
		NewMapPropertySource("test", map[string]string{
			"my.var": "additional",
		}),
		NewFixedPublishEventMapPropertySource("Test1", map[string]string{
			"my.var1": "additional1",
		}),
		NewFixedPublishEventMapPropertySource("Test2", map[string]string{
			"my.var2": "additional2",
		}),
	)

	env := New(ConfigDirs("../testdata"), AdditionalPropertySources(additionalPropertySources))
	fmt.Println(env)

	fmt.Println(env.GetProperty("redis.server"))
	fmt.Println(env.GetProperty("test.name"))
	fmt.Println(env.GetProperty("my.var"))

	fmt.Println("xxxxxxxxx")

	for {
		time.Sleep(time.Duration(1) * time.Second)
	}
}

type UserInfo struct {
	Id       int    `ck:"id"`
	Username string `ck:"username"`
}

func TestStandardEnvironment_BindProperties(t *testing.T) {
	additionalPropertySources := NewMutablePropertySources(
		NewMapPropertySource("test", map[string]string{
			"user.id":       "1",
			"user.username": "Hello_${user.id}",
		}),
	)

	env := New(AdditionalPropertySources(additionalPropertySources))
	user := &UserInfo{}

	_ = env.BindProperties("user.", user)

	assert.Equal(t, 1, user.Id)
	assert.Equal(t, "Hello_1", user.Username)
	fmt.Println(env.GetProperty("GOPATH"))
}

func TestStandardEnvironment_Toml(t *testing.T) {
	env := New(
		ConfigDirs("../testdata/toml"),
		IgnoreUnresolvableNestedPlaceholders(true),
	)

	fmt.Println(env.GetProperty("test.name"))

}

func TestStandardEnvironment_BindLoggerProperties(t *testing.T) {
	additionalPropertySources := NewMutablePropertySources(
		NewMapPropertySource("test", map[string]string{
			"logger.level": "INFO",
		}),
	)

	env := New(AdditionalPropertySources(additionalPropertySources))
	props := &logger.Properties{}

	_ = env.BindProperties("logger.", props)

	fmt.Println(props)
}
