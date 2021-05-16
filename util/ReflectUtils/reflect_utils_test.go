package ReflectUtils

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestSetFieldValue(t *testing.T) {

	type Addr struct {
		country  string
		province string
	}

	type User struct {
		Id       int
		IdPtr    *int
		username string
		Addr     Addr
		AddrPtr  *Addr
	}

	user := &User{}

	// 类型完全一致
	assert.Nil(t, SetFieldValueByName(user, "Id", int(1)))
	assert.Equal(t, 1, user.Id)
	// 结构体一样
	assert.Nil(t, SetFieldValueByName(user, "Addr", Addr{country: "CN", province: "GZ"}))
	assert.NotNil(t, user.Addr)
	assert.Equal(t, "CN", user.Addr.country)
	assert.Equal(t, "GZ", user.Addr.province)

	// 属性不是指针，但是传入的是指针
	id := 2
	idPtr := &id
	assert.Nil(t, SetFieldValueByName(user, "Id", idPtr))
	assert.Equal(t, 2, user.Id)
	assert.Nil(t, SetFieldValueByName(user, "Addr", &Addr{country: "US", province: "WST"}))
	assert.NotNil(t, user.Addr)
	assert.Equal(t, "US", user.Addr.country)
	assert.Equal(t, "WST", user.Addr.province)

	// 属性是指针，传入的不是指针
	assert.Nil(t, SetFieldValueByName(user, "AddrPtr", Addr{country: "A", province: "B"}))
	assert.Equal(t, "A", user.AddrPtr.country)
	assert.Equal(t, "B", user.AddrPtr.province)
	id3 := 3
	assert.Nil(t, SetFieldValueByName(user, "IdPtr", id3))
	assert.Equal(t, 3, *user.IdPtr)
	id3 = 4
	id3Ptr := &id3
	assert.Nil(t, SetFieldValueByName(user, "IdPtr", id3Ptr))
	assert.Equal(t, 4, *user.IdPtr)

	// 字符串的注入
	assert.Nil(t, SetFieldValueByName(user, "username", "Arvin"))
	assert.Equal(t, "Arvin", user.username)
	assert.Nil(t, SetFieldValueByName(user, "Id", "10"))
	assert.Equal(t, 10, user.Id)
	assert.Nil(t, SetFieldValueByName(user, "IdPtr", "10"))
	assert.Equal(t, 10, *user.IdPtr)

}

func TestRegisterType(t *testing.T) {

	assert.NotNil(t, types[StringPtr])
	assert.NotNil(t, types[String])
	assert.NotNil(t, types[Int])
	assert.NotNil(t, types[IntPtr])
}

func TestSliceValue(t *testing.T) {
	type User struct {
		Ids    []int
		IdsPtr []*int
	}

	user := &User{}
	assert.Nil(t, SetFieldValueByName(user, "Ids", "[1,2]"))
	assert.ElementsMatch(t, []int{1, 2}, user.Ids)

	assert.Nil(t, SetFieldValueByName(user, "IdsPtr", "[3,4]"))
	assert.Equal(t, 3, *user.IdsPtr[0])
	assert.Equal(t, 4, *user.IdsPtr[1])
}

type Sex int8

const (
	SexUnknown Sex = 0
	Male       Sex = 1
	Female     Sex = 2
)

// 自定义的类型，需要自己注册转换器了
func TestEnumSetField(t *testing.T) {
	type User struct {
		Sex Sex
	}

	user := &User{}

	RegisterType(SexUnknown, func(value string) (val interface{}, err error) {
		switch value {
		case "0", "Unknown", "SexUnknown":
			return SexUnknown, nil
		case "1", "Male":
			return Male, nil
		case "2", "Female":
			return Female, nil
		}
		return nil, errors.New("invalid sex string: " + value)
	})

	assert.Nil(t, SetFieldValueByName(user, "Sex", "1"))
	assert.Equal(t, Male, user.Sex)

	assert.Nil(t, SetFieldValueByName(user, "Sex", "Female"))
	assert.Equal(t, Female, user.Sex)
}

func sayHello(name string) string {
	return "Hello, " + name
}

func TestMethodInvoke(t *testing.T) {
	var sayHi = func(name string) string {
		return "Hi, " + name
	}

	method := reflect.ValueOf(sayHi)

	values, _ := InvokeMethod(method, "arvin")
	assert.Equal(t, "Hi, arvin", values[0].Interface())

	helloMethod := reflect.ValueOf(sayHello)
	values, _ = InvokeMethod(helloMethod, "arvin")
	assert.Equal(t, "Hello, arvin", values[0].Interface())
}

type HelloService struct {
}

func (s *HelloService) sayHi(name string) string {
	return "Hi, " + name
}

func (s *HelloService) SayHi(name string) string {
	return "Hi, " + name
}

func TestObjectMethodInvoke(t *testing.T) {

	service := &HelloService{}

	values, err := InvokeObjectMethod(service, "SayHi", "arvin")
	assert.Nil(t, err)
	assert.Equal(t, "Hi, arvin", values[0].Interface())

	// 不支持私有函数
	values, err = InvokeObjectMethod(service, "sayHi", "arvin")
	assert.NotNil(t, err)
}

func TestGetTypes(t *testing.T) {
	types := GetTypes(1, "", nil)
	fmt.Println(types)
}

type GetRetInt64Service struct {
}

func (s *GetRetInt64Service) Get1() int {
	return 1
}

func (s *GetRetInt64Service) Get2(id int64) int64 {
	return id
}

func (s *GetRetInt64Service) Get3(id int64, name string) (int64, error, string) {
	if id < 1 {
		return 0, errors.New("ID 应该大于等于 1"), name
	}
	return id, nil, name
}

func TestGetRetInt64(t *testing.T) {
	service := &GetRetInt64Service{}

	ret, err := GetRetInt64(service, "Get1")
	assert.Nil(t, err)
	assert.Equal(t, int64(1), ret)

	ret, err = GetRetInt64(service, "Get2", int64(2))
	assert.Nil(t, err)
	assert.Equal(t, int64(2), ret)

	ret, err = GetRetInt64(service, "Get3", int64(3), "Arvin")
	assert.Nil(t, err)
	assert.Equal(t, int64(3), ret)

	ret, err = GetRetInt64(service, "Get3", int64(-1), "Arvin")
	assert.NotNil(t, err)

}
