package logger

import (
	"github.com/sirupsen/logrus"
)

func IntLogger() {

	l, err := logrus.ParseLevel(ConfigLogLevel())
	logrus.SetReportCaller(true)

	if err != nil {
		logrus.WithError(err).WithField("log_level", ConfigLogLevel()).Error("level not found setting standard logrus level")
	} else {
		logrus.WithField("log_level", ConfigLogLevel()).Info("log level initialized")
		logrus.SetLevel(l)
	}
}
