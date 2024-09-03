package bugsnagperformance

import (
	"fmt"
	"strconv"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"
)

type payloadEncoder struct {
}

func (enc *payloadEncoder) encode(spans []trace.ReadOnlySpan) map[string]interface{} {
	encodedResult := map[string]interface{}{}

	spansArray := []map[string]interface{}{}
	for _, span := range spans {
		spansArray = append(spansArray, enc.spanToMap(span))
	}

	fmt.Printf("Resource attributes\n")
	encodedResult["resourceSpans"] = []interface{}{map[string]interface{}{
		"resource": map[string]interface{}{
			"attributes": enc.attributesToSlice(spans[0].Resource().Attributes()),
		},
		"scopeSpans": []interface{}{
			map[string]interface{}{
				"scope": map[string]interface{}{
					"name":    spans[0].InstrumentationScope().Name,
					"version": spans[0].InstrumentationScope().Version,
				},
				"spans": spansArray,
			},
		},
	},
	}

	return encodedResult
}

func (enc *payloadEncoder) spanToMap(span trace.ReadOnlySpan) map[string]interface{} {

	encodedSpan := map[string]interface{}{
		"name":                   span.Name(),
		"kind":                   int(span.SpanKind()),
		"startTimeUnixNano":      strconv.FormatInt(span.StartTime().UnixNano(), 10),
		"endTimeUnixNano":        strconv.FormatInt(span.EndTime().UnixNano(), 10),
		"traceState":             "",
		"droppedAttributesCount": span.DroppedAttributes(),
		"droppedEventsCount":     span.DroppedEvents(),
		"droppedLinksCount":      span.DroppedLinks(),
		"status": map[string]interface{}{
			"code":    span.Status().Code,
			"message": span.Status().Description,
		},
	}

	if span.Parent().SpanID().IsValid() {
		encodedSpan["parentSpanId"] = span.Parent().SpanID().String()
	}
	if span.SpanContext().HasTraceID() {
		encodedSpan["traceId"] = span.SpanContext().TraceID().String()
	}
	if span.SpanContext().HasSpanID() {
		encodedSpan["spanId"] = span.SpanContext().SpanID().String()
	}

	fmt.Printf("Span attributes\n")
	attr := enc.attributesToSlice(span.Attributes())
	encodedSpan["attributes"] = attr

	events := enc.eventsToSlice(span.Events())
	encodedSpan["events"] = events

	links := enc.linksToSlice(span.Links())
	encodedSpan["links"] = links

	return encodedSpan
}

func (enc *payloadEncoder) attributesToSlice(attr []attribute.KeyValue) []map[string]interface{} {
	encodedAttr := []map[string]interface{}{}

	for _, keyVal := range attr {
		encodedAttr = append(encodedAttr, enc.attributeToMap(keyVal))
	}

	fmt.Printf("ATTRIBUTES: %+v\n", encodedAttr)

	return encodedAttr
}

func (enc *payloadEncoder) attributeToMap(kv attribute.KeyValue) map[string]interface{} {
	singleAttr := map[string]interface{}{}

	singleAttr["key"] = kv.Key
	singleAttr["value"] = enc.attributeValueToMap(kv.Value)

	return singleAttr
}

func (enc *payloadEncoder) attributeValueToMap(val attribute.Value) map[string]interface{} {
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

func (enc *payloadEncoder) eventsToSlice(events []trace.Event) []map[string]interface{} {
	encodedEvents := []map[string]interface{}{}

	fmt.Printf("Event attributes\n")
	for _, event := range events {
		encodedEvents = append(encodedEvents, map[string]interface{}{
			"name":         event.Name,
			"timeUnixNano": event.Time.UnixNano(),
			"attributes":   enc.attributesToSlice(event.Attributes),
		})
	}

	return encodedEvents
}

func (enc *payloadEncoder) linksToSlice(links []trace.Link) []map[string]interface{} {
	encodedLinks := []map[string]interface{}{}

	fmt.Printf("Link attributes\n")
	for _, link := range links {
		encodedLinks = append(encodedLinks, map[string]interface{}{
			"traceId":    link.SpanContext.TraceID().String(),
			"spanId":     link.SpanContext.SpanID().String(),
			"traceState": link.SpanContext.TraceState().String(),
			"attributes": enc.attributesToSlice(link.Attributes),
		})
	}

	return encodedLinks
}
