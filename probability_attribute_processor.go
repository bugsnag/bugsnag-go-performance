package bugsnagperformance

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type ProbabilityAttributeProcessor struct {
	probabilityManager *probabilityManager
}

func createProbabilityAttributeProcessor(pMgr *probabilityManager) *ProbabilityAttributeProcessor {
	return &ProbabilityAttributeProcessor{
		probabilityManager: pMgr,
	}
}

func (pap *ProbabilityAttributeProcessor) OnStart(parent context.Context, s sdktrace.ReadWriteSpan) {
	attributeSet := attribute.NewSet(s.Attributes()...)
	if !attributeSet.HasValue(samplingAttribute) {
		oldAttributes := s.Attributes()
		newAttributes := append(oldAttributes, attribute.KeyValue{
			Key:   samplingAttribute,
			Value: attribute.Float64Value(pap.probabilityManager.getProbability()),
		})
		s.SetAttributes(newAttributes...)
	}
}

func (pap *ProbabilityAttributeProcessor) OnEnd(s sdktrace.ReadOnlySpan) {}

func (pap *ProbabilityAttributeProcessor) Shutdown(ctx context.Context) error {
	return nil
}

func (pap *ProbabilityAttributeProcessor) ForceFlush(ctx context.Context) error {
	return nil
}
