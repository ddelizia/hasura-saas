package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

var EVENTS = []string{
	"charge.succeeded",
	"customer.subscription.created",
	"invoice.payment_action_required",
	"invoice.payment_failed",
	"customer.created",
	"customer.subscription.deleted",
	"customer.updated",
	"invoice.created",
	"invoice.finalized",
	"invoice.paid",
	"invoice.payment_succeeded",
	"payment_intent.created",
	"payment_intent.requires_action",
	"payment_intent.succeeded",
	"payment_method.attached",
}

const PATH string = "pkg/subscription"

const TEMPLATE_FILE string = `package subscription

import (
	"github.com/stripe/stripe-go"
)

func Process{{ .EventName }}(event stripe.Event, id string) error {
	data := &stripe.{{ .TypeName }}{}
	if err := beforeEvent(event, data); err != nil {
		return err
	}

	return nil
}
`

const TEMPLATE_MAPPING string = `package subscription

import (
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go"
)

// Autogenerated by ./cmd/stripe_event_code_generator/main.go

func EventMapping(event stripe.Event, id string) {
	switch event.Type {
	
	{{ range $key, $value := . }}
	case "{{ $key }}":
		Process{{ $value }}(event, id)
	{{- end }}
	
	default:
		// unhandled event type
		logrus.WithField("eventType", event.Type).Warn("event not mapped")
	}
}
`

type DataModel struct {
	TypeName  string
	EventName string
}

func main() {

	mapping := map[string]string{}

	for _, e := range EVENTS {
		typeName := strings.Replace(strings.Title(strings.Replace(strings.Split(e, ".")[0], "_", " ", -1)), " ", "", -1)
		fileNameSuffix := strings.Replace(e, ".", "_", -1)
		fileName := "event_" + fileNameSuffix + ".go"
		eventName := strings.Replace(strings.Title(strings.Replace(fileNameSuffix, "_", " ", -1)), " ", "", -1)
		path := PATH + "/" + fileName

		if _, err := os.Stat(path); os.IsNotExist(err) {

			tmpl, _ := template.New(fileNameSuffix).Parse(TEMPLATE_FILE)

			f, _ := os.Create(path)

			data := DataModel{
				TypeName:  typeName,
				EventName: eventName,
			}
			fmt.Printf("Creating data for %#v\n", data)
			_ = tmpl.Execute(f, data)
		}

		mapping[e] = eventName
	}

	tmplMapping, _ := template.New("mapping").Parse(TEMPLATE_MAPPING)
	f, _ := os.Create(PATH + "/event_mapping.go")
	_ = tmplMapping.Execute(f, mapping)

}
