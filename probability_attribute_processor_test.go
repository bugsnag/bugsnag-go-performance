package bugsnagperformance

import (
	"testing"

	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestDoesNotReplaceExistingAttribute(t *testing.T) {
	probMgr := &probabilityManager{}
	probMgr.probability = 0.25
	testProc := createProbabilityAttributeProcessor(probMgr)

	testExporter := tracetest.NewInMemoryExporter()
	options := []sdktrace.TracerProviderOption{
		sdktrace.WithSpanProcessor(testProc),
		sdktrace.WithSpanProcessor(sdktrace.NewSimpleSpanProcessor(testExporter)),
	}
	tracerProvider := sdktrace.NewTracerProvider(options...)

	makeSpan(tracerProvider, "test1", float64(0.1))
	makeSpan(tracerProvider, "test2", float64(1.0))
	makeSpan(tracerProvider, "test3", float64(0.5))

	stubs := testExporter.GetSpans()
	for _, span := range stubs {
		attributeSet := attribute.NewSet(span.Attributes...)
		if attributeSet.HasValue(samplingAttribute) {
			val, ok := attributeSet.Value(samplingAttribute)
			if !ok || val == attribute.Float64Value(0.25) {
				t.Errorf("Expected attribute processor to not overwrite existing value")
			}
		}
	}
}

func TestAddsProbabilityAttribute(t *testing.T) {
	probMgr := &probabilityManager{}
	probMgr.probability = 0.25
	testProc := createProbabilityAttributeProcessor(probMgr)

	testExporter := tracetest.NewInMemoryExporter()
	options := []sdktrace.TracerProviderOption{
		sdktrace.WithSpanProcessor(testProc),
		sdktrace.WithSpanProcessor(sdktrace.NewSimpleSpanProcessor(testExporter)),
	}
	tracerProvider := sdktrace.NewTracerProvider(options...)

	makeSpan(tracerProvider, "test1")
	makeSpan(tracerProvider, "test2")
	makeSpan(tracerProvider, "test3")

	stubs := testExporter.GetSpans()
	for _, span := range stubs {
		attributeSet := attribute.NewSet(span.Attributes...)
		if attributeSet.HasValue(samplingAttribute) {
			val, ok := attributeSet.Value(samplingAttribute)
			if !ok || val != attribute.Float64Value(0.25) {
				t.Errorf("Expected attribute processor to add attribute, got %v", val)
			}
		}
	}
}
