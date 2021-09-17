#!/bin/bash

set -e

VERSION=$1

go mod tidy
GOPROXY=proxy.golang.org go list -m github.com/ddelizia/hasura-saas@$VERSION
