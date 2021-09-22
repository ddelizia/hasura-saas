#!/bin/bash

set -e

TAG=$1

BUILD_TAG=${TAG:-latest}

build_docker () {
  APP=$1
  IMAGE_NAME=ghcr.io/ddelizia/hasura-saas-$APP
  echo "Building $IMAGE_NAME:$BUILD_TAG"
	docker build --build-arg APP="$APP" --target release -t $IMAGE_NAME:$BUILD_TAG .
}

build_docker rp
build_docker subscription
build_docker saas
