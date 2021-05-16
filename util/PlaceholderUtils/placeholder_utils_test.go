package PlaceholderUtils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResolvePlaceholdersExt(t *testing.T) {
	properties := map[string]string{
		"circular.var1": "${circular.var2}",
		"circular.var2": "${circular.var1}",
		"user.name":     "Arvin",
	}

	// 循环引用
	assert.Panics(t, func() { ResolvePlaceholders("${circular.var1}", properties) })

	assert.Equal(t, "你好:Arvin", ResolvePlaceholders("你好:${user.name}", properties))
	assert.Equal(t, "你好:Go", ResolvePlaceholders("你好:${user.no:Go}", properties))
	assert.Equal(t, "你好:Arvin", ResolvePlaceholders("你好:${user.no:${user.name}}", properties))
	assert.Equal(t, "你好:", ResolvePlaceholders("你好:${user.no:}", properties))
	assert.Equal(t, "你好:--Arvin", ResolvePlaceholders("你好:${user.no:}--${user.name}", properties))
	assert.Equal(t, "你好:${user.no}--Arvin", ResolvePlaceholders("你好:${user.no}--${user.name}", properties))
}
