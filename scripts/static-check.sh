#!/bin/bash

go install honnef.co/go/tools/cmd/staticcheck@latest
$(go env GOPATH)/bin/staticcheck ./...