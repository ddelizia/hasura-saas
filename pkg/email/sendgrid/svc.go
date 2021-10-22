package sendgrid

import (
	"encoding/json"

	"github.com/ddelizia/hasura-saas/pkg/email"
	"github.com/ddelizia/hasura-saas/pkg/hserrorx"
	"github.com/ddelizia/hasura-saas/pkg/logger"
	"github.com/joomcode/errorx"
	"github.com/sendgrid/sendgrid-go"
)

type service struct {
	ApiKey      string
	SandboxMode bool
}

func NewService() email.Service {
	return &service{
		ApiKey: ConfigSendgridApiKey(),
	}
}

func NewServiceSandbox() email.Service {
	return &service{
		ApiKey:      ConfigSendgridApiKey(),
		SandboxMode: true,
	}
}

func (s *service) SendEmail(emails []string, subject, template string, data interface{}) error {

	tos := []*To{}
	for _, e := range emails {
		to := &To{Email: e}
		tos = append(tos, to)
	}
	emailData := &SendgridEmailData{
		From: &From{
			Email: email.ConfigFrom(),
		},
		TemplateID: template,
		Personalizations: []*Personalizations{
			{
				To:                  tos,
				DynamicTemplateData: data,
			},
		},
		MailSettings: &MailSettings{
			SandboxMode: &SandboxMode{
				Enable: s.SandboxMode,
			},
		},
	}

	body, err := json.Marshal(emailData)
	if err != nil {
		return hserrorx.Wrap(
			err,
			errorx.InternalError,
			hserrorx.Fields{
				LOG_PARAM_SENDGRID_REQUEST: logger.PrintStruct(emailData),
			},
			"not able to process data for email", nil)
	}

	request := sendgrid.GetRequest(s.ApiKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = body

	response, err := sendgrid.API(request)

	if err != nil {
		return hserrorx.Wrap(
			err,
			errorx.InternalError,
			hserrorx.Fields{
				LOG_PARAM_SENDGRID_REQUEST:  logger.PrintStruct(emailData),
				LOG_PARAM_SENDGRID_RESPONSE: logger.PrintStruct(response),
			},
			"not able to send email", nil)
	}

	if response.StatusCode > 299 || response.StatusCode < 200 {
		return hserrorx.New(
			errorx.InternalError,
			hserrorx.Fields{
				"ApiKey":                    s.ApiKey,
				LOG_PARAM_SENDGRID_REQUEST:  logger.PrintStruct(emailData),
				LOG_PARAM_SENDGRID_RESPONSE: logger.PrintStruct(response),
			},
			"provider responded with invalid status code", nil)
	}

	return nil
}
