package e2e

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/ddelizia/hasura-saas/pkg/gqlreq"
	"github.com/ddelizia/hasura-saas/pkg/gqlsdk"
	"github.com/onsi/ginkgo"
	"github.com/simonnilsson/ask"
	"github.com/sirupsen/logrus"
)

type SaasAccountData struct {
	SaasAccount []*gqlsdk.SaasAccount `json:"saas_account"`
}

func GetAllTestAccounts() map[string]string {
	result := &SaasAccountData{}
	err := GraphqlService.Execute(
		context.Background(),
		`query GetAllAccountsStartingWithTest($_like: String = "Test%") {
			saas_account(where: {name: {_like: $_like}}) {
				id
				name
			}
		}`,
		nil,
		nil,
		true,
		result)

	if err != nil {
		ginkgo.Fail("error getting info account " + err.Error())
	}

	response := map[string]string{}
	for _, account := range result.SaasAccount {
		response[account.Name] = account.ID
	}

	return response
}

func DeleteTestData() {
	logrus.Info("Delete test data")
	CACHE_ACCOUNT_NAME_TO_ID = map[string]string{}
	CACHE_ACCOUNT_ID_TO_NAME = map[string]string{}

	var resultDeleteMembership map[string]interface{}
	err := GraphqlService.Execute(
		context.Background(),
		fmt.Sprintf(`
			mutation DeleteByPrefix {
				delete_saas_membership(where: {id_user: {_like: "%s"}}) {
					affected_rows
				}
			}
			`, USER_PREFIX+"%"),
		nil,
		nil,
		true,
		resultDeleteMembership)

	if err != nil {
		ginkgo.Fail("failed to remove all members")
	}

	var result map[string]interface{}
	err = GraphqlService.Execute(
		context.Background(),
		fmt.Sprintf(`
		mutation DeleteByPrefix {
			delete_saas_account(where: {name: {_like: "%s"}}) {
				affected_rows
			}
		}
		`, ACCOUNT_PREFIX+"%"),
		nil,
		nil,
		true,
		result)

	if err != nil {
		ginkgo.Fail("failed to remove all accounts with prefix: " + ACCOUNT_PREFIX)
	}
}

func CreateTestData() {

	logrus.Info("Creating test data")

	tmpl, err := template.New("TemplateCreateTestData").Parse(`
	mutation MutationInsertInitialData {
		
		insert_saas_account(objects: [
			{{ range $account, $account_config :=  .Accounts -}}
			{
				name: "{{ $account }}", 
				
				saas_memberships: {
					data: [
						{{ range $user_membership, $membership_config := $account_config.Memberships -}}
							{id_role: "{{ $membership_config.Role }}", id_user: "{{ $user_membership }}"}, 
						{{ end }}
					]
				},
				
				subscription_customer: {
					data: {stripe_customer: "{{ $account_config.Customer}}"}
				},
				
				subscription_status: {
					data: {id_plan: "{{ $account_config.Plan }}", status: "{{ $account_config.Status }}"}
				}
			},
			{{- end }}
		]) 
		{
			affected_rows,
			returning {
				id
				name
			}
		}
	}
	`)
	if err != nil {
		ginkgo.Fail("unable to parse template " + err.Error())
	}

	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, TEST_CONFIGURATION)
	if err != nil {
		ginkgo.Fail("error while applying template " + err.Error())
	}
	mutation := tpl.String()
	logrus.Debug("importing initial data: " + mutation)

	result := map[string]interface{}{}
	err = GraphqlService.Execute(
		context.Background(),
		mutation,
		[]gqlreq.RequestHeader{},
		[]gqlreq.RequestVar{},
		true,
		&result)

	if err != nil {
		ginkgo.Fail("while executing mutation MutationInsertInitialData" + err.Error())
	}

	def := []interface{}{}
	returnResult, isSuccessfult := ask.For(result, "insert_saas_account.returning").Slice(def)

	if !isSuccessfult {
		ginkgo.Fail("not able to parse data")
	}

	for _, v := range returnResult {
		val := v.(map[string]interface{})
		CACHE_ACCOUNT_NAME_TO_ID[val["name"].(string)] = val["id"].(string)
		CACHE_ACCOUNT_ID_TO_NAME[val["id"].(string)] = val["name"].(string)
	}

	if err != nil {
		ginkgo.Fail("error executing creation of account query " + err.Error() + " query: \n" + mutation)
	}
}
