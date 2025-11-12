package utils

import "encoding/json"

func ToJsonString(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}
