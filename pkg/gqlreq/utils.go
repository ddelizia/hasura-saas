package gqlreq

import (
	"encoding/json"
	"time"

	"github.com/joomcode/errorx"
)

func TimestampToGraphQlTimestamWithTimeszone(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	return t.Format(time.RFC3339)
}

func InterfaceToJson(data interface{}) (string, error) {
	jsonString, err := json.Marshal(data)
	if err != nil {
		return "", errorx.IllegalArgument.Wrap(err, "not able to marshal json")
	}
	return string(jsonString), nil
}
