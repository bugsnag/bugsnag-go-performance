package bugsnagperformance

import (
	"fmt"
	"strings"

	"go.opentelemetry.io/otel/sdk/trace"
)

type samplingHeaderEncoder struct {}

func (enc *samplingHeaderEncoder) encode(spans []trace.ReadOnlySpan) string {
	if len(spans) == 0 {
		return "1.0:0"
	}

	mappedValues := map[string]int{}
	for _, span := range spans {
		attributes := span.Attributes()
		found := false
		for _, keyVal := range attributes {
			if keyVal.Key == "bugsnag.sampling.p" {
				value := keyVal.Value.AsString()
				mappedValues[value] += 1
				found = true
				break
			}
		}

		if !found {
			// Bail if the atrribute is missing; we'll warn about this later as it
      // means something has gone wrong
			return ""
		}
	}

	valuesSlice := []string{}
	for key, val := range mappedValues {
		valuesSlice = append(valuesSlice, fmt.Sprintf("%+v:%+v", key, val))
	}

	return strings.Join(valuesSlice[:], ";")
}