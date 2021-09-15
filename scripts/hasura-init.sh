#!/bin/bash

echo "👉 Installing stripe price id"
curl --location --request POST 'http://localhost:8082/v1/graphql' \
    --header 'content-type: application/json' \
    --header "x-hasura-admin-secret: $HASURA_GRAPHQL_ADMIN_SECRET" \
    --data-raw "{\"query\":\"mutation UpdateBasicPlan {update_subscription_plan(where: {id: {_eq: \\\"basic\\\"}}, _set: {stripe_code: \\\"$STRIPE_PRICE_ID\\\"}) {affected_rows}}\",\"variables\":{}}"