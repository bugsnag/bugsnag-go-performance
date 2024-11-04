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
		AppVersion: "1.22.333",
		Resource:   createScenarioResource("basic app", "1"),
	}

	return config, f
}
