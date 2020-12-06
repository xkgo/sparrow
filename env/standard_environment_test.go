package env

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestStandardEnvironment_New(t *testing.T) {
	fmt.Println(filepath.Abs("./"))
}
