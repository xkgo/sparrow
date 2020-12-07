package JsonUtils

import (
	"encoding/json"
	"fmt"
)

func ToJsonString(v interface{}) (content string, err error) {
	bytes, err := json.Marshal(v)

	if err == nil {
		content = string(bytes)
	}
	return
}

func ToJsonStringWithoutError(v interface{}) string {
	bytes, err := json.Marshal(v)

	if err == nil {
		return string(bytes)
	}
	return fmt.Sprint(v)
}

func FromJson(jsonStr string, v interface{}) (interface{}, error) {
	err := json.Unmarshal([]byte(jsonStr), v)
	return v, err
}
