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
	ACCOUNT_PREFIX              = "TestAccount"
	USER_PREFIX                 = "TestUser"

	USER_ON_ACCOUNT_O1_ACCOUNT_OWNER = USER_PREFIX + "_Account_01_account_owner"
	USER_ON_ACCOUNT_O1_ACCOUNT_ADMIN = USER_PREFIX + "_Account_01_account_admin"
	USER_ON_ACCOUNT_O1_USER          = USER_PREFIX + "_Account_01_account_user"

	USER_ON_ACCOUNT_02_ACCOUNT_OWNER = USER_PREFIX + "_Account_02_account_owner"
	USER_ON_ACCOUNT_02_ACCOUNT_ADMIN = USER_PREFIX + "_Account_02_account_admin"
	USER_ON_ACCOUNT_02_USER          = USER_PREFIX + "_Account_02_account_user"

	USER_ON_ACCOUNT_03_ACCOUNT_OWNER = USER_PREFIX + "_Account_03_account_owner"
	USER_ON_ACCOUNT_03_ACCOUNT_ADMIN = USER_PREFIX + "_Account_03_account_admin"
	USER_ON_ACCOUNT_03_USER          = USER_PREFIX + "_Account_03_account_user"

	SHARED_USER_ACCOUNT_01_02_OWNER_ADMIN = USER_PREFIX + "_Account_01_02_owner_admin"
	SHARED_USER_ACCOUNT_01_03_ADMIN_USER  = USER_PREFIX + "_Account_01_03_admin_user"

	ACCOUNT_01 = ACCOUNT_PREFIX + "01"
	ACCOUNT_02 = ACCOUNT_PREFIX + "02"
	ACCOUNT_03 = ACCOUNT_PREFIX + "03"

	ROLE_ADMIN         = "admin"
	ROLE_ACCOUNT_OWNER = "account_owner"
	ROLE_ACCOUNT_ADMIN = "account_admin"
	ROLE_USER          = "account_user"

	PLAN_BASIC   = "basic"
	PLAN_PREMIUM = "premium"

	STATUS_PAID = "PAID"
)

var (
	GraphqlService = gqlreq.NewService()

	CACHE_ACCOUNT_NAME_TO_ID = map[string]string{}
	CACHE_ACCOUNT_ID_TO_NAME = map[string]string{}

	TEST_CONFIGURATION = &TestConfig{
		Users: []string{
			USER_ON_ACCOUNT_O1_ACCOUNT_OWNER,
			USER_ON_ACCOUNT_O1_ACCOUNT_ADMIN,
			USER_ON_ACCOUNT_O1_USER,
			USER_ON_ACCOUNT_02_ACCOUNT_OWNER,
			USER_ON_ACCOUNT_02_ACCOUNT_ADMIN,
			USER_ON_ACCOUNT_02_USER,
			USER_ON_ACCOUNT_03_ACCOUNT_OWNER,
			USER_ON_ACCOUNT_03_ACCOUNT_ADMIN,
			USER_ON_ACCOUNT_03_USER,
			SHARED_USER_ACCOUNT_01_02_OWNER_ADMIN,
			SHARED_USER_ACCOUNT_01_03_ADMIN_USER,
		},
		Accounts: map[string]*AccountData{
			ACCOUNT_01: {
				Memberships: map[string]*MembershipData{
					USER_ON_ACCOUNT_O1_ACCOUNT_OWNER:      {Role: ROLE_ACCOUNT_OWNER},
					USER_ON_ACCOUNT_O1_ACCOUNT_ADMIN:      {Role: ROLE_ACCOUNT_ADMIN},
					USER_ON_ACCOUNT_O1_USER:               {Role: ROLE_USER},
					SHARED_USER_ACCOUNT_01_02_OWNER_ADMIN: {Role: ROLE_ACCOUNT_OWNER},
					SHARED_USER_ACCOUNT_01_03_ADMIN_USER:  {Role: ROLE_ACCOUNT_ADMIN},
				},
				Customer: "MockStripeCustomer" + ACCOUNT_01,
				Plan:     PLAN_BASIC,
				Status:   STATUS_PAID,
			},
			ACCOUNT_02: {
				Memberships: map[string]*MembershipData{
					USER_ON_ACCOUNT_02_ACCOUNT_OWNER:      {Role: ROLE_ACCOUNT_OWNER},
					USER_ON_ACCOUNT_02_ACCOUNT_ADMIN:      {Role: ROLE_ACCOUNT_ADMIN},
					USER_ON_ACCOUNT_02_USER:               {Role: ROLE_USER},
					SHARED_USER_ACCOUNT_01_02_OWNER_ADMIN: {Role: ROLE_ACCOUNT_ADMIN},
				},
				Customer: "MockStripeCustomer" + ACCOUNT_02,
				Plan:     PLAN_BASIC,
				Status:   STATUS_PAID,
			},
			ACCOUNT_03: {
				Memberships: map[string]*MembershipData{
					USER_ON_ACCOUNT_03_ACCOUNT_OWNER:     {Role: ROLE_ACCOUNT_OWNER},
					USER_ON_ACCOUNT_03_ACCOUNT_ADMIN:     {Role: ROLE_ACCOUNT_ADMIN},
					USER_ON_ACCOUNT_03_USER:              {Role: ROLE_USER},
					SHARED_USER_ACCOUNT_01_03_ADMIN_USER: {Role: ROLE_USER},
				},
				Customer: "MockStripeCustomer" + ACCOUNT_03,
				Plan:     PLAN_PREMIUM,
				Status:   STATUS_PAID,
			},
		},
	}
)
