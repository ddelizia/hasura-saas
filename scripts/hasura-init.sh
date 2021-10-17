#!/bin/bash

set -e

echo "👉 Installing free stripe price id"
curl --location --request POST 'http://localhost:8082/v1/graphql' \
    --header 'content-type: application/json' \
    --header "x-hasura-admin-secret: $HASURA_GRAPHQL_ADMIN_SECRET" \
    --data-raw "{\"query\":\"mutation UpdateFreePlan {update_subscription_plan(where: {id: {_eq: \\\"free\\\"}}, _set: {is_active: true, stripe_code: \\\"$STRIPE_FREE_PRICE_ID\\\"}) {affected_rows}}\",\"variables\":{}}"

echo ""
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

echo ""
echo "👉 Create basic plan with trial"
curl --location --request POST 'http://localhost:8082/v1/graphql' \
    --header 'content-type: application/json' \
    --header "x-hasura-admin-secret: $HASURA_GRAPHQL_ADMIN_SECRET" \
    --data-raw "{\"query\":\"mutation UpdateOrCreateFreeBasic { insert_subscription_plan(objects: {description: \\\"Basic with trial\\\", id: \\\"basic_trial\\\", is_active: true, price: \\\"0\\\", stripe_code: \\\"$STRIPE_BASIC_TRIAL_PRICE_ID\\\", trial_days: 30}, on_conflict: {constraint: subscription_plan_pkey, update_columns: [is_active, description, price, stripe_code, trial_days]}) {returning {id}}}\",\"variables\":{}}"