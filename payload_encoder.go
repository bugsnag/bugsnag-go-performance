package bugsnagperformance

import (
	"strconv"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"
)

type payloadEncoder struct {
}

func (enc payloadEncoder) encode(spans []wrappedSpan) map[string]interface{} {
	encodedResult := map[string]interface{}{}
	encodedResourceSpans := []interface{}{}
	encodedScopeSpans := []interface{}{}

	spansByResource := enc.sortSpansByResource(spans)
	for _, resourceSpans := range spansByResource {
		if len(resourceSpans) == 0 {
			continue
		}

		spansByScope := enc.sortSpansByScope(resourceSpans)
		for _, scopeSpansArr := range spansByScope {
			if len(scopeSpansArr) == 0 {
				continue
			}
			encodedScopeSpansArr := []map[string]interface{}{}
			for _, scopedSpan := range scopeSpansArr {
				encodedScopeSpansArr = append(encodedScopeSpansArr, enc.spanToMap(scopedSpan))
			}

			encodedScopeSpans = append(encodedScopeSpans, map[string]interface{}{
				"scope": map[string]interface{}{
					"name":    scopeSpansArr[0].roSpan.InstrumentationScope().Name,
					"version": scopeSpansArr[0].roSpan.InstrumentationScope().Version,
				},
				"spans": encodedScopeSpansArr,
			})
		}

		encodedResourceSpans = append(encodedResourceSpans, map[string]interface{}{
			"resource": map[string]interface{}{
				"attributes": enc.attributesToSlice(resourceSpans[0].roSpan.Resource().Attributes()),
			},
			"scopeSpans": encodedScopeSpans,
		})
	}

	encodedResult["resourceSpans"] = encodedResourceSpans

	return encodedResult
}

func (enc payloadEncoder) sortSpansByResource(spans []wrappedSpan) map[attribute.Distinct][]wrappedSpan {
	spansByResource := map[attribute.Distinct][]wrappedSpan{}
	for _, span := range spans {
		mapKey := span.roSpan.Resource().Equivalent()
		if spansByResource[mapKey] == nil {
			spansByResource[mapKey] = []wrappedSpan{}
		}
		spansByResource[mapKey] = append(spansByResource[mapKey], span)
	}
	return spansByResource
}

func (enc payloadEncoder) sortSpansByScope(spans []wrappedSpan) map[string][]wrappedSpan {
	spansByScope := map[string][]wrappedSpan{}
	for _, span := range spans {
		mapKey := span.roSpan.InstrumentationScope().Name
		if spansByScope[mapKey] == nil {
			spansByScope[mapKey] = []wrappedSpan{}
		}
		spansByScope[mapKey] = append(spansByScope[mapKey], span)
	}
	return spansByScope
}

func (enc payloadEncoder) spanToMap(span wrappedSpan) map[string]interface{} {

	encodedSpan := map[string]interface{}{
		"name":                   span.roSpan.Name(),
		"kind":                   int(span.roSpan.SpanKind()),
		"startTimeUnixNano":      strconv.FormatInt(span.roSpan.StartTime().UnixNano(), 10),
		"endTimeUnixNano":        strconv.FormatInt(span.roSpan.EndTime().UnixNano(), 10),
		"droppedAttributesCount": span.roSpan.DroppedAttributes(),
		"droppedEventsCount":     span.roSpan.DroppedEvents(),
		"droppedLinksCount":      span.roSpan.DroppedLinks(),
		"status": map[string]interface{}{
			"code":    span.roSpan.Status().Code,
			"message": span.roSpan.Status().Description,
		},
	}

	if span.roSpan.Parent().SpanID().IsValid() {
		encodedSpan["parentSpanId"] = span.roSpan.Parent().SpanID().String()
	}
	if span.roSpan.SpanContext().HasTraceID() {
		encodedSpan["traceId"] = span.roSpan.SpanContext().TraceID().String()
	}
	if span.roSpan.SpanContext().HasSpanID() {
		encodedSpan["spanId"] = span.roSpan.SpanContext().SpanID().String()
	}
	if traceState := span.roSpan.Parent().TraceState().String(); traceState != "" {
		encodedSpan["traceState"] = traceState
	}

	attr := enc.attributesToSlice(span.roSpan.Attributes())
	// was resampled
	if span.probAttr != nil {
		// update bugsnag.sampling.p attribute
		found := false
		for _, singleAttr := range attr {
			if singleAttr["key"] == "bugsnag.sampling.p" {
				found = true
				singleAttr["value"] = map[string]interface{}{"doubleValue": *span.probAttr}
			}
		}
		if !found {
			attr = append(attr, map[string]interface{}{
				"key": "bugsnag.sampling.p",
				"value": map[string]interface{}{
					"doubleValue": *span.probAttr,
				},
			})
		}
	}
	encodedSpan["attributes"] = attr

	events := enc.eventsToSlice(span.roSpan.Events())
	encodedSpan["events"] = events

	links := enc.linksToSlice(span.roSpan.Links())
	encodedSpan["links"] = links

	return encodedSpan
}

func (enc payloadEncoder) attributesToSlice(attr []attribute.KeyValue) []map[string]interface{} {
	encodedAttr := []map[string]interface{}{}

	for _, keyVal := range attr {
		encodedAttr = append(encodedAttr, enc.attributeToMap(keyVal))
	}
	return encodedAttr
}

func (enc payloadEncoder) attributeToMap(kv attribute.KeyValue) map[string]interface{} {
	singleAttr := map[string]interface{}{}

	singleAttr["key"] = string(kv.Key)
	singleAttr["value"] = enc.attributeValueToMap(kv.Value)

	return singleAttr
}

func (enc payloadEncoder) attributeValueToMap(val attribute.Value) map[string]interface{} {
	singleVal := map[string]interface{}{}
	arrayValues := []interface{}{}

	switch val.Type() {
	case attribute.INT64:
		singleVal["intValue"] = val.AsInt64()
	case attribute.INT64SLICE:
		for _, arrayItem := range val.AsInt64Slice() {
			arrayValues = append(arrayValues, map[string]interface{}{"intValue": arrayItem})
		}
	case attribute.BOOL:
		singleVal["boolValue"] = val.AsBool()
	case attribute.BOOLSLICE:
		for _, arrayItem := range val.AsBoolSlice() {
			arrayValues = append(arrayValues, map[string]interface{}{"boolValue": arrayItem})
		}
	case attribute.FLOAT64:
		singleVal["doubleValue"] = val.AsFloat64()
	case attribute.FLOAT64SLICE:
		for _, arrayItem := range val.AsFloat64Slice() {
			arrayValues = append(arrayValues, map[string]interface{}{"doubleValue": arrayItem})
		}
	case attribute.STRING:
		singleVal["stringValue"] = val.AsString()
	case attribute.STRINGSLICE:
		for _, arrayItem := range val.AsStringSlice() {
			arrayValues = append(arrayValues, map[string]interface{}{"stringValue": arrayItem})
		}
	case attribute.INVALID:
		return singleVal
	}

	if len(arrayValues) > 0 {
		singleVal["arrayValue"] = map[string]interface{}{"values": arrayValues}
	}

	return singleVal
}

func (enc payloadEncoder) eventsToSlice(events []trace.Event) []map[string]interface{} {
	encodedEvents := []map[string]interface{}{}

	for _, event := range events {
		encodedEvents = append(encodedEvents, map[string]interface{}{
			"name":         event.Name,
			"timeUnixNano": event.Time.UnixNano(),
			"attributes":   enc.attributesToSlice(event.Attributes),
		})
	}

	return encodedEvents
}

func (enc payloadEncoder) linksToSlice(links []trace.Link) []map[string]interface{} {
	encodedLinks := []map[string]interface{}{}

	for _, link := range links {
		encodedLink := map[string]interface{}{
			"attributes": enc.attributesToSlice(link.Attributes),
		}
		if link.SpanContext.HasTraceID() {
			encodedLink["traceId"] = link.SpanContext.TraceID().String()
		}
		if link.SpanContext.HasSpanID() {
			encodedLink["spanId"] = link.SpanContext.SpanID().String()
		}
		if traceState := link.SpanContext.TraceState(); traceState.String() != "" {
			encodedLink["traceState"] = traceState.String()
		}

		encodedLinks = append(encodedLinks, encodedLink)
	}

	return encodedLinks
}
