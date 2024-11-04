package main

import (
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
)

func createScenarioResource(srvName, srvVer, deviceID string) *resource.Resource {
	attr := []attribute.KeyValue{
		{
			Key:   attribute.Key("device.id"),
			Value: attribute.StringValue(deviceID),
		},
		{
			Key:   attribute.Key("service.name"),
			Value: attribute.StringValue(srvName),
		},
		{
			Key:   attribute.Key("service.version"),
			Value: attribute.StringValue(srvVer),
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
