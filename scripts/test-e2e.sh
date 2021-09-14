#!/bin/bash

./scripts/hasura-download-schema.sh

docker-compose up -d

./scripts/hasura-init.sh

env GRAPHQL.HASURA.ADMINSECRET=mysecrethasura EXECUTE_E2E=true go test -v ./pkg/e2e