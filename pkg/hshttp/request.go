package hshttp

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/ddelizia/hasura-saas/pkg/logger"
	"github.com/joomcode/errorx"
	"github.com/sirupsen/logrus"
)

func GetBody(r *http.Request, data interface{}) error {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.WithError(err).WithField("payload", logger.PrintStruct(data)).Error("payload cannot be readed from bytes")
		return errorx.IllegalArgument.Wrap(err, "invalid payload")
	}

	err = json.Unmarshal(reqBody, data)
	if err != nil {
		logrus.WithError(err).WithField("payload", logger.PrintStruct(data)).Error("request body is not valid")
		return errorx.IllegalArgument.Wrap(err, "invalid json")
	}
	logrus.WithField("payload", logger.PrintStruct(data)).Debug("received request payload")
	return nil
}

func WriteBody(w http.ResponseWriter, data interface{}) error {

	d, err := json.Marshal(data)
	if err != nil {
		logrus.WithError(err).WithField("body", logger.PrintStruct(data)).Error("not able to generate response")
		return errorx.InternalError.Wrap(err, "not able to create output")
	}
	w.Write(d)
	logrus.WithField("body", logger.PrintStruct(data)).Debug("sending body")
	return nil
}
