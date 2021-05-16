package sparrow

import (
	"fmt"
	"github.com/xkgo/sparrow/env"
	"reflect"
	"strconv"
	"testing"
)

type IGet interface {
	Get(id int64) string
}

type UserApi interface {
	IGet
	Add(id int64) bool
}

type UserDao struct {
	tableName string `@Value:"key=${table.name.user:t_user},required=true"`
}

func (d *UserDao) insert() {
	fmt.Println("DAO: 插入数据。。。。。。: table: ", d.tableName)
}

func (d *UserDao) Add(id int64) bool {
	fmt.Println("添加用户Add(", id, ")")
	return true
}

func (d UserDao) Get(id int64) string {
	fmt.Println("查询用户 Get(", id, ")")
	return "User_" + strconv.FormatInt(id, 10)
}

func (d *UserDao) GetOrder() int64 {
	return 1
}

type MomentDao struct {
}

func (m *MomentDao) Get(id int64) string {
	return "MomentDao.Get: " + strconv.FormatInt(id, 10)
}

func (m *MomentDao) GetOrder() int64 {
	return 2
}

type UserService struct {
	userDaos []UserDao `@Inject:"required=true"`
	userDao  UserDao   `@Inject:""` //`@Inject:"name=userDao"`
	//userApi  UserApi   `@Inject:"required=true"`
	userApi2      *UserApi       `@Inject:"required=true"`
	iget          IGet           `@Inject:"required=true"`
	igets         []IGet         `@Inject:"required=true"`
	igetPtrs      []*IGet        `@Inject:"required=true"`
	appProperties *AppProperties `@Inject:"required=true"`
	//userApis []UserApi `@Inject:"required=true"`
}

func (s *UserService) add() {
	//fmt.Println("Service: Add 数据、、、、、、: ", len(s.userDaos))
	fmt.Println("Service: Add 数据、、、、、、: ")
	//if nil != s.userApi {
	//	s.userApi.Add(1)
	//}
	if nil != s.userApi2 {
		(*s.userApi2).Add(2)
	}
	if nil != s.iget {
		s.iget.Get(1000)
	}
	s.userDao.insert()

	fmt.Println("-------- userDaos 测试")
	if len(s.userDaos) > 0 {
		for idx, dao := range s.userDaos {
			dao.Get(int64(1000 + idx))
		}
	}
	fmt.Println("-------- igets 测试")
	if len(s.igets) > 0 {
		for idx, dao := range s.igets {
			dao.Get(int64(1000 + idx))
		}
	}
	fmt.Println("-------- igetPtrs 测试")
	if len(s.igetPtrs) > 0 {
		for idx, dao := range s.igetPtrs {
			(*dao).Get(int64(1000 + idx))
		}
	}
	fmt.Println("------------------- AppProperties")
	fmt.Println("properties.Name: ", s.appProperties.Name)
}

func (s *UserService) Init(environment env.Environment, app *Application) {
	fmt.Println("UserService 初始化： Init 方法, env:", environment, ", APP: ", app)
}

type AppProperties struct {
	Name string
}

func TestRegisterBean(t *testing.T) {

	RegisterBean(&UserDao{}, "userDao", true)
	RegisterBean(&UserDao{}, "userDao2", false)
	RegisterBean(&MomentDao{}, "momentDao", true)

	RegisterBean(&UserService{}, "userService", true)
	RegisterPropertiesBean(&AppProperties{}, "appProperties", "app.", true)

	_ = Run(env.New(env.ConfigDirs("./testdata")),
		WithName("test"),
		WithRunner(func(app *Application) (err error) {
			fmt.Println("App Run.....")

			userService := app.GetBeanByName("userService").(*UserService)
			userService.add()

			if bean, err := app.GetBeanByType(&UserService{}); err == nil {
				service := bean.(*UserService)
				service.add()
			}

			userServiceType := reflect.TypeOf(&UserService{})
			if bean, err := app.GetBeanByType(userServiceType); err == nil {
				service := bean.(*UserService)
				service.add()
			}

			beans := app.GetBeansOfType(reflect.TypeOf(&UserDao{}))
			fmt.Println("根据类型获取实例~~~")
			for beanName, bean := range beans {
				fmt.Println(beanName+" : ", bean)
			}

			properties := app.GetBeanByName("appProperties").(*AppProperties)
			fmt.Println("properties.Name: ", properties.Name)

			return nil
		}),
	)
}
