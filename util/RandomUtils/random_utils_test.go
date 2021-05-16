package RandomUtils

import (
	"fmt"
	"testing"
)

func TestRandomString(t *testing.T) {

	fmt.Println(RandomString(10, true, true, true))
	fmt.Println(RandomString(10, true, false, false))
	fmt.Println(RandomString(10, false, true, false))
	fmt.Println(RandomString(10, false, false, true))

}
