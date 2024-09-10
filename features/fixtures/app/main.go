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
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

var scenariosMap = map[string]func() (string, func()){
	"HandledScenario": HandledScenario,
}

func HandledScenario() (string, func()) {
	f := func() {
		fmt.Println("HELLO WORLD")
		_, span := otel.GetTracerProvider().Tracer("maze-test").Start(context.Background(), "HandledScenario")

		// TODO - hardcoded sampling attribute
		span.SetAttributes(attribute.KeyValue{
			Key: "bugsnag.sampling.p",
			Value: attribute.Float64Value(1.0),
		})
		span.End()
	}
	return "OUTPUT", f
}

func configureOtel(addr string) {
	otelOptions := []trace.TracerProviderOption{}

	_, bsgExporter, err := bsgperf.Configure(bsgperf.Configuration{
		APIKey:   "a35a2a72bd230ac0aa0f52715bbdc6aa",
		Endpoint: fmt.Sprintf("%v/traces", addr),
	})
	if err != nil {
		fmt.Printf("Error while creating bugsnag-go-performance: %+v\n", err)
		return
	}

	if bsgExporter != nil {
		otelOptions = append(otelOptions, trace.WithSpanProcessor(bsgExporter))
	}

	traceRes, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("basic app"),
		semconv.ServiceVersion("1.22.333"),
		semconv.DeviceID("1"),
		semconv.DeploymentEnvironment("production"),
	))
	if err != nil {
		fmt.Printf("Error while creating resource: %+v\n", err)
	}
	otelOptions = append(otelOptions, trace.WithResource(traceRes))

	tracerProvider := trace.NewTracerProvider(otelOptions...)
	otel.SetTracerProvider(tracerProvider)
}

func main() {
	fmt.Println("[Bugsnag] Starting testapp")
	// Listening to the OS Signals
	signalsChan := make(chan os.Signal, 1)
	signal.Notify(signalsChan, syscall.SIGINT, syscall.SIGTERM)
	ticker := time.NewTicker(1 * time.Second)

	addr := os.Getenv("DEFAULT_MAZE_ADDRESS")
	if addr == "" {
		addr = DEFAULT_MAZE_ADDRESS
	}

	configureOtel(addr)

	for {
		select {
		case <-ticker.C:
			fmt.Println("[Bugsnag] Get command")
			command := GetCommand(addr)
			fmt.Printf("[Bugsnag] Received command: %+v\n", command)

			if command.Action == "run-scenario" {
				prepareScenarioFunc, ok := scenariosMap[command.ScenarioName]
				if ok {
					_, scenarioFunc := prepareScenarioFunc()
					//bugsnag.Configure(config)
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
