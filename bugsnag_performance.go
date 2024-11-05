package bugsnagperformance

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

// Version defines the version of this Bugsnag performance module
const Version = "0.1.0"

// Config is the configuration for the default Bugsnag performance module
var Config Configuration

var readEnvConfigOnce sync.Once

func init() {
	fmt.Println("Starting bugsnag performance.")
	Config = Configuration{
		ReleaseStage: "production",
		Logger:       log.New(os.Stdout, "[BugsnagPerformance] ", log.LstdFlags),
		MainContext:  context.TODO(),
	}
}

// Configure Bugsnag. The only required setting is the APIKey, which can be
// obtained by clicking on "Settings" in your Bugsnag dashboard.
// Returns OTeL sampler, probability attribute processor, trace exporter and error
func Configure(config Configuration) ([]trace.TracerProviderOption, error) {
	readEnvConfigOnce.Do(Config.loadEnv)
	Config.update(&config)
	err := Config.validate()
	if err != nil {
		return nil, err
	}

	otelOptions := createBugsnagOtelOptions()

	return otelOptions, nil
}

func createBugsnagOtelOptions() []trace.TracerProviderOption {
	delivery := createDelivery()
	probabilityManager := createProbabilityManager(Config.MainContext, delivery)
	sampler := createSampler(probabilityManager)

	otelOptions := []trace.TracerProviderOption{}
	probAttrProcessor := createProbabilityAttributeProcessor(probabilityManager)
	otelOptions = append(otelOptions, trace.WithSpanProcessor(probAttrProcessor))

	// enter unmanaged mode if the OTel sampler environment variable has been set
	// note: we assume any value means a non-default sampler will be used because
	//       we don't control what the valid values are
	unmanagedMode := false
	if customSampler := os.Getenv("OTEL_TRACES_SAMPLER"); customSampler != "" || Config.CustomSampler != nil {
		unmanagedMode = true
		otelOptions = append(otelOptions, trace.WithSampler(Config.CustomSampler))
	} else {
		otelOptions = append(otelOptions, trace.WithSampler(sampler))
	}

	if Config.isReleaseStageEnabled() {
		spanExporter := createSpanExporter(probabilityManager, sampler, delivery, unmanagedMode)
		otelOptions = append(otelOptions, trace.WithSpanProcessor(trace.NewBatchSpanProcessor(spanExporter)))
	}

	otelOptions = append(otelOptions, trace.WithResource(createBugsnagMergedResource()))

	return otelOptions
}

func createBugsnagMergedResource() *resource.Resource {
	customResource := Config.Resource
	if customResource == nil {
		customResource = resource.Default()
	}

	attr := []attribute.KeyValue{
		{
			Key:   "deployment.environment",
			Value: attribute.StringValue(Config.ReleaseStage),
		},
		{
			Key:   "service.version",
			Value: attribute.StringValue(Config.AppVersion),
		}}
	bsgResource, err := resource.Merge(
		customResource,
		resource.NewSchemaless(attr...),
	)
	if err != nil {
		Config.Logger.Printf("Error while merging resource: %+v\n", err)
		return customResource
	}

	return bsgResource
}
