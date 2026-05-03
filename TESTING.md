# Testing the Go BugSnag performance SDK

## Unit tests

```
go test ./...
```

## End-to-end tests

These tests are implemented with our internal testing tool [Maze Runner](https://github.com/bugsnag/maze-runner).

End to end tests are written in cucumber-style `.feature` files, and need Ruby-backed "steps" in order to know what to run. The tests are located in the top level [`features`](./features/) directory.

The Maze Runner test fixtures are containerised so you'll need Docker and Docker Compose to run them.

### Running the end to end tests

Install Maze Runner:

```sh
$ bundle install
```

Configure the tests to be run in the following way:

- Determine the Go version to be tested using the environment variable `GO_VERSION`, e.g. `GO_VERSION=1.25`
- Determine the Open Telemetry SDK version using the environment variable `OTEL_VERSION`, e.g. `OTEL_VERSION=1.43.0`

The current module directly requires OTel v1.43.0, so local and CI runs should use:
* Go 1.25
* OTel 1.43.0

Use the Maze Runner CLI to run the tests:

```sh
$ GO_VERSION=1.25 OTEL_VERSION=1.43.0 bundle exec maze-runner
```