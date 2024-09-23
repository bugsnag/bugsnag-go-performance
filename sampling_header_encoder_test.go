package bugsnagperformance

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func makeSpan(tp *trace.TracerProvider, name string, probability ...float64) {
	_, span := tp.Tracer("test").Start(context.Background(), name)

	if len(probability) > 0 {
		span.SetAttributes(attribute.Float64(BUGSNAG_SAMPLING_ATTRIBUTE, probability[0]))
	}

	span.End()
}

func getSpans(exporter *tracetest.InMemoryExporter) []managedSpan {
	spans := exporter.GetSpans().Snapshots()
	wrappedSpans := []managedSpan{}
	for _, span := range spans {
		wrappedSpans = append(wrappedSpans, managedSpan{span: span})
	}
	return wrappedSpans
}

func TestDefaultHeaderValue(t *testing.T) {
	enc := &samplingHeaderEncoder{}

	result := enc.encode([]managedSpan{})
	if result != "1.0:0" {
		t.Errorf("Expected '1.0:0', got %s", result)
	}
}

func TestMissingAttribute(t *testing.T) {
	enc := &samplingHeaderEncoder{}
	testExporter := tracetest.NewInMemoryExporter()
	tracerProvider := trace.NewTracerProvider(trace.WithSpanProcessor(trace.NewSimpleSpanProcessor(testExporter)))

	makeSpan(tracerProvider, "test")
	makeSpan(tracerProvider, "test2")
	makeSpan(tracerProvider, "test3")

	spans := getSpans(testExporter)

	result := enc.encode(spans)
	if result != "" {
		t.Errorf("Expected '', got %s", result)
	}
}

type samplingHeaderTestData struct {
	probabilities []float64
	expected      string
}

var attrTests = []samplingHeaderTestData{
	{[]float64{1.0}, "1.0:1"},
	{[]float64{1.0, 1.0, 1.0}, "1.0:3"},
	{[]float64{1.0, 0.1, 1.0}, "0.1:1;1.0:2"},
	{[]float64{0.1, 1.0, 1.0}, "0.1:1;1.0:2"},
	{[]float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9}, "0.1:1;0.2:1;0.3:1;0.4:1;0.5:1;0.6:1;0.7:1;0.8:1;0.9:1"},
}

func TestAttributes(t *testing.T) {
	//prepare huge array of probabilities
	probabilities := make([]float64, 300)
	for i := 0; i < 100; i++ {
		probabilities[i] = 0.1
	}
	for i := 100; i < 150; i++ {
		probabilities[i] = 0.2
	}
	for i := 150; i < 175; i++ {
		probabilities[i] = 0.88
	}
	for i := 175; i < 300; i++ {
		probabilities[i] = 0.456
	}
	attrTests = append(attrTests, samplingHeaderTestData{probabilities, "0.1:100;0.2:50;0.456:125;0.88:25"})

	for _, item := range attrTests {
		enc := &samplingHeaderEncoder{}
		testExporter := tracetest.NewInMemoryExporter()
		tracerProvider := trace.NewTracerProvider(trace.WithSpanProcessor(trace.NewSimpleSpanProcessor(testExporter)))

		for _, probability := range item.probabilities {
			makeSpan(tracerProvider, "test", probability)
		}

		spans := getSpans(testExporter)
		result := enc.encode(spans)
		if result != item.expected {
			t.Errorf("Expected %v, got %s", item.expected, result)
		}
	}
}

func TestResampled(t *testing.T) {
	enc := &samplingHeaderEncoder{}
	testExporter := tracetest.NewInMemoryExporter()
	tracerProvider := trace.NewTracerProvider(trace.WithSpanProcessor(trace.NewSimpleSpanProcessor(testExporter)))

	makeSpan(tracerProvider, "test", 0.1)
	makeSpan(tracerProvider, "test2", 0.2)
	makeSpan(tracerProvider, "test3", 0.1)

	spans := getSpans(testExporter)
	newProbability := 0.2
	spans[0].samplingProbability = &newProbability

	result := enc.encode(spans)
	if result != "0.1:1;0.2:2" {
		t.Errorf("Expected '0.1:1;0.2:2', got %s", result)
	}
}
