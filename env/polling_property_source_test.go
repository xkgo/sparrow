package env

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestPollingPropertySource_Subscribe(t *testing.T) {

	rand.Seed(time.Now().Unix())
	var reader PropertyReader = NewPropertyReader(func() (kvs map[string]string, err error) {
		return map[string]string{
			"name": strconv.FormatInt(int64(rand.Int()), 10),
		}, nil
	})

	source, _ := NewPollingPropertySource("test", 1, reader)

	source.Subscribe("", func(event *KeyChangeEvent) {
		fmt.Println(event.Key + ":" + event.Nv)
	})

	for {
		time.Sleep(time.Duration(10) * time.Second)
	}

}
