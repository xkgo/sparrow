package env

import (
	"fmt"
	"testing"
)

func TestNewSystemEnvironmentPropertySource(t *testing.T) {

	source := NewSystemEnvironmentPropertySource()

	source.Each(func(key string, value string) (stop bool) {
		fmt.Println(key, ":", value)
		return false
	})
}
