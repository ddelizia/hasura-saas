#!/bin/bash

set -e

echo "👉 Installing basic stripe price id"
curl --location --request POST 'http://localhost:8082/v1/graphql' \
    --header 'content-type: application/json' \
    --header "x-hasura-admin-secret: $HASURA_GRAPHQL_ADMIN_SECRET" \
    --data-raw "{\"query\":\"mutation UpdateBasicPlan {update_subscription_plan(where: {id: {_eq: \\\"basic\\\"}}, _set: {is_active: true, stripe_code: \\\"$STRIPE_BASIC_PRICE_ID\\\"}) {affected_rows}}\",\"variables\":{}}"

echo ""
echo "👉 Installing premium stripe price id"
curl --location --request POST 'http://localhost:8082/v1/graphql' \
    --header 'content-type: application/json' \
    --header "x-hasura-admin-secret: $HASURA_GRAPHQL_ADMIN_SECRET" \
    --data-raw "{\"query\":\"mutation UpdatePremiumPlan {update_subscription_plan(where: {id: {_eq: \\\"premium\\\"}}, _set: {is_active: true, stripe_code: \\\"$STRIPE_PREMIUM_PRICE_ID\\\"}) {affected_rows}}\",\"variables\":{}}"