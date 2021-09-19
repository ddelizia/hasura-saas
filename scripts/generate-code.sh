#!/bin/bash

set -e

echo "Generating graphql classes"
$(go env GOPATH)/bin/gqlgenc

echo "Generating stripe webhook"
go run ./cmd/stripe_event_code_generator/main.go