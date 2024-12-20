package bugsnagperformance

import (
	"encoding/binary"
	"math"
	"math/big"

	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var (
	PROBABILITY_SCALE_FACTOR_64 = new(big.Float).SetUint64(math.MaxUint64) // (2 ** 64) - 1
)

const (
	PROBABILITY_SCALE_FACTOR_32 float64 = 4_294_967_295 // (2 ** 32) - 1
)

type Sampler struct {
	probMgr *probabilityManager
	parser  *tracestateParser
}

func createSampler(probManager *probabilityManager) *Sampler {
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
		Attributes: []attribute.KeyValue{{Key: samplingAttribute, Value: attribute.Float64Value(probability)}},
		Tracestate: traceState,
	}
}

func (s *Sampler) resample(span sdktrace.ReadOnlySpan) (managedSpan, bool) {
	managedSpan := managedSpan{span: span}
	attributes := attribute.NewSet(span.Attributes()...)

	// sample all spans that are missing the p value attribute
	if attributes.Len() == 0 || !attributes.HasValue(samplingAttribute) {
		return managedSpan, true
	}

	probability := s.probMgr.getProbability()
	value, _ := attributes.Value(samplingAttribute)
	value64 := value.AsFloat64()
	if value64 > probability {
		value64 = probability
		managedSpan.samplingProbability = &value64
	}

	result := s.sampleUsingProbabilityAndTrace(value64, span.SpanContext().TraceState(), span.SpanContext().TraceID())
	return managedSpan, result
}

func (s *Sampler) sampleUsingProbabilityAndTrace(probability float64, traceState trace.TraceState, traceID trace.TraceID) bool {
	parsedState := s.parser.parse(traceState)

	if parsedState.isValid() {
		if parsedState.isValue32() {
			rValue := parsedState.getRValue32()
			pValue := uint32(math.Floor(probability * PROBABILITY_SCALE_FACTOR_32))
			return pValue >= rValue
		} else {
			rValue := parsedState.getRValue64()
			probabilityBig := new(big.Float).SetFloat64(probability)
			pValueRes := new(big.Float)
			pValueRes = pValueRes.Mul(probabilityBig, PROBABILITY_SCALE_FACTOR_64)
			pValue, _ := pValueRes.Uint64()
			return pValue >= rValue
		}
	} else {
		traceIDRaw := [16]byte(traceID)
		rValue := binary.BigEndian.Uint64(traceIDRaw[8:])
		probabilityBig := new(big.Float).SetFloat64(probability)
		pValueRes := new(big.Float)
		pValueRes = pValueRes.Mul(probabilityBig, PROBABILITY_SCALE_FACTOR_64)
		pValue, _ := pValueRes.Uint64()
		return pValue >= rValue
	}
}

func (s *Sampler) Description() string {
	return "Bugsnag Go Performance SDK Sampler"
}
