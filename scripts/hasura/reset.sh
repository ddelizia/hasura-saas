#!/bin/bash

set -e

source ./scripts/hasura/functions.sh

DATABASE=saas

echo "👉 Initializing the database"
hasura_init $DATABASE