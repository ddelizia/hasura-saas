#!/bin/bash

set -e

source ./scripts/hasura/functions.sh

DATABASE="saas"
WAIT=10

echo "⌛️ Wait $WAIT s"
sleep $WAIT

echo "🧪 Initializing the database"
hasura_init $DATABASE

echo "🧪 Executing all down migrations"
hasura_down_migrations $DATABASE

echo "🧪 Executing all up migrations"
hasura_up_migrations $DATABASE