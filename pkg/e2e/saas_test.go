package e2e_test

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/ddelizia/hasura-saas/pkg/e2e"
	"github.com/ddelizia/hasura-saas/pkg/gqlreq"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/simonnilsson/ask"
	"github.com/sirupsen/logrus"
)

func setAccount(user, role, sourceAccountName, targetAccountName string) (map[string]interface{}, error) {
	bodySet := map[string]interface{}{}
	err := e2e.GraphqlService.Execute(
		context.Background(),
		fmt.Sprintf(`
			mutation SetCurrentAccount {
				saas_set_current_account ( data: { id_account: "%0s" } ) {
					id_account,
				}
			}`, e2e.CACHE_ACCOUNT_NAME_TO_ID[targetAccountName]),
		[]gqlreq.RequestHeader{
			{Key: "x-hasura-account-id", Value: e2e.CACHE_ACCOUNT_NAME_TO_ID[sourceAccountName]},
			{Key: "x-hasura-role", Value: role},
			{Key: "x-hasura-user-id", Value: user},
		},
		[]gqlreq.RequestVar{},
		true,
		&bodySet,
	)

	return bodySet, err
}

func getAccountInfo(id_user string) (map[string]interface{}, error) {
	bodyGet := map[string]interface{}{}
	err := e2e.GraphqlService.Execute(
		context.Background(),
		fmt.Sprintf(`
			query GetCurrentAccount {
				saas_get_current_account ( data: { id_user: "%0s" } ) {
					id_account,
					id_role
				}
			}`, id_user),
		[]gqlreq.RequestHeader{},
		[]gqlreq.RequestVar{},
		true,
		&bodyGet,
	)

	return bodyGet, err
}

var _ = Describe("saas e2e", func() {

	logrus.SetOutput(ioutil.Discard)

	BeforeEach(func() {
		e2e.DeleteTestData()
		e2e.CreateTestData()
	})

	AfterEach(func() {
		e2e.DeleteTestData()
	})

	It("should be able to get the current account when it is set", func() {

		_, errSet := setAccount(e2e.SHARED_USER_ACCOUNT_01_02_OWNER_ADMIN, e2e.ROLE_ACCOUNT_OWNER, e2e.ACCOUNT_01, e2e.ACCOUNT_02)

		bodyGet, errGet := getAccountInfo(e2e.SHARED_USER_ACCOUNT_01_02_OWNER_ADMIN)

		accountID, _ := ask.For(bodyGet, "saas_get_current_account.id_account").String("")

		Expect(errSet).To(BeNil())
		Expect(errGet).To(BeNil())
		fmt.Printf("Cache %v", e2e.CACHE_ACCOUNT_NAME_TO_ID)
		Expect(accountID).To(Equal(e2e.CACHE_ACCOUNT_NAME_TO_ID[e2e.ACCOUNT_02]))
	})

	It("should not be able change to an account that you do not belong to", func() {

		_, err := setAccount(e2e.SHARED_USER_ACCOUNT_01_02_OWNER_ADMIN, e2e.ROLE_ACCOUNT_OWNER, e2e.ACCOUNT_01, e2e.ACCOUNT_03)

		Expect(err).ToNot(BeNil())
	})

})
