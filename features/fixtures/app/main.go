package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	bsgperf "github.com/bugsnag/bugsnag-go-performance"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

var scenariosMap = map[string]func() (bsgperf.Configuration, func()){
	"ManualTraceScenario":          ManualTraceScenario,
	"DisabledReleaseStageScenario": DisabledReleaseStageScenario,
	"EnvironmentConfigScenario":    EnvironmentConfigScenario,
}

func configureOtel(ctx context.Context, addr string, config bsgperf.Configuration) {
	config.MainContext = ctx
	config.Endpoint = fmt.Sprintf("%v/traces", addr)
	bsgOptions, err := bsgperf.Configure(config)
	if err != nil {
		fmt.Printf("Error while creating bugsnag-go-performance: %+v\n", err)
		return
	}
	tracerProvider := trace.NewTracerProvider(bsgOptions...)
	otel.SetTracerProvider(tracerProvider)
}

func main() {
	fmt.Println("[Bugsnag] Starting testapp")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Listening to the OS Signals
	signalsChan := make(chan os.Signal, 1)
	signal.Notify(signalsChan, syscall.SIGINT, syscall.SIGTERM)
	ticker := time.NewTicker(1 * time.Second)

	addr := os.Getenv("DEFAULT_MAZE_ADDRESS")
	if addr == "" {
		addr = DEFAULT_MAZE_ADDRESS
	}

	for {
		select {
		case <-ticker.C:
			fmt.Println("[Bugsnag] Get command")
			command := GetCommand(addr)
			fmt.Printf("[Bugsnag] Received command: %+v\n", command)

			if command.Action == "run-scenario" {
				prepareScenarioFunc, ok := scenariosMap[command.ScenarioName]
				if ok {
					config, scenarioFunc := prepareScenarioFunc()
					configureOtel(ctx, addr, config)
					scenarioFunc()
				}
			}
		case <-signalsChan:
			fmt.Println("[Bugsnag] Signal received, closing")
			ticker.Stop()
			return
		}
	}
}
