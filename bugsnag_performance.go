package bugsnagperformance

import (
	"context"
	"fmt"
	"os"
	"sync"

	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// Version defines the version of this Bugsnag performance module
const Version = "0.1.0"

// Config is the configuration for the default Bugsnag performance module
var Config Configuration

var readEnvConfigOnce sync.Once

func init() {
	fmt.Println("Starting bugsnag performance")
	Config = Configuration{
		ReleaseStage: "production",
	}
}

type BugsnagPerformance struct {
	Sampler    trace.Sampler
	Processors []trace.SpanProcessor
	Resource   *resource.Resource
}

// Configure Bugsnag. The only required setting is the APIKey, which can be
// obtained by clicking on "Settings" in your Bugsnag dashboard.
// Returns OTeL sampler, probability attribute processor, trace exporter and error
func Configure(config Configuration) (*BugsnagPerformance, error) {
	readEnvConfigOnce.Do(Config.loadEnv)
	Config.update(&config)
	err := Config.validate()
	if err != nil {
		return nil, err
	}

	delivery := createDelivery()
	ctx := context.Background()
	if Config.MainContext != nil {
		ctx = Config.MainContext
	}
	probabilityManager := createProbabilityManager(ctx, delivery)
	sampler := createSampler(probabilityManager)
	probAttrProcessor := createProbabilityAttributeProcessor(probabilityManager)
	processors := []trace.SpanProcessor{probAttrProcessor}

	// enter unmanaged mode if the OTel sampler environment variable has been set
	// note: we assume any value means a non-default sampler will be used because
	//       we don't control what the valid values are
	unmanagedMode := false
	if customSampler := os.Getenv("OTEL_TRACES_SAMPLER"); customSampler != "" {
		fmt.Printf("UNMANAGED MODE ENABLED: %v\n", customSampler)
		unmanagedMode = true
	}

	// Create an exporter only if the configured release stage is enabled
	if Config.isReleaseStageEnabled() {
		fmt.Printf("RELEASE STAGE IS ENABLED: %v\n", Config.ReleaseStage)
		spanExporter := createSpanExporter(probabilityManager, sampler, delivery, unmanagedMode)
		// Batch processor with default settings
		bsgSpanProcessor := trace.NewBatchSpanProcessor(spanExporter)
		processors = append(processors, bsgSpanProcessor)
	}

	bsgResource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.DeploymentEnvironment(Config.ReleaseStage),
	)

	performanceItems := &BugsnagPerformance{
		Sampler:    sampler,
		Processors: processors,
		Resource:   bsgResource,
	}

	return performanceItems, nil
}
