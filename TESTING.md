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

- Determine the Go version to be tested using the environment variable `GO_VERSION`, e.g. `GO_VERSION=1.22`
- Determine the OpenTelemetry SDK version using the environment variable `OTEL_VERSION`, e.g. `OTEL_VERSION=1.33.0`

The minimum versions are Go 1.22.0 and OpenTelemetry 1.33.0. The CI matrix tests
each OpenTelemetry minor release listed in these ranges, excluding incompatible
Go x OTeL combinations:

| Go version | OpenTelemetry versions |
| --- | --- |
| 1.22 | 1.33.0 - 1.35.0 |
| 1.23 | 1.33.0 - 1.38.0 |
| 1.24 | 1.33.0 - 1.41.0 |
| 1.25 | 1.33.0 - 1.46.0 |
| 1.26 | 1.33.0 - 1.47.0 |
| 1.27 | 1.33.0 - 1.47.0 |

Both unit tests and E2E fixtures use the selected OpenTelemetry version.

The `test-results` CI job succeeds only when every test matrix job succeeds.
Use this stable check name for required test checks in branch rules, rather than
individual matrix job names that change when the supported versions change.

Use the Maze Runner CLI to run the tests:

```sh
$ GO_VERSION=1.22 OTEL_VERSION=1.33.0 bundle exec maze-runner
```
