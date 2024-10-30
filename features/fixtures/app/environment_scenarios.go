package main

import (
	"fmt"

	bsgperf "github.com/bugsnag/bugsnag-go-performance"
)

func EnvironmentConfigScenario() (bsgperf.Configuration, func()) {
	f := func() {
		fmt.Println("[Bugsnag] EnvironmentConfigScenario")
		createSpans("EnvironmentConfigScenario")
	}
	config := bsgperf.Configuration{
		Resource: createScenarioResource("basic app", "1.22.333", "1"),
	}

	return config, f
}
