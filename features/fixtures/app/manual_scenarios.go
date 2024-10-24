package main

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func createSpans(scenarioName string) {
	for i := 0; i < 5; i++ {
		_, span := otel.GetTracerProvider().Tracer("maze-test").Start(context.Background(), scenarioName)
		span.SetName(fmt.Sprintf("test span %v", i+1))
		span.SetAttributes(attribute.KeyValue{
			Key:   attribute.Key("span.custom.age"),
			Value: attribute.IntValue(i * 10),
		})
		span.End()
	}
}

func ManualTraceScenario() (resourceData, func()) {
	f := func() {
		fmt.Println("[Bugsnag] ManualTraceScenario")
		createSpans("ManualTraceScenario")
	}
	resource := resourceData{
		serviceName:           "basic app",
		serviceVersion:        "1.22.333",
		deviceID:              "1",
		deploymentEnvironment: "staging",
	}
	return resource, f
}

func DisabledReleaseStageScenario() (resourceData, func()) {
	f := func() {
		fmt.Println("[Bugsnag] ManualTraceScenario")
		createSpans("DisabledReleaseStageScenario")
	}
	resource := resourceData{
		serviceName:           "basic app",
		serviceVersion:        "1.22.333",
		deviceID:              "1",
		deploymentEnvironment: "development",
	}
	return resource, f
}
