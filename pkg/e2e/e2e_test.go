package e2e_test

import (
	"os"
	"testing"

	"github.com/ddelizia/hasura-saas/pkg/e2e"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func Test_E2e(t *testing.T) {
	RegisterFailHandler(Fail)

	os.Setenv("GRAPHQL.HASURA.ADMINSECRET", os.Getenv("HASURA_GRAPHQL_ADMIN_SECRET"))
	os.Setenv("SUBSCRIPTION.STRIPE.APIKEY", os.Getenv("STRIPE_KEY"))
	os.Setenv("SUBSCRIPTION.STRIPE.WEBHOOKSECRET", os.Getenv("STRIPE_WEBHOOK_SECRET"))

	if os.Getenv("EXECUTE_E2E") != "true" {
		print("Skipping [e2e] tests! To run them set EXECUTE_E2E=true")
	} else {
		if os.Getenv("CREATE_INITIAL_DATA") == "true" {
			e2e.CreateTestData()
		}
		RunSpecs(t, "[e2e] test suite")
	}
}
