package main

import (
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

var scenariosMap = map[string]func() (resourceData, func()){
	"ManualTraceScenario":          ManualTraceScenario,
	"DisabledReleaseStageScenario": DisabledReleaseStageScenario,
}

type resourceData struct {
	serviceName           string
	serviceVersion        string
	deviceID              string
	deploymentEnvironment string
}

func configureOtel(addr string, resourceData resourceData) {
	otelOptions := []trace.TracerProviderOption{}

	_, processors, err := bsgperf.Configure(bsgperf.Configuration{
		APIKey:               "a35a2a72bd230ac0aa0f52715bbdc6aa",
		Endpoint:             fmt.Sprintf("%v/traces", addr),
		EnabledReleaseStages: []string{"production", "staging"},
		ReleaseStage:         resourceData.deploymentEnvironment,
	})
	if err != nil {
		fmt.Printf("Error while creating bugsnag-go-performance: %+v\n", err)
		return
	}

	for _, processor := range processors {
		otelOptions = append(otelOptions, trace.WithSpanProcessor(processor))
	}

	// TODO - return resource object from Configure?
	traceRes, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(resourceData.serviceName),
			semconv.ServiceVersion(resourceData.serviceVersion),
			semconv.DeviceID(resourceData.deviceID),
			semconv.DeploymentEnvironment(resourceData.deploymentEnvironment),
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

	for {
		select {
		case <-ticker.C:
			fmt.Println("[Bugsnag] Get command")
			command := GetCommand(addr)
			fmt.Printf("[Bugsnag] Received command: %+v\n", command)

			if command.Action == "run-scenario" {
				prepareScenarioFunc, ok := scenariosMap[command.ScenarioName]
				if ok {
					resourceData, scenarioFunc := prepareScenarioFunc()
					configureOtel(addr, resourceData)
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
