package utils

import "encoding/json"

func AsPrettyJson(input interface{}) string {
	jsonB, _ := json.MarshalIndent(input, "", "  ")
	return string(jsonB)
}
