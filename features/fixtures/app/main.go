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
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

var scenariosMap = map[string]func() (resourceData, bsgperf.Configuration, func()){
	"ManualTraceScenario":          ManualTraceScenario,
	"DisabledReleaseStageScenario": DisabledReleaseStageScenario,
	"EnvironmentConfigScenario":    EnvironmentConfigScenario,
}

type resourceData struct {
	serviceName    string
	serviceVersion string
	deviceID       string
}

func configureOtel(ctx context.Context, addr string, resData resourceData, config bsgperf.Configuration) {
	otelOptions := []trace.TracerProviderOption{}

	config.MainContext = ctx
	config.Endpoint = fmt.Sprintf("%v/traces", addr)
	bsgPerformance, err := bsgperf.Configure(config)
	if err != nil {
		fmt.Printf("Error while creating bugsnag-go-performance: %+v\n", err)
		return
	}

	for _, processor := range bsgPerformance.Processors {
		otelOptions = append(otelOptions, trace.WithSpanProcessor(processor))
	}

	if bsgPerformance.Sampler != nil {
		otelOptions = append(otelOptions, trace.WithSampler(bsgPerformance.Sampler))
	}

	// normal creation
	traceRes, err := resource.Merge(
		resource.Default(),
		bsgPerformance.Resource,
	)
	if err != nil {
		fmt.Printf("Error while merging resource: %+v\n", err)
	}

	// setup data from tests to be merged with the resource
	traceRes, err = resource.Merge(
		traceRes,
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(resData.serviceName),
			semconv.ServiceVersion(resData.serviceVersion),
			semconv.DeviceID(resData.deviceID)),
	)
	if err != nil {
		fmt.Printf("Error while merging resource: %+v\n", err)
	}

	otelOptions = append(otelOptions, trace.WithResource(traceRes))

	tracerProvider := trace.NewTracerProvider(otelOptions...)
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
					resData, config, scenarioFunc := prepareScenarioFunc()
					configureOtel(ctx, addr, resData, config)
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
