# Hasura SaaS

[![GitHub stars](https://img.shields.io/github/stars/ddelizia/hasura-saas.svg?style=social&label=Star&maxAge=2592000)](https://GitHub.com/ddelizia/hasura-saas/stargazers/) 
[![GitHub forks](https://img.shields.io/github/forks/ddelizia/hasura-saas.svg?style=social&label=Fork&maxAge=2592000)](https://GitHub.com/ddelizia/hasura-saas/network/) 

[![PkgGoDev](https://pkg.go.dev/badge/github.com/ddelizia/hasura-saas)](https://pkg.go.dev/github.com/ddelizia/hasura-saas)
[![Test Actions Status](https://github.com/ddelizia/hasura-saas/workflows/ci/badge.svg)](https://github.com/ddelizia/hasura-saas/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/ddelizia/hasura-saas)](https://goreportcard.com/report/github.com/ddelizia/hasura-saas)
[![MIT license](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/ddelizia/channelify/blob/main/LICENSE)

## About the project

This repository is an intent to create a set of microservices to be able to spin up a SaaS application backed by Hasura.

### Status

Work in progress

## Architecture

![Architecture Overview](docs/images/Architecture.png)

## Configuration

### Authentication and Authorization

Authentication and Authorization are the most important aspect of an SaaS. Hasura does not provide any out of the box authentication mechanism but it provides the fine grained authorization capabilities. Anyway hasura is able to read jwt tokens and it is able to inject claims has part of the hasura request.

This project relies on 3rd party authentication like Auth0 (for the moment it has been tested with Auth0 and maybe in the future will support additional providers mechanisms besides jwt), anyway it should be able to work with other openid connect providers such as Cognito or Firebase Auth.

#### Auth0

In oreder to configure Auth0, you need to follow the steps provided in this page https://hasura.io/docs/latest/graphql/core/guides/integrations/auth0-jwt.html. With some small changes that will be listed below.

When your are at he step `Configure Auth0 Rules & Callback URLs` you will need to use the following snippet:

```javascript
function (user, context, callback) {
  const namespace = "https://hasura.io/jwt/claims";
  context.idToken[namespace] =
    {
      'x-hasura-default-role': 'user',
      // do some custom logic to decide allowed roles
      'x-hasura-allowed-roles': ['user'],
      'x-hasura-user-id': user.id_user
    };
  callback(null, user, context);
}
```


### Subscription and Payment

