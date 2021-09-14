FROM golang:1.16-alpine AS base

RUN apk add build-base
WORKDIR /app-build
COPY . .


#Build 
FROM base as build
ARG APP
RUN go get -t ./cmd/${APP}
RUN go build -o /dist/app ./cmd/${APP}


# Release
FROM alpine AS release
WORKDIR /dist
COPY --from=build /dist/ /dist/
ENTRYPOINT ./app
