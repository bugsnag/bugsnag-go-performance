package main

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	bsgperf "github.com/bugsnag/bugsnag-go-performance"
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

func ManualTraceScenario() (resourceData, bsgperf.Configuration, func()) {
	f := func() {
		fmt.Println("[Bugsnag] ManualTraceScenario")
		createSpans("ManualTraceScenario")
	}
	config := bsgperf.Configuration{
		APIKey:               "a35a2a72bd230ac0aa0f52715bbdc6aa",
		EnabledReleaseStages: []string{"production", "staging"},
		ReleaseStage:         "staging",
	}

	resource := resourceData{
		serviceName:    "basic app",
		serviceVersion: "1.22.333",
		deviceID:       "1",
	}
	return resource, config, f
}

func DisabledReleaseStageScenario() (resourceData, bsgperf.Configuration, func()) {
	f := func() {
		fmt.Println("[Bugsnag] ManualTraceScenario")
		createSpans("DisabledReleaseStageScenario")
	}

	config := bsgperf.Configuration{
		APIKey:               "a35a2a72bd230ac0aa0f52715bbdc6aa",
		EnabledReleaseStages: []string{"production", "staging"},
		ReleaseStage:         "development",
	}

	resource := resourceData{
		serviceName:    "basic app",
		serviceVersion: "1.22.333",
		deviceID:       "1",
	}
	return resource, config, f
}
