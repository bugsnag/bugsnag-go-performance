Feature: bugsnag errors integration

#Scenario: It picks up configuration from bugsnag errors configuration

Scenario: It picks up configuration from bugsnag errors' environment variables
  Given I set environment variable "BUGSNAG_API_KEY" to "ab123456789012345678901234567890"
  Given I set environment variable "BUGSNAG_RELEASE_STAGE" to "developroduction"
  Given I set environment variable "BUGSNAG_APP_VERSION" to "1.2.3"
  When I start the service "app"
  Then I run "EnvironmentConfigScenario"
  And I wait to receive a trace
  And the trace "Bugsnag-Api-Key" header equals "ab123456789012345678901234567890"
  And the trace payload field "resourceSpans.0.resource" string attribute "deployment.environment" equals "developroduction"
  And the trace payload field "resourceSpans.0.resource" string attribute "service.version" equals "1.2.3"


Scenario: It picks up configuration from performance environment variables over errors'
  Given I set environment variable "BUGSNAG_API_KEY" to "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
  Given I set environment variable "BUGSNAG_RELEASE_STAGE" to "development"
  Given I set environment variable "BUGSNAG_APP_VERSION" to "1.2.3"
  Given I set environment variable "BUGSNAG_PERFORMANCE_API_KEY" to "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
  Given I set environment variable "BUGSNAG_PERFORMANCE_RELEASE_STAGE" to "production"
  Given I set environment variable "BUGSNAG_PERFORMANCE_SERVICE_NAME" to "srv2"
  Given I set environment variable "BUGSNAG_PERFORMANCE_APP_VERSION" to "3.2.1"
  When I start the service "app"
  Then I run "EnvironmentConfigScenario"
  And I wait to receive a trace
  And the trace "Bugsnag-Api-Key" header equals "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
  And the trace payload field "resourceSpans.0.resource" string attribute "deployment.environment" equals "production"
  And the trace payload field "resourceSpans.0.resource" string attribute "service.name" equals "srv2"
  And the trace payload field "resourceSpans.0.resource" string attribute "service.version" equals "3.2.1"
