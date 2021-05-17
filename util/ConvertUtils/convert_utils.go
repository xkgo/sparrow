package ConvertUtils

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var trimRegex = regexp.MustCompile("(^\\s+)|(\\s+$)")

// 转成字符串 分隔符，默认是 |
func JoinAsString(vs ...interface{}) string {
	return JoinAsStringWithSep("|", vs...)
}

// 转成字符串
// sep 分隔符，默认是 |
func JoinAsStringWithSep(sep string, vs ...interface{}) string {
	if len(vs) < 1 {
		return ""
	}
	builder := strings.Builder{}
	vmidx := len(vs) - 1
	for idx, v := range vs {
		builder.WriteString(ToString(v))
		if vmidx > idx {
			builder.WriteString(sep)
		}
	}
	return builder.String()
}

func ToString(v interface{}) string {
	if nil == v {
		return ""
	}
	var content = ""
	switch v.(type) {
	case string:
		content = v.(string)
	case int:
		content = strconv.FormatInt(int64(v.(int)), 10)
	case int8:
		content = strconv.FormatInt(int64(v.(int8)), 10)
	case int16:
		content = strconv.FormatInt(int64(v.(int16)), 10)
	case int32:
		content = strconv.FormatInt(int64(v.(int32)), 10)
	case int64:
		content = strconv.FormatInt(v.(int64), 10)
	case uint:
		content = strconv.FormatUint(uint64(v.(uint)), 10)
	case uint8:
		content = strconv.FormatUint(uint64(v.(uint8)), 10)
	case uint16:
		content = strconv.FormatUint(uint64(v.(uint16)), 10)
	case uint32:
		content = strconv.FormatUint(uint64(v.(uint32)), 10)
	case uint64:
		content = strconv.FormatUint(v.(uint64), 10)
	case float32:
		content = strconv.FormatFloat(float64(v.(float32)), 'f', -1, 64)
	case float64:
		content = strconv.FormatFloat(v.(float64), 'f', -1, 64)
	case bool:
		content = strconv.FormatBool(v.(bool))
	default:
		content = fmt.Sprint(v)
	}
	return content
}

// 将前后缀的空白字符串都替换掉
func TrimBlank(str string) string {
	if len(str) < 1 {
		return str
	}
	return trimRegex.ReplaceAllString(str, "")
}

func ToInt(val string) (ret int, err error) {
	ret, err = strconv.Atoi(TrimBlank(val))
	return
}

func ToIntWithDef(val string, def int) int {
	ret, err := strconv.Atoi(TrimBlank(val))
	if err != nil {
		return def
	}
	return ret
}

func ToInt8(val string) (ret int8, err error) {
	v, err := strconv.ParseInt(TrimBlank(val), 10, 8)
	if err == nil {
		ret = int8(v)
	}
	return
}

func ToInt8WithDef(val string, def int8) int8 {
	ret, err := ToInt8(val)
	if err != nil {
		return def
	}
	return ret
}

func ToInt16(val string) (ret int16, err error) {
	v, err := strconv.ParseInt(TrimBlank(val), 10, 16)
	if err == nil {
		ret = int16(v)
	}
	return
}

func ToInt16WithDef(val string, def int16) int16 {
	ret, err := ToInt16(val)
	if err != nil {
		return def
	}
	return ret
}

func ToInt32(val string) (ret int32, err error) {
	v, err := strconv.ParseInt(TrimBlank(val), 10, 32)
	if err == nil {
		ret = int32(v)
	}
	return
}

func ToInt32WithDef(val string, def int32) int32 {
	ret, err := ToInt32(val)
	if err != nil {
		return def
	}
	return ret
}

func ToInt64(val string) (ret int64, err error) {
	ret, err = strconv.ParseInt(TrimBlank(val), 10, 64)
	return
}

func ToInt64WithDef(val string, def int64) int64 {
	ret, err := ToInt64(val)
	if err != nil {
		return def
	}
	return ret
}

func ToUint(val string) (ret uint, err error) {
	v, err := strconv.ParseUint(TrimBlank(val), 10, 0)
	if err == nil {
		ret = uint(v)
	}
	return
}

func ToUintWithDef(val string, def uint) uint {
	ret, err := ToUint(val)
	if err != nil {
		return def
	}
	return ret
}

func ToUint8(val string) (ret uint8, err error) {
	v, err := strconv.ParseUint(TrimBlank(val), 10, 8)
	if err == nil {
		ret = uint8(v)
	}
	return
}

func ToUint8WithDef(val string, def uint8) uint8 {
	ret, err := ToUint8(val)
	if err != nil {
		return def
	}
	return ret
}

func ToUint16(val string) (ret uint16, err error) {
	v, err := strconv.ParseUint(TrimBlank(val), 10, 16)
	if err == nil {
		ret = uint16(v)
	}
	return
}

func ToUint16WithDef(val string, def uint16) uint16 {
	ret, err := ToUint16(val)
	if err != nil {
		return def
	}
	return ret
}

func ToUint32(val string) (ret uint32, err error) {
	v, err := strconv.ParseUint(TrimBlank(val), 10, 32)
	if err == nil {
		ret = uint32(v)
	}
	return
}

func ToUint32WithDef(val string, def uint32) uint32 {
	ret, err := ToUint32(val)
	if err != nil {
		return def
	}
	return ret
}

func ToUint64(val string) (ret uint64, err error) {
	ret, err = strconv.ParseUint(TrimBlank(val), 10, 64)
	return
}

func ToUint64WithDef(val string, def uint64) uint64 {
	ret, err := ToUint64(val)
	if err != nil {
		return def
	}
	return ret
}

func ToFloat32(val string) (ret float32, err error) {
	v, err := strconv.ParseFloat(TrimBlank(val), 32)
	if err == nil {
		ret = float32(v)
	}
	return
}

func ToFloat32WithDef(val string, def float32) float32 {
	ret, err := ToFloat32(val)
	if err != nil {
		return def
	}
	return ret
}

func ToFloat64(val string) (ret float64, err error) {
	ret, err = strconv.ParseFloat(TrimBlank(val), 64)
	return
}

func ToFloat64WithDef(val string, def float64) float64 {
	ret, err := ToFloat64(val)
	if err != nil {
		return def
	}
	return ret
}

var boolTrueVals = []string{"yes", "true", "1", "ok", "on", "open"}
var boolFalseVals = []string{"no", "false", "0", "close"}

// yes|1|true|ok|on|open 都为 true
// no|0|false|close 都为 false
// 其余则为解析失败
func ToBool(val string) (ret bool, err error) {
	val = TrimBlank(val)
	slen := len(val)
	if slen < 1 || slen > 5 {
		err = errors.New("converter:ToBool:" + val + ":Invalid boolean text")
		return
	}
	val = strings.ToLower(val)
	for _, trueVal := range boolTrueVals {
		if trueVal == val {
			ret = true
			return
		}
	}
	for _, falseVal := range boolFalseVals {
		if falseVal == val {
			ret = false
			return
		}
	}

	err = errors.New("converter:ToBool:" + val + ":Invalid boolean text")
	return
}

func ToBoolWithDef(val string, def bool) bool {
	ret, err := ToBool(val)
	if err != nil {
		return def
	}
	return ret
}

func Int64MapToList(m map[int64]bool) []int64 {
	list := make([]int64, 0)
	if len(m) < 1 {
		return list
	}
	for k, _ := range m {
		list = append(list, k)
	}
	return list
}

func StringMapToList(m map[string]bool) []string {
	list := make([]string, 0)
	if len(m) < 1 {
		return list
	}
	for k, _ := range m {
		list = append(list, k)
	}
	return list
}

func Int64ArrToSetList(arr []int64) []int64 {
	if len(arr) < 1 {
		return arr
	}
	finalList := make([]int64, 0)
	exists := make(map[int64]bool)
	for _, item := range arr {
		if _, ok := exists[item]; !ok {
			finalList = append(finalList, item)
		}
	}
	return finalList
}

func ContainsInt64(val int64, arr []int64) bool {
	if len(arr) < 1 {
		return false
	}
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

func SplitAsInt64s(str, sep string, errToZero bool) ([]int64, error) {
	if len(str) < 1 {
		return []int64{}, nil
	}
	arr := strings.Split(str, sep)
	return StringArrayAsInt64s(arr, errToZero)
}

func StringArrayAsInt64s(arr []string, errToZero bool) ([]int64, error) {
	retList := make([]int64, 0)
	for _, item := range arr {
		v, err := strconv.ParseInt(item, 10, 64)
		if nil != err {
			if !errToZero {
				return retList, err
			}
		}
		retList = append(retList, v)
	}
	return retList, nil
}

func ToInt64ExistsMap(arr []int64) map[int64]bool {
	return ToInt64BoolMap(arr, true)
}

func ToInt64BoolMap(arr []int64, def bool) map[int64]bool {
	if len(arr) < 1 {
		return map[int64]bool{}
	}
	emap := make(map[int64]bool)
	for _, v := range arr {
		emap[v] = def
	}
	return emap
}

func ToStringBoolMap(arr []string, def bool) map[string]bool {
	if len(arr) < 1 {
		return map[string]bool{}
	}
	emap := make(map[string]bool)
	for _, v := range arr {
		emap[v] = def
	}
	return emap
}

func ToStringIntMap(arr []string, def int) map[string]int {
	if len(arr) < 1 {
		return map[string]int{}
	}
	emap := make(map[string]int)
	for _, v := range arr {
		emap[v] = def
	}
	return emap
}
