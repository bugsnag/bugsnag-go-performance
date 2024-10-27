package main

import (
	"fmt"

	bsgperf "github.com/bugsnag/bugsnag-go-performance"
)

func EnvironmentConfigScenario() (resourceData, bsgperf.Configuration, func()) {
	f := func() {
		fmt.Println("[Bugsnag] EnvironmentConfigScenario")
		createSpans("EnvironmentConfigScenario")
	}
	resource := resourceData{
		serviceName:    "basic app",
		serviceVersion: "1.2.3",
		deviceID:       "1",
	}
	return resource, bsgperf.Configuration{}, f
}
