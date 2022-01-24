package model

import "github.com/ddelizia/hasura-saas/pkg/hstype"

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
	IDPaymentMethod hstype.String
}

type CreateOutput struct {
	IDAccount string
	IsActive  bool
}

type RetryInput struct {
	IDAccount       string
	IDPaymentMethod hstype.String
}

type RetryOutput struct {
	IDAccount string
	IsActive  bool
}
