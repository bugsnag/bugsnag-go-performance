package main

import (
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
)

func createScenarioResource(deviceID string) *resource.Resource {
	attr := []attribute.KeyValue{
		{
			Key:   attribute.Key("device.id"),
			Value: attribute.StringValue(deviceID),
		},
	}
	traceRes, err := resource.Merge(
		resource.Default(),
		resource.NewSchemaless(attr...),
	)
	if err != nil {
		fmt.Printf("Error while merging resource: %+v\n", err)
		return nil
	}
	return traceRes
}
