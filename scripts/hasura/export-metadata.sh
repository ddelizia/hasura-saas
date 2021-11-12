#!/bin/bash

set -e

hasura metadata export --project $HASURA_SCHEMA_FOLDER