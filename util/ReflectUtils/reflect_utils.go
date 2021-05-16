package ReflectUtils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/xkgo/sparrow/util/ConvertUtils"
	"reflect"
	"strconv"
	"unsafe"
)

/**
https://www.jianshu.com/p/7b3638b47845
https://zhuanlan.zhihu.com/p/135224673
*/

/**
给对象设置属性值
*/
func SetFieldValueByName(object interface{}, fieldName string, value interface{}) (err error) {
	ot := reflect.TypeOf(object)
	if ot.Kind() != reflect.Ptr {
		return errors.New("待设置对象必须是指针类型~~~")
	}
	ot = ot.Elem()
	ov := reflect.ValueOf(object).Elem()

	field, exists := ot.FieldByName(fieldName)
	if !exists {
		return errors.New("对象(" + ot.String() + ")不存在属性：[" + fieldName + "]")
	}
	fieldValue := ov.FieldByName(fieldName)
	return SetFieldValueByField(field, fieldValue, value)
}

/**
设置属性值
*/
func doSetFieldValueByField(field reflect.StructField, fieldValue reflect.Value, value interface{}) (err error) {
	fieldType := field.Type

	if !fieldValue.CanSet() {
		fieldValue = reflect.NewAt(fieldValue.Type(), unsafe.Pointer(fieldValue.UnsafeAddr())).Elem()
	}

	nvalue, err := ConvertTo(value, fieldType)
	if err != nil {
		return err
	}
	fieldValue.Set(nvalue)
	return
}

/**
类型转换
*/
func ConvertTo(value interface{}, targetType reflect.Type) (ret reflect.Value, err error) {

	if val, ok := value.(reflect.Value); ok {
		return val, nil
	}

	var emptyValue = reflect.Value{}
	vtype := reflect.TypeOf(value)

	if vtype == targetType {
		// 类型完全相同
		return reflect.ValueOf(value), nil
	}
	// 属性是指针类型，但是给的值是非指针, 但是本质上是同一个类型
	if targetType.Kind() == reflect.Ptr && vtype.Kind() != reflect.Ptr && reflect.PtrTo(vtype) == targetType {
		nvalue := reflect.New(targetType.Elem()) // 可寻址
		pvalue := nvalue.Elem()                  // 可寻址
		pvalue.Set(reflect.ValueOf(value))
		return nvalue, nil
	}

	// 属性不是指针，但是给过来的是指针，但是本质上是一样类型
	if targetType.Kind() != reflect.Ptr && vtype.Kind() == reflect.Ptr && reflect.PtrTo(targetType) == vtype {
		// 类型完全相同
		return reflect.ValueOf(value).Elem(), nil
	}

	// 属性是接口类型的情况
	if targetType.Kind() == reflect.Interface {
		// 如果目标结果是指针类型，直接进行转换
		val := reflect.ValueOf(value).Convert(targetType)
		if val.IsValid() {
			return val, nil
		}
		return emptyValue, errors.New(fmt.Sprint("无法将指定对象(", value, ")转成目标接口类型(", targetType, ")"))
	}

	// 属性是接口的指针类型
	if targetType.Kind() == reflect.Ptr && targetType.Elem().Kind() == reflect.Interface {
		nvalue := reflect.New(targetType.Elem()) // 可寻址
		pvalue := nvalue.Elem()                  // 可寻址
		rval := reflect.ValueOf(value)
		val := rval.Convert(pvalue.Type())
		pvalue.Set(val)
		if val.IsValid() {
			return nvalue, nil
		}
		return emptyValue, errors.New(fmt.Sprint("无法将指定对象(", rval.String(), ")转成目标接口类型指针(", targetType.String(), ")"))
	}

	if strVal, ok := value.(string); ok {
		if reflectType, ok := types[targetType]; ok {
			rvalue, err := reflectType.Converter(strVal)
			if err != nil {
				return emptyValue, err
			}
			ret, err = ConvertTo(rvalue, targetType)
			if nil != err {
				return emptyValue, err
			}
			return ret, nil
		} else {
			// 其他的，使用 JSON 转换
			rval := reflect.New(targetType)
			aval := rval.Interface()

			if len(strVal) > 0 {
				err = json.Unmarshal([]byte(strVal), aval)
				if err == nil {
					return rval.Elem(), nil
				}
			} else {
				return rval.Elem(), nil
			}
		}
	}
	return emptyValue, errors.New(fmt.Sprint("无法将指定对象(", value, ")转成目标类型(", targetType, ")"))
}

func SetFieldValueByField(field reflect.StructField, fieldValue reflect.Value, value interface{}) (err error) {
	return doSetFieldValueByField(field, fieldValue, value)
}

/**
执行给定对象的方法，并返回方法返回值
*/
func InvokeObjectMethod(object interface{}, methodName string, args ...interface{}) (rvalues []reflect.Value, err error) {
	ov := reflect.ValueOf(object)
	if ov.Kind() != reflect.Ptr {
		return nil, errors.New("待设置对象必须是指针类型~~~")
	}
	method := ov.MethodByName(methodName)
	if !method.IsValid() {
		// 可能是私有方法
		ov = reflect.NewAt(ov.Type(), unsafe.Pointer(ov.Elem().UnsafeAddr())).Elem()
		method = ov.MethodByName(methodName)
		if !method.IsValid() {
			return nil, errors.New("对象(" + ov.String() + ")不存在方法：[" + methodName + "]")
		}
	}
	return InvokeMethod(method, args...)
}

/**
执行方法
*/
func InvokeMethod(method reflect.Value, args ...interface{}) (rvalues []reflect.Value, err error) {
	if !method.IsValid() {
		return nil, errors.New("method is invalid")
	}
	in := make([]reflect.Value, 0)
	if len(args) > 0 {
		for _, arg := range args {
			in = append(in, reflect.ValueOf(arg))
		}
	}
	rvalues = method.Call(in)
	return
}

/**
判断是否可以将 typeA 转换成 typeB 类型
*/
func CanConvertTo(typeA, typeB reflect.Type) bool {
	if typeA.Kind() == reflect.Ptr {
		typeA = typeA.Elem()
	}
	if typeB.Kind() == reflect.Ptr {
		typeB = typeB.Elem()
	}
	if typeA == typeB {
		return true
	}
	return IsImplements(typeA, typeB)
}

/**
判断 typeA 是否实现了 interfaceB 接口
*/
func IsImplements(typeA, interfaceB reflect.Type) bool {
	if interfaceB.Kind() != reflect.Interface {
		return false
	}

	typeAPtr := typeA
	noneTypeA := typeA
	if typeA.Kind() != reflect.Ptr {
		typeAPtr = reflect.PtrTo(typeA)
	} else {
		noneTypeA = typeA.Elem()
	}

	if interfaceB.NumMethod() < 1 || typeAPtr.NumMethod() < interfaceB.NumMethod() {
		return false
	}

	for i := 0; i < interfaceB.NumMethod(); i++ {
		iMethod := interfaceB.Method(i)
		sMethod, exists := typeAPtr.MethodByName(iMethod.Name)
		if !exists {
			return false
		}
		iMethodType := iMethod.Type
		sMethodType := sMethod.Type

		//fmt.Println(iMethodType.Kind(), ":", iMethod.Name, ":", iMethodType.NumIn(), ":", iMethodType.NumOut(), " ==> ", iMethodType.String())
		//fmt.Println(sMethodType.Kind(), ":", sMethod.Name, ":", sMethodType.NumIn(), ":", sMethodType.NumOut(), " ==> ", sMethodType.String())

		// 检查参数数量是否一致, 这里需要注意一下，如果 typeA 也是接口的话，那么这里就不一样了, 如果是结构体，那么应该从第二个参数开始匹配
		if noneTypeA.Kind() == reflect.Struct {
			if iMethodType.NumIn() != sMethodType.NumIn()-1 {
				return false
			}
		} else {
			if iMethodType.NumIn() != sMethodType.NumIn() {
				return false
			}
		}

		// 检查输出参数数量是否一致
		if iMethodType.NumOut() != sMethodType.NumOut() {
			return false
		}

		var inOffset int = 0
		if noneTypeA.Kind() == reflect.Struct {
			inOffset = 1
		}
		// 检查每个参数是否一致, 结构体和接口不一样
		for j := 0; j < iMethodType.NumIn(); j++ {
			iIn := iMethodType.In(j)
			sIn := sMethodType.In(j + inOffset)
			if iIn != sIn {
				return false
			}
		}

		// 检查每个输出参数是否一致
		for j := 0; j < iMethodType.NumOut(); j++ {
			iOut := iMethodType.Out(j)
			sOut := sMethodType.Out(j)
			if iOut != sOut {
				return false
			}
		}
	}

	return true
}

/**
获取类型
*/
func GetTypes(v ...interface{}) []reflect.Type {
	types := make([]reflect.Type, 0)
	if nil == v || len(v) < 1 {
		return types
	}
	for _, item := range v {
		t := reflect.TypeOf(item)
		types = append(types, t)
	}
	return types
}

/**
执行方法并返回int64 类型返回值
1. 要求返回值中有且仅有一个数字类型的，int,int8,int16,int32,int64,uint,uint8,uint16,uint32,uint64
*/
func GetRetInt64(object interface{}, methodName string, args ...interface{}) (ret int64, err error) {
	argTypes := GetTypes(args...)
	ov := reflect.ValueOf(object)

	if ov.Kind() != reflect.Ptr {
		return 0, errors.New("对象(" + ov.String() + ")必须是指针类型！")
	}

	method := ov.MethodByName(methodName)
	if !method.IsValid() {
		// 可能是私有方法
		ov = reflect.NewAt(ov.Type(), unsafe.Pointer(ov.Elem().UnsafeAddr())).Elem()
		method = ov.MethodByName(methodName)
		if !method.IsValid() {
			return 0, errors.New("对象(" + ov.String() + ")不存在方法：[" + methodName + "]")
		}
	}
	// 计算参数是否匹配
	methodType := method.Type()
	argOffset := methodType.NumIn() - len(argTypes)
	if argOffset != 0 && argOffset != 1 {
		return 0, errors.New("method args count not matched")
	}

	for i := 0; i < methodType.NumIn(); i++ {
		mArgType := methodType.In(i + argOffset)
		argType := argTypes[i+argOffset]
		if argType != nil && mArgType != argType {
			idxStr := strconv.FormatInt(int64(i), 10)
			return 0, errors.New("method args type not matched, arg[" + idxStr + "]:" + mArgType.String() + ", input[" + idxStr + "]:" + argType.String())
		}
	}

	// 输出参数检查
	errIndex := -1
	retIndex := -1
	for i := 0; i < methodType.NumOut(); i++ {
		out := methodType.Out(i)
		hit := false
		switch out {
		case Int, Int8, Int16, Int32, Int64, Uint, Uint8, Uint16, Uint32, Uint64, IntPtr, Int8Ptr, Int16Ptr, Int32Ptr, Int64Ptr, UintPtr, Uint8Ptr, Uint16Ptr, Uint32Ptr, Uint64Ptr:
			if retIndex >= 0 {
				return 0, errors.New("return params too much int type")
			} else {
				retIndex = i
			}
			hit = true
		case ErrorType, ErrorTypePtr:
			if errIndex >= 0 {
				return 0, errors.New("return params too much error type")
			} else {
				errIndex = i
			}
			hit = true
		}
		if errIndex < 0 && !hit {
			if CanConvertTo(out, ErrorType) {
				errIndex = i
			}
		}
	}

	if errIndex < 0 && retIndex < 0 {
		return 0, errors.New("could not found int64 return param")
	}
	argValues := make([]reflect.Value, 0)
	for _, a := range args {
		if av, ok := a.(reflect.Value); ok {
			argValues = append(argValues, av)
		} else {
			argValues = append(argValues, reflect.ValueOf(a))
		}
	}

	methodFn := reflect.ValueOf(object).MethodByName(methodName)
	values := methodFn.Call(argValues)
	if errIndex >= 0 {
		errVal := values[errIndex]
		if errVal.IsValid() && !errVal.IsNil() && !errVal.IsZero() {
			if errVal.Type().Kind() == reflect.Ptr {
				errVal = errVal.Elem()
			}
			return 0, errVal.Interface().(error)
		}
	}

	valVal := values[retIndex]
	if valVal.Type() == Int64 {
		return valVal.Interface().(int64), nil
	}
	if valVal.Type() == Int64Ptr {
		return *(valVal.Interface().(*int64)), nil
	}
	valStr := fmt.Sprint(valVal.Interface())
	return ConvertUtils.ToInt64(valStr)
}
