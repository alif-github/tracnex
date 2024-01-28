package util

import (
	"encoding/json"
	"fmt"
)

func StructToJSON(input interface{}) (output string) {
	b, err := json.Marshal(input)
	if err != nil {
		fmt.Println(err)
		output = ""
		return
	}
	output = string(b)
	return
}

func StringTruncate(data string, maxLength int) string {
	if len(data) < maxLength {
		return data
	} else {
		return data[0 : maxLength-1]
	}
}
