package bugsnagperformance

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

const (
	PROBABILITY_SCALE_FACTOR_64 float64 = 18_446_744_073_709_551_615 // (2 ** 64) - 1
	PROBABILITY_SCALE_FACTOR_32 float32 = 4_294_967_295              // (2 ** 32) - 1
)

type Sampler struct {
	probMgr *probabilityManager
	parser  *tracestateParser
}

func CreateSampler(probManager *probabilityManager) sdktrace.Sampler {
	sampler := Sampler{
		probMgr: probManager,
		parser:  &tracestateParser{},
	}

	return &sampler
}

func (s *Sampler) ShouldSample(parameters sdktrace.SamplingParameters) sdktrace.SamplingResult {
	// NOTE: the probability could change at any time so we _must_ only read
	//       it once in this method, otherwise we could use different values
	//       for the sampling decision & p value attribute which would result
	//       in inconsistent data
	probability := s.probMgr.getProbability()
	parentSpanCtx := trace.SpanContextFromContext(parameters.ParentContext)
	traceState := parentSpanCtx.TraceState()

	var decision sdktrace.SamplingDecision
	if s.sampleUsingProbabilityAndTrace(probability, traceState, parameters.TraceID) {
		decision = sdktrace.RecordAndSample
	} else {
		decision = sdktrace.Drop
	}

	return sdktrace.SamplingResult{
		Decision:   decision,
		Attributes: []attribute.KeyValue{{Key: "bugsnag.sampling.p", Value: attribute.Float64Value(probability)}},
		Tracestate: traceState,
	}
}

func (s *Sampler) resample(span sdktrace.ReadOnlySpan) bool {
	attributes := attribute.NewSet(span.Attributes()...)

	// sample all spans that are missing the p value attribute
	if attributes.Len() == 0 || !attributes.HasValue("bugsnag.sampling.p") {
		return true
	}

	probability := s.probMgr.getProbability()
	value, _ := attributes.Value("bugsnag.sampling.p")
	value64 := value.AsFloat64()
	if value64 > probability {
		// TODO update the p value attribute
		// but the span is readonly...

		value64 = probability
	}

	return s.sampleUsingProbabilityAndTrace(value64, span.SpanContext().TraceState(), span.SpanContext().TraceID())
}

func (s *Sampler) sampleUsingProbabilityAndTrace(probability float64, traceState trace.TraceState, traceID trace.TraceID) bool {
	parsedState, err := s.parser.parse(traceState)
	if err != nil {
		fmt.Printf("Error parsing tracestate: %v\n", err)
		return true
	}

	if parsedState.isValid() {
		if parsedState.isValue32() {
			rValue := parsedState.getRValue32()
			pValue := uint32(float32(probability) * PROBABILITY_SCALE_FACTOR_32)
			return pValue >= rValue
		} else {
			rValue := parsedState.getRValue64()
			pValue := uint64(probability * PROBABILITY_SCALE_FACTOR_64)
			return pValue >= rValue
		}
	} else {
		var rValue uint64
		err = binary.Read(bytes.NewBuffer(traceID[:]), binary.LittleEndian, &rValue)
		if err != nil {
			fmt.Printf("Error parsing trace ID: %v\n", err)
			return true
		}
		pValue := uint64(probability * PROBABILITY_SCALE_FACTOR_64)
		return pValue >= rValue
	}
}

func (s *Sampler) Description() string {
	return "Bugsnag Go Performance SDK Sampler"
}
