package bugsnagperformance

import (
	"fmt"
	"sync"

	"go.opentelemetry.io/otel/sdk/trace"
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

// Configure Bugsnag. The only required setting is the APIKey, which can be
// obtained by clicking on "Settings" in your Bugsnag dashboard.
// Returns OTeL sampler, probability attribute processor, trace exporter and error
func Configure(config Configuration) (trace.Sampler, []trace.SpanProcessor, error) {
	readEnvConfigOnce.Do(Config.loadEnv)
	Config.update(&config)
	err := Config.validate()
	if err != nil {
		return nil, nil, err
	}

	probabilityManager := createProbabilityManager()
	sampler := createSampler(probabilityManager)
	spanExporter := createSpanExporter(probabilityManager, sampler)
	probAttrProcessor := createProbabilityAttributeProcessor(probabilityManager)
	// Batch processor with default settings
	bsgSpanProcessor := trace.NewBatchSpanProcessor(spanExporter)
	processors := []trace.SpanProcessor{probAttrProcessor, bsgSpanProcessor}

	return sampler, processors, nil
}
