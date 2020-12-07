package env

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResolveProfileDirs(t *testing.T) {
	dirs := resolveProfileDirs([]string{"../testdata"})
	assert.ElementsMatch(t, []string{"../testdata"}, dirs)

	dirs = resolveProfileDirs([]string{"./noexists"})
	assert.Empty(t, dirs)
}

func TestGetFirstDefaultApplicationProfileInfo(t *testing.T) {

	pi := getFirstDefaultApplicationProfileInfo([]string{"../testdata"})

	fmt.Println(pi.profile)
	fmt.Println(pi.path)
	fmt.Println(pi.extension)
}

func TestGetNotDefaultProfileInfoWithExtension(t *testing.T) {
	pis := getNotDefaultProfileInfoWithExtension([]string{"../testdata"}, "")

	fmt.Println(pis)
}

func TestReadLocalFileAsPropertySource(t *testing.T) {
	source, _ := ReadLocalFileAsPropertySource("test", "../testdata/application.properties")
	source.Each(func(key, value string) (stop bool) {
		fmt.Println(key, ":", value)
		return false
	})
	fmt.Println("========================================================================================================================")
	source, _ = ReadLocalFileAsPropertySource("test", "../testdata/yml/application.yml")
	source.Each(func(key, value string) (stop bool) {
		fmt.Println(key, ":", value)
		return false
	})

	fmt.Println("========================================================================================================================")
	source, err := ReadLocalFileAsPropertySource("test", "../testdata/toml/application.toml")
	fmt.Println(err)
	source.Each(func(key, value string) (stop bool) {
		fmt.Println(key, ":", value)
		return false
	})

	fmt.Println("========================================================================================================================")
}
