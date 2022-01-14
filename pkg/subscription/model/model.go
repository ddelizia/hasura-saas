package model

type InitInput struct {
	AccountName string
	IDPlan      string
	IDUser      string
}

type InitOutput struct {
	IDAccount string
}
