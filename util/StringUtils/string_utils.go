package StringUtils

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
)

var blankRegex, _ = regexp.Compile("^\\s+$")
var trimRegex = regexp.MustCompile("(^\\s+)|(\\s+$)")

// 邮箱地址正则
var emailRegex = regexp.MustCompile(`^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`)

/**
检查给定的字符串从 index 开始，是否一直匹配 substring
*/
func SubstringMatch(str string, index int, substring string) bool {
	if index+len(substring) > len(str) {
		return false
	}
	suffix := str[index:]
	idx := strings.Index(suffix, substring)
	return idx == 0
}

/**
从源串中查找给定串，并且从 fromIndex 开始查询
*/
func IndexFrom(source, findstr string, fromIndex int) (index int, err error) {
	if len(findstr) < 1 || len(source) < 1 {
		return -1, nil
	}
	if len(source) < fromIndex+len(findstr) {
		return -1, errors.New("find str out of range for source")
	}

	subfix := source[fromIndex:]
	index = strings.Index(subfix, findstr)
	if index > -1 {
		index = index + fromIndex
	}
	return index, nil
}

/**
从 fromIndex（含） 开始替换成 newVal
*/
func Replace(source, newVal string, fromIndex int) (val string, err error) {
	if len(newVal) < 1 {
		return source, nil
	}
	if len(source) < fromIndex+len(newVal) {
		return source, errors.New("newVal out of range for source")
	}

	prefix := source[0:fromIndex]
	suffix := source[fromIndex+len(newVal):]
	return prefix + newVal + suffix, nil
}

/**
从 fromIndex（含） 开始替换成 newVal
fromIndex include
endIndex exclude
*/
func ReplaceRange(source, newVal string, fromIndex, endIndex int) (val string, err error) {
	if len(source) <= fromIndex {
		return source, errors.New("fromIndex out of range for source")
	}
	if len(source) < endIndex {
		return source, errors.New("endIndex out of range for source")
	}

	prefix := source[0:fromIndex]
	suffix := source[endIndex:]
	return prefix + newVal + suffix, nil
}

/**
使用正则进行 split
*/
func SplitByRegex(source, regex string) []string {
	spaceRe, _ := regexp.Compile(regex)
	return spaceRe.Split(source, -1)
}

/**
使用正则进行 split
*/
func Split(source, sep string) []string {
	return strings.Split(source, sep)
}

/**
是否是空白字符串，长度为0，nil，空格、\t,\n\r 等字符构成
*/
func IsBlank(str string) bool {
	if len(str) < 1 {
		return true
	}
	return blankRegex.MatchString(str)
}

/**
是否是 非空白字符串，长度为非 0，nil，空格、\t,\n\r 等字符构成
*/
func IsNotBlank(str string) bool {
	return !IsBlank(str)
}

/*
 任意一个是 Blank 就返回true
 参数为空则返回false，表示没有任何为空
*/
func IsAnyBlank(strs ...string) bool {
	if len(strs) < 1 {
		return false
	}
	for _, str := range strs {
		if IsBlank(str) {
			return true
		}
	}
	return false
}

/*
 全部都是 Blank 就返回true
 参数为空则返回false，表示没有任何为空
*/
func IsAllBlank(strs ...string) bool {
	if len(strs) < 1 {
		return false
	}
	for _, str := range strs {
		if !IsBlank(str) {
			return false
		}
	}
	return true
}

/*
 有任意一个不为空则返回 true
 参数为空则返回false，表示没有任何为空
*/
func IsAnyNotBlank(strs ...string) bool {
	if len(strs) < 1 {
		return false
	}
	for _, str := range strs {
		if !IsBlank(str) {
			return true
		}
	}
	return false
}

/*
 所有都不为空则返回 true
 参数为空则返回false，表示没有任何为空
*/
func IsAllNotBlank(strs ...string) bool {
	if len(strs) < 1 {
		return false
	}
	for _, str := range strs {
		if IsBlank(str) {
			return false
		}
	}
	return true
}

/**
去掉空白字符串，空白包含： 长度0，nil，空格、\t,\n\r 等字符构成
*/
func Trim(str string) string {
	if len(str) < 1 {
		return ""
	}
	return trimRegex.ReplaceAllString(str, "")
}

// 首字母小写
func FirstLetterLower(s string) string {
	if len(s) < 1 {
		return s
	}
	runes := []rune(s)
	if unicode.IsUpper(runes[0]) {
		runes[0] = unicode.ToLower(runes[0])
		return string(runes)
	}
	return s
}

// 首字母大写
func FirstLetterUpper(s string) string {
	if len(s) < 1 {
		return s
	}
	runes := []rune(s)
	if unicode.IsLower(runes[0]) {
		runes[0] = unicode.ToUpper(runes[0])
		return string(runes)
	}
	return s
}

func EndWith(str, suffix string) bool {
	return strings.HasSuffix(str, suffix)
}

func StartWith(str, prefix string) bool {
	return strings.HasPrefix(str, prefix)
}

func Contains(str, substr string) bool {
	return strings.Contains(str, substr)
}

func IsEmail(str string) bool {
	if len(str) < 1 {
		return false
	}
	return emailRegex.MatchString(str)
}

// 返回 str 的after字符串之后的后缀，如果 after 为 nil 则返回 ""
func GetSuffix(str, after string) string {
	strLen, afterLen := len(str), len(after)
	if strLen < 1 || afterLen < 1 || strLen <= afterLen {
		return ""
	}

	index := strings.LastIndex(str, after)
	if index < 0 || strLen-index+1 < afterLen {
		return ""
	}

	return str[index+1:]
}

func EqualsIgnoreCase(str1, str2 string) bool {
	if str1 == str2 {
		return true
	}
	return strings.ToLower(str1) == strings.ToLower(str2)
}

func GetLowerLetter(str string, index int) (letter rune, exists bool) {
	if index < 0 {
		return 0, false
	}
	if index >= len(str) {
		return 0, false
	}
	ch := rune(str[index])
	if ch >= 'a' && ch <= 'z' {
		return ch, true
	}
	return 0, false
}

func GetUpperLetter(str string, index int) (letter rune, exists bool) {
	if index < 0 {
		return 0, false
	}
	if index >= len(str) {
		return 0, false
	}
	ch := rune(str[index])
	if ch >= 'A' && ch <= 'Z' {
		return ch, true
	}
	return 0, false
}

/**
判断首字母是否小写
*/
func IsFirstLetterLowerCase(str string) bool {
	if len(str) < 1 {
		return false
	}
	if len(str) == 1 {
		return strings.ToLower(str) == str
	}
	if _, ok := GetLowerLetter(str, 0); ok {
		return true
	}
	return false
}

/**
判断首字母是否是大写
*/
func IsFirstLetterUpperCase(str string) bool {
	if len(str) < 1 {
		return false
	}
	if len(str) == 1 {
		return strings.ToUpper(str) == str
	}
	if _, ok := GetUpperLetter(str, 0); ok {
		return true
	}
	return false
}
