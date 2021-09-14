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

	if os.Getenv("EXECUTE_E2E") != "true" {
		print("Skipping [e2e] tests! To run them set EXECUTE_E2E=true")
	} else {
		RunSpecs(t, "[e2e] test suite")
		if os.Getenv("CREATE_INITIAL_DATA") == "true" {
			e2e.CreateTestConfiguration()
		}
	}
}
