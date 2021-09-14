package logger

import (
	"encoding/json"
	"fmt"
)

func PrintStruct(o interface{}) string {
	empJSON, err := json.Marshal(o)
	if err != nil {
		return fmt.Sprintf("error while marshalling struct %#v", o)
	}
	return string(empJSON)
}
