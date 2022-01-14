package hserrorx

import (

	"github.com/joomcode/errorx"
)

func AssertTrue(isTrue bool, fields Fields, message string, privateMessage *string) error {
	if !isTrue {
		return New(errorx.IllegalState, fields, message, privateMessage)
	}
	return nil
}
