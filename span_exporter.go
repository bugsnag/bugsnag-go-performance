package bugsnagperformance

import (
	"context"
	"encoding/json"
	"fmt"

	"go.opentelemetry.io/otel/sdk/trace"
)

type SpanExporter struct {
	disabled       bool
	unmanagedMode bool
	loggedFirstBatchDestination bool
	probabilityManager interface{}
	delivery *delivery
	samplingHeaderEncoder
	payloadEncoder
}

func CreateSpanExporter() trace.SpanExporter {
	delivery := createDelivery(Config.Endpoint, Config.APIKey)

	sp := SpanExporter{
		disabled: false,
		loggedFirstBatchDestination: false,
		probabilityManager: nil,
		delivery: delivery,
	}

	return &sp
}

func (sp *SpanExporter) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	if sp.disabled {
		return nil
	}
	
	 managedStatus := "managed"
	 if sp.unmanagedMode {
		managedStatus = "unmanaged"
	 }

	headers := map[string]string{}
	if !sp.unmanagedMode {

		samplingHeader := sp.samplingHeaderEncoder.encode(spans)

		if samplingHeader == "" {
			fmt.Println("One or more spans are missing the 'bugsnag.sampling.p' attribute. This trace will be sent as unmanaged")
			managedStatus = "unmanaged"
		} else {
			headers["Bugsnag-Span-Sampling"] = samplingHeader
		}
	}

	if !sp.loggedFirstBatchDestination {
		fmt.Printf("Sending %+v spans to %+v\n", managedStatus, "url")
		sp.loggedFirstBatchDestination = true
	}

	// encode to JSON
	encodedPayload := sp.payloadEncoder.encode(spans)
	payload, err := json.Marshal(encodedPayload)
	if err != nil {
		fmt.Printf("Error encoding spans: %v\n", err)
	}

	// send payload
	sp.delivery.send(headers, payload)

	// update sampling probability in ProbabilityManager

	return nil
}

func (sp *SpanExporter) Shutdown(ctx context.Context) error {
	return nil
}
