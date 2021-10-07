package e2e_test

import (
	"context"
	"fmt"

	"github.com/ddelizia/hasura-saas/pkg/e2e"
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

func generateAccountName() string {
	return "subscription_test_account_" + uuid.NewString()
}

func createPaymentMethod(number, expMonth, expYear, cvc string) string {

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

	err := e2e.GraphqlService.Execute(
		context.Background(),
		fmt.Sprintf(`
		mutation InitSubscription {
			subscription_init ( data: { account_name: "%0s", id_plan: "%1s" } ) {
				id_account
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
	accountID, _ := ask.For(bodyInit, "subscription_init.id_account").String("")

	Expect(err).To(BeNil())
	Expect(accountID).NotTo(BeEmpty())

	return bodyInit
}

func createSubscription(a, cardN, cardM, cardY, cardCVC string) map[string]interface{} {
	bodyCreate := map[string]interface{}{}

	err := e2e.GraphqlService.Execute(
		context.Background(),
		fmt.Sprintf(`
		mutation CreateSubscription {
			subscription_create ( data: { payment_method_id: "%0s" } ) {
				id_account,
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
	accountID, _ := ask.For(bodyCreate, "subscription_create.id_account").String("")

	// Then
	Expect(err).To(BeNil())
	Expect(accountID).ToNot(BeEmpty())

	return bodyCreate
}

var _ = Describe("subscription e2e", func() {

	stripe.Key = subscription.ConfigApiKey()
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
	})

	It("should be able to execute complete payment flow with the default card", func() {

		responseInit := initSubscription()
		accountID, _ := ask.For(responseInit, "subscription_init.id_account").String("")

		responseCreate := createSubscription(accountID, "4242424242424242", "01", "2030", "314")
		isActive, _ := ask.For(responseCreate, "subscription_create.is_active").Bool(false)

		Expect(isActive).To(BeTrue())

	})

	It("should be able to execute a payment with 3d authentication", func() {

		responseInit := initSubscription()
		accountID, _ := ask.For(responseInit, "subscription_init.id_account").String("")

		responseCreate := createSubscription(accountID, "4000002760003184", "01", "2030", "314")
		isActive, _ := ask.For(responseCreate, "subscription_create.is_active").Bool(true)

		Expect(isActive).To(BeFalse())
	})

	It("should be able to retry when payment is not successful", func() {

		responseInit := initSubscription()
		accountID, _ := ask.For(responseInit, "subscription_init.id_account").String("")

		responseCreate := createSubscription(accountID, "4000002760003184", "01", "2030", "314")
		isActiveCreate, _ := ask.For(responseCreate, "subscription_create.is_active").Bool(true)
		accountCreate, _ := ask.For(responseCreate, "subscription_create.id_account").String("id_account_create")

		bodyRetry := map[string]interface{}{}
		err := e2e.GraphqlService.Execute(
			context.Background(),
			fmt.Sprintf(`
			mutation RetrySubscription {
				subscription_retry ( data: { payment_method_id: "%0s" } ) {
					id_account,
					is_active
				}
			}`, createPaymentMethod("4242424242424242", "01", "2030", "314")),
			[]gqlreq.RequestHeader{
				{Key: "x-hasura-account-id", Value: accountID},
				{Key: "x-hasura-role", Value: "account_owner"},
				{Key: "x-hasura-user-id", Value: "user1"},
			},
			[]gqlreq.RequestVar{},
			true,
			&bodyRetry,
		)
		isActiveRetry, _ := ask.For(bodyRetry, "subscription_retry.is_active").Bool(false)
		accountRetry, _ := ask.For(bodyRetry, "subscription_retry.id_account").String("id_account_retry")

		Expect(err).To(BeNil())
		Expect(isActiveCreate).To(BeFalse())
		Expect(isActiveRetry).To(BeTrue())
		Expect(accountCreate).To(Equal(accountRetry))
	})

	It("should be able to change an existing plan", func() {

		responseInit := initSubscription()
		accountID, _ := ask.For(responseInit, "subscription_init.id_account").String("")

		responseCreate := createSubscription(accountID, "4242424242424242", "01", "2030", "314")
		isActiveCreate, _ := ask.For(responseCreate, "subscription_create.is_active").Bool(false)
		accountCreate, _ := ask.For(responseCreate, "subscription_create.id_account").String("id_account_create")

		bodyChange := map[string]interface{}{}
		err := e2e.GraphqlService.Execute(
			context.Background(),
			`mutation ChangeSubscription {
				subscription_change ( data: { id_plan: "premium" } ) {
					id_account,
					is_active
				}
			}`,
			[]gqlreq.RequestHeader{
				{Key: "x-hasura-account-id", Value: accountID},
				{Key: "x-hasura-role", Value: "account_owner"},
				{Key: "x-hasura-user-id", Value: "user1"},
			},
			[]gqlreq.RequestVar{},
			true,
			&bodyChange,
		)
		isActiveRetry, _ := ask.For(bodyChange, "subscription_change.is_active").Bool(false)
		accountChange, _ := ask.For(bodyChange, "subscription_change.id_account").String("id_account_change")

		Expect(err).To(BeNil())
		Expect(isActiveCreate).To(BeTrue())
		Expect(isActiveRetry).To(BeTrue())
		Expect(accountCreate).To(Equal(accountChange))
	})

	It("should be able to crete and cancel subscription", func() {

		responseInit := initSubscription()
		accountID, _ := ask.For(responseInit, "subscription_init.id_account").String("")

		_ = createSubscription(accountID, "4242424242424242", "01", "2030", "314")

		bodyCancel := map[string]interface{}{}

		err := e2e.GraphqlService.Execute(
			context.Background(),
			`mutation CancelSubscription {
				subscription_cancel {
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
		status, _ := ask.For(bodyCancel, "subscription_cancel.status").String("")

		Expect(err).To(BeNil())
		Expect(status).To(Equal("canceled"))
	})

	It("should be able to cancel subscription if it is not account owner", func() {

		responseInit := initSubscription()
		accountID, _ := ask.For(responseInit, "subscription_init.id_account").String("")

		_ = createSubscription(accountID, "4242424242424242", "01", "2030", "314")

		bodyCancel := map[string]interface{}{}

		err := e2e.GraphqlService.Execute(
			context.Background(),
			`mutation CancelSubscription {
				subscription_cancel {
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
