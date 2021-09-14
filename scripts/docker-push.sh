#!/bin/bash

S_TAG=$1
T_TAG=$2

SOURCE_TAG=${S_TAG:-latest}
TARGET_TAG=${T_TAG:-latest}

push_docker () {
  APP=$1
  IMAGE_NAME=ghcr.io/ddelizia/hasura-saas-$APP
  echo "Pushing $APP"
  docker tag $IMAGE_NAME:$SOURCE_TAG $IMAGE_NAME:$TARGET_TAG
  docker push $IMAGE_NAME:$TARGET_TAG
}

build_docker rp
build_docker subscription
