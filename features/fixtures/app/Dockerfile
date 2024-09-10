ARG GO_VERSION
ARG OTEL_VERSION
FROM golang:${GO_VERSION}-alpine

RUN apk update && apk upgrade && apk add git bash build-base

ENV GOPATH /app
ENV GO111MODULE="on"

COPY . /app/src/github.com/bugsnag/bugsnag-go-performance
WORKDIR /app/src/github.com/bugsnag/bugsnag-go-performance

# Ensure subsequent steps are re-run if the GO_VERSION variable changes
ARG GO_VERSION
ARG OTEL_VERSION

# Get bugsnag dependencies
RUN go get $(go list ./... | grep -v /features/)
WORKDIR /app/src/github.com/bugsnag/bugsnag-go-performance/features/fixtures/app

# Create app module - avoid locking bugsnag dep by not checking it in
RUN go mod init && go mod tidy; \
    go get go.opentelemetry.io/otel@v${OTEL_VERSION}; \
    go get go.opentelemetry.io/otel/sdk@v${OTEL_VERSION}; \
    echo "replace github.com/bugsnag/bugsnag-go-performance => ../../../" >> go.mod; \
    go mod tidy
