package model

type InitInput struct {
	AccountName string
	IDPlan      string
	IDUser      string
}

type InitOutput struct {
	IDAccount string
}

type CreateInput struct {
	IDAccount       string
	IDPaymentMethod string
}

type CreateOutput struct {
	IDAccount string
	IsActive  bool
}
