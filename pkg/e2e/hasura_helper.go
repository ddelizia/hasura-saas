package e2e

import (
	"bytes"
	"context"
	"text/template"

	"github.com/ddelizia/hasura-saas/pkg/gqlsdk"
	"github.com/onsi/ginkgo"
	"github.com/sirupsen/logrus"
)

type SaasAccountData struct {
	SaasAccount []*gqlsdk.SaasAccount `json:"saas_account"`
}

func GetAllTestAccounts() map[string]string {
	result := &SaasAccountData{}
	err := GqlService.Execute(
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

func CreateTestConfiguration() {

	logrus.Info("Creating test data")

	tmpl, err := template.New("TemplateInitQuery").Parse(`
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
			affected_rows
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

	var result map[string]interface{}
	err = GqlService.Execute(
		context.Background(),
		mutation,
		nil,
		nil,
		true,
		result)

	if err != nil {
		ginkgo.Fail("error executing creation of account query " + err.Error() + " query: \n" + mutation)
	}
}
