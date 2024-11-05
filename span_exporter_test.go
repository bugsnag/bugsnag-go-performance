package bugsnagperformance

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestSpanExporter(t *testing.T) {
	resetEnv()
	testSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if header := r.Header.Get(samplingRequestHeader); header == "" {
			t.Errorf("Expected %s header to be set", samplingRequestHeader)
		}
	}))
	defer testSrv.Close()
	Config.Endpoint = testSrv.URL

	probMgr := &probabilityManager{probability: 1.0}
	sampler := createSampler(probMgr)

	commonSpanExporterTest(t, probMgr, sampler)
}

type testCustomSampler struct{}

func (s *testCustomSampler) ShouldSample(parameters sdktrace.SamplingParameters) sdktrace.SamplingResult {
	return sdktrace.SamplingResult{
		Decision: sdktrace.RecordAndSample,
	}
}

func (s *testCustomSampler) Description() string {
	return "testCustomSampler"
}

func TestSpanExporterWithCustomSampler(t *testing.T) {
	resetEnv()
	testSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if header := r.Header.Get(samplingRequestHeader); header != "" {
			t.Errorf("Expected %s header not to be set", samplingRequestHeader)
		}
	}))
	defer testSrv.Close()
	Config.Endpoint = testSrv.URL

	// custom sampler
	Config.CustomSampler = &testCustomSampler{}
	probMgr := &probabilityManager{probability: 1.0}
	sampler := createSampler(probMgr)

	commonSpanExporterTest(t, probMgr, sampler)
}

func TestSpanExporterCustomSpanIsOurs(t *testing.T) {
	resetEnv()
	testSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("REQUEST: %v\n", r)
		if header := r.Header.Get(samplingRequestHeader); header == "" {
			t.Errorf("Expected %s header to be set", samplingRequestHeader)
		}
	}))
	defer testSrv.Close()
	Config.Endpoint = testSrv.URL

	probMgr := &probabilityManager{probability: 1.0}
	sampler := createSampler(probMgr)
	Config.CustomSampler = sampler

	commonSpanExporterTest(t, probMgr, sampler)
}

func commonSpanExporterTest(t *testing.T, probMgr *probabilityManager, sampler *Sampler) {
	testExporter := tracetest.NewInMemoryExporter()
	options := []sdktrace.TracerProviderOption{
		sdktrace.WithSpanProcessor(sdktrace.NewSimpleSpanProcessor(testExporter)),
	}
	tracerProvider := sdktrace.NewTracerProvider(options...)

	makeSpan(tracerProvider, "test1", float64(0.1))
	makeSpan(tracerProvider, "test2", float64(1.0))
	makeSpan(tracerProvider, "test3", float64(0.5))

	spanExporter := SpanExporter{
		disabled:                    false,
		loggedFirstBatchDestination: false,
		probabilityManager:          probMgr,
		delivery:                    createDelivery(),
		sampler:                     sampler,
		unmanagedMode:               false,
	}

	spans := testExporter.GetSpans().Snapshots()
	err := spanExporter.ExportSpans(context.Background(), spans)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}
