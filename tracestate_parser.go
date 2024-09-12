package bugsnagperformance

import (
	"fmt"
	"strconv"
	"strings"

	"go.opentelemetry.io/otel/trace"
)

type parsedTracestate struct {
	version  *string
	rValue32 *uint32
	rValue64 *uint64
}

func (pts parsedTracestate) isValid() bool {
	return pts.version != nil && (pts.rValue32 != nil || pts.rValue64 != nil)
}

func (pts parsedTracestate) isValue32() bool {
	return pts.rValue32 != nil
}

func (pts parsedTracestate) getRValue32() uint32 {
	return *(pts.rValue32)
}

func (pts parsedTracestate) getRValue64() uint64 {
	return *(pts.rValue64)
}

type tracestateParser struct{}

func (tsp *tracestateParser) parse(tracestate trace.TraceState) (parsedTracestate, error) {
	state := parsedTracestate{}

	sbValues := tracestate.Get("sb")
	if sbValues == "" {
		return state, fmt.Errorf("tracestate does not contain 'sb' key")
	}

	sbParts := strings.Split(sbValues, ";")
	for _, pair := range sbParts {
		splitPair := strings.Split(pair, ":")
		switch splitPair[0] {
		case "v":
			state.version = &splitPair[1]
		case "r32":
			parsedR, err := strconv.ParseUint(splitPair[1], 10, 32)
			if err == nil {
				parsedR32 := uint32(parsedR)
				state.rValue32 = &parsedR32
			}
		case "r64":
			parsedR64, err := strconv.ParseUint(splitPair[1], 10, 64)
			if err == nil {
				state.rValue64 = &parsedR64
			}
		}
	}

	return state, nil
}
