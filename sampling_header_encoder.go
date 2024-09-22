package bugsnagperformance

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type samplingHeaderEncoder struct{}

func (enc *samplingHeaderEncoder) encode(spans []managedSpan) string {
	if len(spans) == 0 {
		return "1.0:0"
	}

	mappedValues := map[float64]int{}
	for _, span := range spans {
		attributes := span.span.Attributes()
		found := false
		for _, keyVal := range attributes {
			if keyVal.Key == "bugsnag.sampling.p" {
				// was resampled
				if span.samplingProbability != nil {
					mappedValues[*span.samplingProbability] += 1
				} else {
					mappedValues[keyVal.Value.AsFloat64()] += 1
				}
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

	// Sort the keys so the result is deterministic
	keysSlice := []float64{}
	for key := range mappedValues {
		keysSlice = append(keysSlice, key)
	}
	sort.Float64s(keysSlice)

	valuesSlice := []string{}
	for _, key := range keysSlice {
		keyStr := strconv.FormatFloat(key, 'g', -1, 64)
		if keyStr == "1" {
			keyStr = "1.0"
		}
		valuesSlice = append(valuesSlice, fmt.Sprintf("%+v:%+v", keyStr, mappedValues[key]))
	}

	return strings.Join(valuesSlice[:], ";")
}
