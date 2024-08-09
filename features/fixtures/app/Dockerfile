ARG GO_VERSION
FROM golang:${GO_VERSION}-alpine

RUN apk update && apk upgrade && apk add git bash build-base

ENV GOPATH /app

COPY . /app/src/github.com/bugsnag/bugsnag-go-performance
WORKDIR /app/src/github.com/bugsnag/bugsnag-go-performance

# Ensure subsequent steps are re-run if the GO_VERSION variable changes
ARG GO_VERSION

# Get bugsnag dependencies using a conditional call to run go get or go install based on the go version
RUN if [[ $(echo -e "1.11\n$GO_VERSION\n1.16" | sort -V | head -2 | tail -1) == "$GO_VERSION" ]]; then \
        echo "Version is between 1.11 and 1.16, running go get"; \
        go get $(go list ./... | grep -v /features/); \
    else \
        echo "Version is greater than 1.16, running go install"; \
        go install $(go list ./... | grep -v /features/); \
    fi

WORKDIR /app/src/github.com/bugsnag/bugsnag-go-performance/features/fixtures/app

# Create app module - avoid locking bugsnag dep by not checking it in
# Skip on old versions of Go which pre-date modules
RUN if [[ $GO_VERSION != '1.11' && $GO_VERSION != '1.12' ]]; then \
        go mod init && go mod tidy; \
        #echo "replace github.com/bugsnag/bugsnag-go-performance => ../../../" >> go.mod; \
        go mod tidy; \
    fi
