package env

import (
	"fmt"
	"testing"
)

func TestResolveProfileDirs(t *testing.T) {
	dirs := resolveProfileDirs([]string{"../testdata"})

	fmt.Println(dirs)
}
