package util

import (
	"fmt"
	"testing"
)

type Obj struct {
	order int
}

func (o *Obj) GetOrder() int {
	return o.order
}

func (o *Obj) String() string {
	return fmt.Sprintf("%d", o.order)
}

func TestSortByOrdered(t *testing.T) {
	objs1 := make([]Obj, 0)
	objs1 = append(objs1, Obj{order: 1})
	objs1 = append(objs1, Obj{order: 3})
	objs1 = append(objs1, Obj{order: 2})

	fmt.Println(objs1)
	SortByOrdered(objs1)
	fmt.Println(objs1)

	fmt.Println("-------------------------------------------")

	objs := make([]*Obj, 0)
	objs = append(objs, &Obj{order: 1})
	objs = append(objs, &Obj{order: 3})
	objs = append(objs, &Obj{order: 2})

	fmt.Println(objs)
	SortByOrdered(objs)
	fmt.Println(objs)

	fmt.Println("-------------------------------------------")
}
