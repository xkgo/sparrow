package util

import (
	"math"
	"reflect"
	"sort"
)

// 最低优先级，越大值表示优先级越低
const OrderedLowest int = math.MaxInt32

// 最高优先级，越小表示优先级越高
const OrderedHighest int = math.MinInt32

/**
排序接口
*/
type Ordered interface {

	/**
	返回排序值，越小优先级越高
	*/
	GetOrder() int
}

func SortByOrdered(slices interface{}) {
	value := reflect.ValueOf(slices)
	if value.Type().Kind() == reflect.Ptr {
		value = value.Elem()
	}
	sort.Slice(slices, func(i, j int) bool {

		v1 := value.Index(i)
		v2 := value.Index(j)

		order1 := 0
		order2 := 0

		if o1, ok := v1.Interface().(Ordered); ok {
			order1 = o1.GetOrder()
		} else if v1.CanAddr() {
			if o11, ok1 := v1.Addr().Interface().(Ordered); ok1 {
				order1 = o11.GetOrder()
			}
		}
		if o2, ok := value.Index(j).Interface().(Ordered); ok {
			order2 = o2.GetOrder()
		} else if v2.CanAddr() {
			if o21, ok2 := value.Index(j).Addr().Interface().(Ordered); ok2 {
				order2 = o21.GetOrder()
			}
		}
		return order1 < order2
	})
}
