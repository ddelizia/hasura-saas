package hserrorx

import (
	"github.com/joomcode/errorx"
	"github.com/sirupsen/logrus"
)

type Fields = map[string]interface{}

func Wrap(e error, errorType *errorx.Type, withFields Fields, message string, privateMessage *string) error {
	log := logrus.WithError(e).WithFields(withFields)
	if privateMessage != nil {
		log.Error(*privateMessage)
	} else {
		log.Error(message)
	}
	return withProperties(errorType.New(message), withFields)
}

func New(errorType *errorx.Type, withFields Fields, message string, privateMessage *string) error {
	log := logrus.WithFields(withFields)
	if privateMessage != nil {
		log.Error(*privateMessage)
	} else {
		log.Error(message)
	}
	return withProperties(errorType.New(message), withFields)
}

func withProperties(e *errorx.Error, withFields Fields) error {
	for k, v := range withFields {
		e = e.WithProperty(errorx.RegisterPrintableProperty(k), v)
	}
	return e
}
