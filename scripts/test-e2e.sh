#!/bin/bash

./scripts/hasura-download-schema.sh

docker-compose up -d

echo "⌛️ Waiting 30s"
sleep 30

./scripts/hasura-init.sh

EXECUTE_E2E=true go test -v ./pkg/e2e