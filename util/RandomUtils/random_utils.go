package RandomUtils

import (
	"math/rand"
)

const (
	upperCaseLettersIdx = iota
	lowerCaseLettersIdx
	numbersIdx
)

var charsMap = map[int]string{
	0: "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
	1: "abcdefghijklmnopqrstuvwxyz",
	2: "0123456789",
}

/**
生产随机字符串
*/
func RandomString(length int, includeUpperCase, includeLowerCase, includeNumber bool) string {
	bytes := make([]byte, length)

	for i := 0; i < length; i++ {
		subByte := random(includeUpperCase, includeLowerCase, includeNumber)
		bytes[i] = subByte
	}
	return string(bytes)
}

func RandomNumberString(length int) string {
	return RandomString(length, false, false, true)
}

func RandomUpperCaseLetterString(length int) string {
	return RandomString(length, true, false, false)
}

func RandomUpperCaseLetterAndNumberString(length int) string {
	return RandomString(length, true, false, true)
}

func RandomLowerCaseLetterString(length int) string {
	return RandomString(length, false, true, false)
}

func RandomLowerCaseLetterAndNumberString(length int) string {
	return RandomString(length, false, true, true)
}

/**
大小写
*/
func RandomLetterString(length int) string {
	return RandomString(length, true, true, false)
}

/**
大小写 + 数字
*/
func RandomLetterAndNumberString(length int) string {
	return RandomString(length, true, true, true)
}

/**
随机一个字符
*/
func random(includeUpperCase, includeLowerCase, includeNumber bool) (b byte) {
	if !includeLowerCase && !includeUpperCase && !includeNumber {
		includeUpperCase, includeLowerCase, includeNumber = true, true, true
	}
	types := make([]int, 0)
	if includeUpperCase {
		types = append(types, upperCaseLettersIdx)
	}
	if includeLowerCase {
		types = append(types, lowerCaseLettersIdx)
	}
	if includeNumber {
		types = append(types, numbersIdx)
	}
	var chars string
	if len(types) == 1 {
		chars = charsMap[types[0]]
	} else {
		chars = charsMap[types[rand.Intn(len(types))]]
	}
	return chars[rand.Intn(len(chars))]
}
