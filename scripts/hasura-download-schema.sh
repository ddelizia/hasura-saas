#!/bin/bash

VERSION="v0.0.1-alpha.1"

download () {
  FILE=$1

  echo "👉 Downloading artifact $FILE"
  (mkdir -p hasura/$FILE && cd ./hasura/$FILE && wget https://github.com/ddelizia/hasura-saas-starter/releases/download/$VERSION/$FILE.zip && unzip $FILE.zip && rm $FILE.zip)
} 

rm -Rf ./hasura

download "migrations"
download "metadata"