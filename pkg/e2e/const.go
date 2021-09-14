package e2e

import "github.com/ddelizia/hasura-saas/pkg/gqlreq"

type MembershipData struct {
	Role string
}

type CustomerData struct {
	Customer string
}

type AccountData struct {
	Memberships map[string]*MembershipData
	Customer    string
	Plan        string
	Status      string
}

type TestConfig struct {
	Accounts map[string]*AccountData
	Users    []string
}

const (
	TABLE_SAAS_ACCOUNT          = "saas_account"
	TABLE_SUBSCRIPTION_CUSTOMER = "subscription_customer"
	TABLE_SUBSCRIPTION_STATUS   = "subscription_status"

	USER_ACCOUNT_O1_01 = "TestUserAccount01_01"
	USER_ACCOUNT_O1_02 = "TestUserAccount01_02"
	USER_ACCOUNT_O1_03 = "TestUserAccount01_03"
	USER_ACCOUNT_O1_04 = "TestUserAccount01_04"

	USER_ACCOUNT_02_01 = "TestUserAccount02_01"
	USER_ACCOUNT_02_02 = "TestUserAccount02_02"
	USER_ACCOUNT_02_03 = "TestUserAccount02_03"
	USER_ACCOUNT_02_04 = "TestUserAccount02_04"

	USER_ACCOUNT_03_01 = "TestUserAccount03_01"
	USER_ACCOUNT_03_02 = "TestUserAccount03_02"
	USER_ACCOUNT_03_03 = "TestUserAccount03_03"
	USER_ACCOUNT_03_04 = "TestUserAccount03_04"

	ACCOUNT_01 = "TestAccount01"
	ACCOUNT_02 = "TestAccount02"
	ACCOUNT_03 = "TestAccount03"

	ROLE_ADMIN         = "admin"
	ROLE_ACCOUNT_OWNER = "account_owner"
	ROLE_ACCOUNT_ADMIN = "account_admin"
	ROLE_USER          = "account_user"

	PLAN_BASIC = "basic"

	STATUS_PAID = "PAID"
)

var (
	GqlService = gqlreq.NewService()

	TEST_CONFIGURATION = &TestConfig{
		Users: []string{
			USER_ACCOUNT_O1_01,
			USER_ACCOUNT_O1_02,
			USER_ACCOUNT_O1_03,
			USER_ACCOUNT_O1_04,
		},
		Accounts: map[string]*AccountData{
			ACCOUNT_01: {
				Memberships: map[string]*MembershipData{
					USER_ACCOUNT_O1_01: {Role: ROLE_ADMIN},
					USER_ACCOUNT_O1_02: {Role: ROLE_ACCOUNT_OWNER},
					USER_ACCOUNT_O1_03: {Role: ROLE_ACCOUNT_ADMIN},
					USER_ACCOUNT_O1_04: {Role: ROLE_USER},
				},
				Customer: "MockStripeCustomer" + ACCOUNT_01,
				Plan:     PLAN_BASIC,
				Status:   STATUS_PAID,
			},
			ACCOUNT_02: {
				Memberships: map[string]*MembershipData{
					USER_ACCOUNT_02_01: {Role: ROLE_ADMIN},
					USER_ACCOUNT_02_02: {Role: ROLE_ACCOUNT_OWNER},
					USER_ACCOUNT_02_03: {Role: ROLE_ACCOUNT_ADMIN},
					USER_ACCOUNT_02_04: {Role: ROLE_USER},
				},
				Customer: "MockStripeCustomer" + ACCOUNT_02,
				Plan:     PLAN_BASIC,
				Status:   STATUS_PAID,
			},
			ACCOUNT_03: {
				Memberships: map[string]*MembershipData{
					USER_ACCOUNT_03_01: {Role: ROLE_ADMIN},
					USER_ACCOUNT_03_02: {Role: ROLE_ACCOUNT_OWNER},
					USER_ACCOUNT_03_03: {Role: ROLE_ACCOUNT_ADMIN},
					USER_ACCOUNT_03_04: {Role: ROLE_USER},
				},
				Customer: "MockStripeCustomer" + ACCOUNT_03,
				Plan:     PLAN_BASIC,
				Status:   STATUS_PAID,
			},
		},
	}
)
