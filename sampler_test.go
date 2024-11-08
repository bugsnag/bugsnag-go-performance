package bugsnagperformance

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"testing"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

func TestShouldSampleOnProbability1(t *testing.T) {
	resetEnv()
	tracestate, _ := trace.ParseTraceState("")
	traceID, _ := trace.TraceIDFromHex("0102030405060708090a0b0c0d0e0f10")
	sampler := createSampler(nil)
	decision := sampler.sampleUsingProbabilityAndTrace(1.0, tracestate, traceID)
	if !decision {
		t.Errorf("Expected true, got false")
	}
}

func TestShouldNotSampleOnProbability0(t *testing.T) {
	resetEnv()
	tracestate, _ := trace.ParseTraceState("")
	traceID, _ := trace.TraceIDFromHex("0102030405060708090a0b0c0d0e0f10")
	sampler := createSampler(nil)
	decision := sampler.sampleUsingProbabilityAndTrace(0.0, tracestate, traceID)
	if decision {
		t.Errorf("Expected false, got true")
	}
}

func TestSampleWithSpecificTraceID(t *testing.T) {
	resetEnv()
	tracestate, _ := trace.ParseTraceState("")
	traceID, _ := trace.TraceIDFromHex("2b0eb6c82ae431ad7fdc00306faebef6")
	sampler := createSampler(nil)
	decision := sampler.sampleUsingProbabilityAndTrace(float64(0.5), tracestate, traceID)
	if !decision {
		t.Errorf("Expected true, got false")
	}
}

func TestNotSampleWithSpecificTraceID(t *testing.T) {
	resetEnv()
	tracestate, _ := trace.ParseTraceState("")
	traceID, _ := trace.TraceIDFromHex("98e03bf7fc2715bdcf426f549ca74150")
	sampler := createSampler(nil)
	decision := sampler.sampleUsingProbabilityAndTrace(float64(0.5), tracestate, traceID)
	if decision {
		t.Errorf("Expected false, got true")
	}
}

func TestShouldSampleHalfOfSpans(t *testing.T) {
	resetEnv()
	probMgr := &probabilityManager{}
	probMgr.probability = 0.5

	tracestate, _ := trace.ParseTraceState("")
	ctxConfig := trace.SpanContextConfig{
		TraceState: tracestate,
	}
	spanCtx := trace.NewSpanContext(ctxConfig)
	ctx := trace.ContextWithSpanContext(context.Background(), spanCtx)

	sampler := createSampler(probMgr)

	sampledCounter := 0
	for i := 0; i < 50_000; i++ {
		traceIDBytes := make([]byte, 16)
		rand.Read(traceIDBytes)
		traceID, _ := trace.TraceIDFromHex(hex.EncodeToString(traceIDBytes))

		res := sampler.ShouldSample(sdktrace.SamplingParameters{
			ParentContext: ctx,
			TraceID:       traceID,
			Name:          "test",
			Kind:          trace.SpanKindServer,
			Attributes:    nil,
			Links:         nil,
		})
		if res.Decision == sdktrace.RecordAndSample {
			sampledCounter++
		}
	}

	if sampledCounter < 24_500 || sampledCounter > 25_500 {
		t.Errorf("Expected around half samples, got %d", sampledCounter)
	}
}

func TestResample(t *testing.T) {
	resetEnv()
	probMgr := &probabilityManager{}
	probMgr.probability = 0.5
	sampler := createSampler(probMgr)

	testExporter := tracetest.NewInMemoryExporter()
	tracerProvider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sdktrace.NewSimpleSpanProcessor(testExporter)))

	makeSpan(tracerProvider, "test")
	makeSpan(tracerProvider, "test2", float64(0.1))
	makeSpan(tracerProvider, "test3", float64(1.0))

	// enforce good traceID
	stubs := testExporter.GetSpans()
	traceID, _ := trace.TraceIDFromHex("2b0eb6c82ae431ad7fdc00306faebef6")
	stubs[1].SpanContext = trace.SpanContext{}.WithTraceID(traceID)
	stubs[2].SpanContext = trace.SpanContext{}.WithTraceID(traceID)

	expected := []bool{true, false, true}

	wrappedSpans := []managedSpan{}
	for _, span := range stubs {
		wrappedSpans = append(wrappedSpans, managedSpan{span: span.Snapshot()})
	}
	for i, span := range wrappedSpans {
		_, result := sampler.resample(span.span)
		if result != expected[i] {
			t.Errorf("Expected %v, got %v, for %v", expected[i], result, i)
		}
	}
}

type samplerTestData struct {
	tracestateStr string
	probability   float64
	expected      sdktrace.SamplingDecision
}

var samplerTests = []samplerTestData{
	{"sb=v:1;r32:1234", 0.00000029, sdktrace.RecordAndSample},
	{"sb=v:1;r32:2000", 0.00000030, sdktrace.Drop},
	{"sb=v:1;r32:999999999", 1.0, sdktrace.RecordAndSample},
	{"sb=v:1;r32:999999999", 0.00005, sdktrace.Drop},
	{"sb=v:1;r64:1234", 0.00000029, sdktrace.RecordAndSample},
	{"sb=v:1;r64:5534023222113", 0.00000030, sdktrace.Drop},
	{"sb=v:1;r64:999999999", 1.0, sdktrace.RecordAndSample},
	{"sb=v:1;r64:999999999", 0.00000000005, sdktrace.Drop},
}

func TestSampleUsingTracestate(t *testing.T) {
	resetEnv()
	for _, item := range samplerTests {
		probMgr := &probabilityManager{}
		probMgr.probability = item.probability

		tracestate, _ := trace.ParseTraceState(item.tracestateStr)
		ctxConfig := trace.SpanContextConfig{
			TraceState: tracestate,
		}
		spanCtx := trace.NewSpanContext(ctxConfig)
		ctx := trace.ContextWithSpanContext(context.Background(), spanCtx)
		traceIDBytes := make([]byte, 16)
		rand.Read(traceIDBytes)
		traceID, _ := trace.TraceIDFromHex(hex.EncodeToString(traceIDBytes))

		sampler := createSampler(probMgr)

		res := sampler.ShouldSample(sdktrace.SamplingParameters{
			ParentContext: ctx,
			TraceID:       traceID,
			Name:          "test",
			Kind:          trace.SpanKindServer,
			Attributes:    nil,
			Links:         nil,
		})
		if res.Decision != item.expected {
			t.Errorf("For item: %+v, expected %+v, got %+v", item, item.expected, res.Decision)
		}
	}
}
