package StringUtils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSubstringMatch(t *testing.T) {
	assert.Equal(t, true, SubstringMatch("012345", 0, "012"))
	assert.Equal(t, true, SubstringMatch("012345", 1, "12"))
	assert.Equal(t, false, SubstringMatch("012345", 1, "012"))
	assert.Equal(t, true, SubstringMatch("中国人民", 1, "国"))
}

func TestReplace(t *testing.T) {
	actual, err := Replace("0123456", "aaa", 0)
	assert.Equal(t, "aaa3456", actual)

	actual, err = Replace("0123456", "aaa", 1)
	assert.Equal(t, "0aaa456", actual)

	actual, err = Replace("0123456", "aaa", 10)
	assert.NotNil(t, err)

	actual, err = Replace("0123456", "aaa", 4)
	assert.Equal(t, "0123aaa", actual)

	actual, err = Replace("0123456", "aaa", 5)
	assert.NotNil(t, err)

}

func TestIndexFrom(t *testing.T) {
	index, err := IndexFrom("0123456", "aa", 0)
	assert.Nil(t, err)
	assert.Equal(t, -1, index)

	index, err = IndexFrom("0123456", "23", 0)
	assert.Equal(t, 2, index)

	index, err = IndexFrom("0123456", "23", 2)
	assert.Equal(t, 2, index)

	index, err = IndexFrom("0123456", "23", 1)
	assert.Equal(t, 2, index)

	index, err = IndexFrom("0123456", "23", 3)
	assert.Equal(t, -1, index)

	index, err = IndexFrom("0123456", "23", 6)
	assert.NotNil(t, err)

}