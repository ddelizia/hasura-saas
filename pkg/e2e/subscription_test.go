package e2e_test

import (
	"context"
	"fmt"
	"os"

	"github.com/ddelizia/hasura-saas/pkg/gqlreq"
	"github.com/ddelizia/hasura-saas/pkg/logger"
	"github.com/ddelizia/hasura-saas/pkg/subscription"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/simonnilsson/ask"
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/paymentmethod"
)

const (
	GRAPH_QL_URL = "http://localhost:8082/v1/graphql"
)

var (
	graphqlService = gqlreq.NewService()
)

func generateAccountName() string {
	return "subscription_test_account_" + uuid.NewString()
}

func createPaymentMethod(number, expMonth, expYear, cvc string) string {

	os.Setenv("SUBSCRIPTION.STRIPE.API_KEY", os.Getenv("STRIPE_KEY"))

	params := &stripe.PaymentMethodParams{
		Type: stripe.String("card"),
		Card: &stripe.PaymentMethodCardParams{
			Number:   stripe.String(number),
			ExpMonth: stripe.String(expMonth),
			ExpYear:  stripe.String(expYear),
			CVC:      stripe.String(cvc),
		},
		BillingDetails: &stripe.BillingDetailsParams{},
	}
	paymentMethodResp, _ := paymentmethod.New(params)

	logrus.WithField("stripe.paymentmethod", logger.PrintStruct(paymentMethodResp)).Info("payment created")

	return paymentMethodResp.ID
}

func initSubscription() map[string]interface{} {
	bodyInit := map[string]interface{}{}

	err := graphqlService.Execute(
		context.Background(),
		fmt.Sprintf(`
		mutation InitSubscription {
			init_subscription ( data: { account_name: "%0s", id_plan: "%1s" } ) {
				account_id
			}
		}`, generateAccountName(), "basic"),
		[]gqlreq.RequestHeader{
			{Key: "x-hasura-account-id", Value: "no-account"},
			{Key: "x-hasura-role", Value: "logged_in"},
			{Key: "x-hasura-user-id", Value: "user1"},
		},
		[]gqlreq.RequestVar{},
		true,
		&bodyInit,
	)
	accountID, _ := ask.For(bodyInit, "init_subscription.account_id").String("")

	Expect(err).To(BeNil())
	Expect(accountID).NotTo(BeEmpty())

	return bodyInit
}

func createSubscription(a, cardN, cardM, cardY, cardCVC string) map[string]interface{} {
	bodyCreate := map[string]interface{}{}

	err := graphqlService.Execute(
		context.Background(),
		fmt.Sprintf(`
		mutation CreateSubscription {
			create_subscription ( data: { payment_method_id: "%0s" } ) {
				account_id,
				is_active
			}
		}`, createPaymentMethod(cardN, cardM, cardY, cardCVC)),
		[]gqlreq.RequestHeader{
			{Key: "x-hasura-account-id", Value: a},
			{Key: "x-hasura-role", Value: "account_owner"},
			{Key: "x-hasura-user-id", Value: "user1"},
		},
		[]gqlreq.RequestVar{},
		true,
		&bodyCreate,
	)
	accountID, _ := ask.For(bodyCreate, "create_subscription.account_id").String("")

	// Then
	Expect(err).To(BeNil())
	Expect(accountID).ToNot(BeEmpty())

	return bodyCreate
}

var _ = Describe("Subscription e2e", func() {

	stripe.Key = subscription.ConfigApiKey()
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
	})

	os.Setenv("GRAPHQL.HASURA.ADMINSECRET", os.Getenv("HASURA_GRAPHQL_ADMIN_SECRET"))
	os.Setenv("SUBSCRIPTION.STRIPE.APIKEY", os.Getenv("STRIPE_KEY"))
	os.Setenv("SUBSCRIPTION.STRIPE.WEBHOOKSECRET", os.Getenv("STRIPE_WEBHOOK_SECRET"))

	It("I should be able to execute complete payment flow with the default card", func() {

		responseInit := initSubscription()
		accountID, _ := ask.For(responseInit, "init_subscription.account_id").String("")

		responseCreate := createSubscription(accountID, "4242424242424242", "01", "2030", "314")
		isActive, _ := ask.For(responseCreate, "create_subscription.is_active").Bool(false)

		Expect(isActive).To(BeTrue())

	})

	It("I should be able to execute a payment with 3d authentication", func() {

		responseInit := initSubscription()
		accountID, _ := ask.For(responseInit, "init_subscription.account_id").String("")

		responseCreate := createSubscription(accountID, "4000002760003184", "01", "2030", "314")
		isActive, _ := ask.For(responseCreate, "create_subscription.is_active").Bool(true)

		Expect(isActive).To(BeFalse())
	})

	It("I should be able to crete and cancel subscription", func() {

		responseInit := initSubscription()
		accountID, _ := ask.For(responseInit, "init_subscription.account_id").String("")

		_ = createSubscription(accountID, "4242424242424242", "01", "2030", "314")

		bodyCancel := map[string]interface{}{}

		err := graphqlService.Execute(
			context.Background(),
			`mutation CancelSubscription {
				cancel_subscription {
					status
				}
			}`,
			[]gqlreq.RequestHeader{
				{Key: "x-hasura-account-id", Value: accountID},
				{Key: "x-hasura-role", Value: "account_owner"},
				{Key: "x-hasura-user-id", Value: "user1"},
			},
			[]gqlreq.RequestVar{},
			true,
			&bodyCancel,
		)
		status, _ := ask.For(bodyCancel, "cancel_subscription.status").String("")

		Expect(err).To(BeNil())
		Expect(status).To(Equal("canceled"))
	})

	It("I should be able to cancel subscription if it is not account owner", func() {

		responseInit := initSubscription()
		accountID, _ := ask.For(responseInit, "init_subscription.account_id").String("")

		_ = createSubscription(accountID, "4242424242424242", "01", "2030", "314")

		bodyCancel := map[string]interface{}{}

		err := graphqlService.Execute(
			context.Background(),
			`mutation CancelSubscription {
				cancel_subscription {
					status
				}
			}`,
			[]gqlreq.RequestHeader{
				{Key: "x-hasura-account-id", Value: accountID},
				{Key: "x-hasura-role", Value: "account_admin"},
				{Key: "x-hasura-user-id", Value: "user1"},
			},
			[]gqlreq.RequestVar{},
			true,
			&bodyCancel,
		)

		Expect(err).ToNot(BeNil())
	})

})
