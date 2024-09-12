package bugsnagperformance

import (
	"testing"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func TestDefaultHeaderValue(t *testing.T) {
	enc := &samplingHeaderEncoder{}

	result := enc.encode([]sdktrace.ReadOnlySpan{})
	if result != "1.0:0" {
		t.Errorf("Expected '1.0:0', got %s", result)
	}
}

// TODO check how to create spans and readonly spans