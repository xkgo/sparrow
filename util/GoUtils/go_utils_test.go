package GoUtils

import (
	"context"
	"fmt"
	"testing"
)

func TestGetContext(t *testing.T) {

	ctx := context.Background()
	BindContext(&ctx)

	ctx1 := GetContext()
	ctx2 := GetContext()

	fmt.Println("ctx:", &ctx, ", ctx1:", ctx1, ", ctx2:", ctx2)

}
