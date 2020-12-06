package deploy

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetInfo(t *testing.T) {

	info := GetInfo()
	assert.Equal(t, Dev, info.Env)

	AddFirst("Test", func() *Info {
		return &Info{
			Env:        Test,
			Properties: nil,
		}
	})

	info = GetInfo()
	assert.Equal(t, Test, info.Env)

}
