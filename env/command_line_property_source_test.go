package env

import (
	"fmt"
	"testing"
)

func TestNewCommandLinePropertySource(t *testing.T) {
	commandLine := "--app.name=Test --active.profiles=dev,wuxi --equal=1-2=x"

	source := NewCommandLinePropertySource(commandLine)

	source.Each(func(key string, value string) (stop bool) {
		fmt.Println(key, ":", value)
		return false
	})

}
