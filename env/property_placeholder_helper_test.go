package env

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPropertyPlaceholderHelper_ReplacePlaceholders(t *testing.T) {
	helper := NewPropertyPlaceholderHelper(DefaultPlaceholderPrefix, DefaultPlaceholderSuffix, DefaultPlaceholderValueSeparator, true)

	properties := map[string]string{
		"circular.var1": "${circular.var2}",
		"circular.var2": "${circular.var1}",
		"user.name":     "Arvin",
	}

	placeholderResolver := func(key string) string {
		return properties[key]
	}
	// 循环引用
	assert.Panics(t, func() { helper.ReplacePlaceholders("${circular.var1}", placeholderResolver) })

	assert.Equal(t, "你好:Arvin", helper.ReplacePlaceholders("你好:${user.name}", placeholderResolver))
	assert.Equal(t, "你好:Go", helper.ReplacePlaceholders("你好:${user.no:Go}", placeholderResolver))
	assert.Equal(t, "你好:Arvin", helper.ReplacePlaceholders("你好:${user.no:${user.name}}", placeholderResolver))
	assert.Equal(t, "你好:", helper.ReplacePlaceholders("你好:${user.no:}", placeholderResolver))
	assert.Equal(t, "你好:--Arvin", helper.ReplacePlaceholders("你好:${user.no:}--${user.name}", placeholderResolver))
	assert.Equal(t, "你好:${user.no}--Arvin", helper.ReplacePlaceholders("你好:${user.no}--${user.name}", placeholderResolver))
}

func TestPropertyPlaceholderHelper2_ReplacePlaceholders(t *testing.T) {
	helper := NewPropertyPlaceholderHelper("#{", "}", "::", true)

	properties := map[string]string{
		"circular.var1": "#{circular.var2}",
		"circular.var2": "#{circular.var1}",
		"user.name":     "Arvin",
	}

	placeholderResolver := func(key string) string {
		return properties[key]
	}
	// 循环引用
	assert.Panics(t, func() { helper.ReplacePlaceholders("#{circular.var1}", placeholderResolver) })

	assert.Equal(t, "你好:Arvin", helper.ReplacePlaceholders("你好:#{user.name}", placeholderResolver))
	assert.Equal(t, "你好:Go", helper.ReplacePlaceholders("你好:#{user.no::Go}", placeholderResolver))
	assert.Equal(t, "你好:Arvin", helper.ReplacePlaceholders("你好:#{user.no::#{user.name}}", placeholderResolver))
	assert.Equal(t, "你好:", helper.ReplacePlaceholders("你好:#{user.no::}", placeholderResolver))
	assert.Equal(t, "你好:--Arvin", helper.ReplacePlaceholders("你好:#{user.no::}--#{user.name}", placeholderResolver))
	assert.Equal(t, "你好:#{user.no}--Arvin", helper.ReplacePlaceholders("你好:#{user.no}--#{user.name}", placeholderResolver))
}
