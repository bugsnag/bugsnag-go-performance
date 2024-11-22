Feature: Manual spans

Scenario: It runs the basic app
  When I start the service "app"
  Then I run "ManualTraceScenario"
  When I wait to receive a trace
  Then the sampling request "Bugsnag-Span-Sampling" header equals "1.0:0"
  And the trace "Bugsnag-Span-Sampling" header equals "1.0:5"

  And the trace payload field "resourceSpans.0.resource" string attribute "service.name" equals "basic app"
  And the trace payload field "resourceSpans.0.resource" string attribute "service.version" equals "1.22.333"
  And the trace payload field "resourceSpans.0.resource" string attribute "device.id" equals "1"
  And the trace payload field "resourceSpans.0.resource" string attribute "deployment.environment" equals "staging"

  And a span named "test span 1" has the following properties:
    | property               | value       |
    | kind                   | 1           |
    | traceState             |             |
    | droppedAttributesCount | 0           |
    | droppedEventsCount     | 0           |
    | droppedLinksCount      | 0           |
    | status.code            | 0           |
    | status.message         |             |
  And a span named "test span 1" contains the attributes:
    | attribute                | type        | value |
    | span.custom.age          | intValue    | 0     |
    | bugsnag.sampling.p       | doubleValue | 1     |
    | bugsnag.span.first_class | boolValue   | true  |

  And a span named "test span 2" has the following properties:
    | property               | value       |
    | kind                   | 1           |
    | traceState             |             |
    | droppedAttributesCount | 0           |
    | droppedEventsCount     | 0           |
    | droppedLinksCount      | 0           |
    | status.code            | 0           |
    | status.message         |             |
  And a span named "test span 2" contains the attributes:
    | attribute                | type        | value |
    | span.custom.age          | intValue    | 10    |
    | bugsnag.sampling.p       | doubleValue | 1     |
    | bugsnag.span.first_class | boolValue   | true  |

  And a span named "test span 3" has the following properties:
    | property               | value       |
    | kind                   | 1           |
    | traceState             |             |
    | droppedAttributesCount | 0           |
    | droppedEventsCount     | 0           |
    | droppedLinksCount      | 0           |
    | status.code            | 0           |
    | status.message         |             |
  And a span named "test span 3" contains the attributes:
    | attribute                | type        | value |
    | span.custom.age          | intValue    | 20    |
    | bugsnag.sampling.p       | doubleValue | 1     |
    | bugsnag.span.first_class | boolValue   | true  |

  And a span named "test span 4" has the following properties:
    | property               | value       |
    | kind                   | 1           |
    | traceState             |             |
    | droppedAttributesCount | 0           |
    | droppedEventsCount     | 0           |
    | droppedLinksCount      | 0           |
    | status.code            | 0           |
    | status.message         |             |
  And a span named "test span 4" contains the attributes:
    | attribute                | type        | value |
    | span.custom.age          | intValue    | 30    |
    | bugsnag.sampling.p       | doubleValue | 1     |
    | bugsnag.span.first_class | boolValue   | true  |

  And a span named "test span 5" has the following properties:
    | property               | value       |
    | kind                   | 1           |
    | traceState             |             |
    | droppedAttributesCount | 0           |
    | droppedEventsCount     | 0           |
    | droppedLinksCount      | 0           |
    | status.code            | 0           |
    | status.message         |             |
  And a span named "test span 5" contains the attributes:
    | attribute                | type        | value |
    | span.custom.age          | intValue    | 40    |
    | bugsnag.sampling.p       | doubleValue | 1     |
    | bugsnag.span.first_class | boolValue   | true  |

Scenario: It picks up default service name from OTEL
  When I start the service "app"
  Then I run "ServiceNameScenario"
  And I wait to receive a trace
  And the trace payload field "resourceSpans.0.resource" string attribute "service.name" equals "unknown_service:app"

Scenario: It picks up service name from OTEL env var
  Given I set environment variable "OTEL_SERVICE_NAME" to "mysrvname"
  When I start the service "app"
  Then I run "ServiceNameScenario"
  And I wait to receive a trace
  And the trace payload field "resourceSpans.0.resource" string attribute "service.name" equals "mysrvname"

Scenario: It does not export spans when the release stage is disabled
  When I start the service "app"
  Then I run "DisabledReleaseStageScenario"
  Then I should receive no traces