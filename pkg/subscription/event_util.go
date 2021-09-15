package subscription

import (
	"encoding/json"
	"fmt"

	"github.com/ddelizia/hasura-saas/pkg/logger"
	"github.com/joomcode/errorx"
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go"
)

func beforeEvent(event stripe.Event, data interface{}) error {
	err := json.Unmarshal(event.Data.Raw, &data)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"eventType": event.Type,
			"eventData": logger.PrintStruct(event.Data.Raw),
		}).Error("error while processing the event")
		return errorx.InternalError.Wrap(err, fmt.Sprintf("not able to process event of type %s", event.Type))
	}
	return nil
}
