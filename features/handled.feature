Feature: Dummy mazerunner tests

Scenario: Run first scenario
  When I start the service "app"
  Then I run "HandledScenario"
  And I wait to receive a trace
  And the trace payload field "resourceSpans.0.resource" string attribute "service.name" equals "basic app"
  And the trace payload field "resourceSpans.0.resource" string attribute "service.version" equals "1.22.333"
  And the trace payload field "resourceSpans.0.resource" string attribute "device.id" equals "1"
  And the trace payload field "resourceSpans.0.resource" string attribute "deployment.environment" equals "production"

