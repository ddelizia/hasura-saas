#!/bin/bash

set -e

hasura_down_migrations () {
  DATABASE=$1

  echo "👉 Execute all down migrations..."
  hasura migrate apply --database-name $DATABASE --disable-interactive --project $HASURA_SCHEMA_FOLDER --down all 
  
  echo "👉 Checking migration status..."
  hasura migrate status --database-name $DATABASE --disable-interactive --project $HASURA_SCHEMA_FOLDER 
}

hasura_up_migrations () {
  DATABASE=$1

  echo "👉 Execute all up migrations..."
  hasura migrate apply --database-name $DATABASE --disable-interactive --project $HASURA_SCHEMA_FOLDER --up all 
  
  echo "👉 Checking migration status..."
  hasura migrate status --database-name $DATABASE --disable-interactive --project $HASURA_SCHEMA_FOLDER 
} 

hasura_restart_migrations () {
  DATABASE=$1

  hasura_down_migrations $DATABASE

  hasura_up_migrations $DATABASE
}

hasura_init () {
  DATABASE=$1

  echo "👉 Applying metadata..."
  hasura metadata apply --project $HASURA_SCHEMA_FOLDER

  echo "👉 Checking metadata diffs..."
  hasura metadata diff --project $HASURA_SCHEMA_FOLDER

  hasura_down_migrations $DATABASE

  hasura_up_migrations $DATABASE

  echo "👉 Reload metadata..."
  hasura metadata reload --project $HASURA_SCHEMA_FOLDER
}

