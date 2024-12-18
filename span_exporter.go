package bugsnagperformance

import (
	"context"
	"encoding/json"

	"go.opentelemetry.io/otel/sdk/trace"
)

type SpanExporter struct {
	disabled                    bool
	unmanagedMode               bool
	loggedFirstBatchDestination bool
	probabilityManager          *probabilityManager
	sampler                     *Sampler
	delivery                    *delivery
	sampleHeaderEnc             *samplingHeaderEncoder
	paylodEnc                   *payloadEncoder
}

type managedSpan struct {
	samplingProbability *float64
	span                trace.ReadOnlySpan
}

func createSpanExporter(probMgr *probabilityManager, sampler *Sampler, delivery *delivery, unmanaged bool) trace.SpanExporter {
	sp := SpanExporter{
		disabled:                    false,
		loggedFirstBatchDestination: false,
		probabilityManager:          probMgr,
		delivery:                    delivery,
		sampler:                     sampler,
		unmanagedMode:               unmanaged,
	}

	return &sp
}

func (sp *SpanExporter) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	if sp.disabled {
		return nil
	}

	sp.maybe_enter_unmanaged_mode()

	managedStatus := "managed"
	if sp.unmanagedMode {
		managedStatus = "unmanaged"
	}

	filteredSpans := []managedSpan{}
	headers := map[string]string{}

	if !sp.unmanagedMode {
		// resample spans
		for _, span := range spans {
			managedSpan, accepted := sp.sampler.resample(span)
			if accepted {
				filteredSpans = append(filteredSpans, managedSpan)
			}
		}

		samplingHeader := sp.sampleHeaderEnc.encode(filteredSpans)
		if samplingHeader == "" {
			Config.Logger.Printf("One or more spans are missing the 'bugsnag.sampling.p' attribute. This trace will be sent as unmanaged.\n")
			managedStatus = "unmanaged"
		} else {
			headers[samplingRequestHeader] = samplingHeader
		}
	} else {
		for _, span := range spans {
			filteredSpans = append(filteredSpans, managedSpan{span: span})
		}
	}

	if !sp.loggedFirstBatchDestination {
		Config.Logger.Printf("Sending %+v spans to %+v\n", managedStatus, sp.delivery.uri)
		sp.loggedFirstBatchDestination = true
	}

	// encode to JSON
	encodedPayload := sp.paylodEnc.encode(filteredSpans)
	payload, err := json.Marshal(encodedPayload)
	if err != nil {
		Config.Logger.Printf("Error encoding spans: %v\n", err)
	}

	// send payload
	resp, err := sp.delivery.send(headers, payload)
	if err != nil {
		Config.Logger.Printf("Error sending payload: %v\n", err)
	}

	// update sampling probability in ProbabilityManager
	if resp != nil {
		parsedResp := newParsedResponse(*resp)
		if parsedResp.samplingProbablity != nil {
			sp.probabilityManager.setProbability(*parsedResp.samplingProbablity)
		}
	}

	return nil
}

func (sp *SpanExporter) Shutdown(ctx context.Context) error {
	return nil
}

func (sp *SpanExporter) maybe_enter_unmanaged_mode() {
	if sp.unmanagedMode {
		return
	}

	if Config.CustomSampler != nil && sp.sampler != Config.CustomSampler {
		sp.unmanagedMode = true
	} else {
		sp.unmanagedMode = false
	}
}
