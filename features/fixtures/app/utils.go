package main

import (
	"fmt"

	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func createScenarioResource(srvName, srvVer, deviceID string) *resource.Resource {
	traceRes, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(srvName),
			semconv.ServiceVersion(srvVer),
			semconv.DeviceID(deviceID)),
	)
	if err != nil {
		fmt.Printf("Error while merging resource: %+v\n", err)
		return nil
	}
	return traceRes
}
