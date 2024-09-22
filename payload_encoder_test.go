package bugsnagperformance

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var encodedTestLinks = []map[string]interface{}{
	{
		"traceId":    "68656c6c6f0000000000000000000000",
		"spanId":     "0506070800000000",
		"traceState": "p=2,r=5",
		"attributes": []map[string]interface{}{},
	},
	{
		"traceId":    "776f726c640000000000000000000000",
		"spanId":     "0506070800000000",
		"traceState": "p=2,r=5",
		"attributes": []map[string]interface{}{},
	}}

var encodedTestEvents = []map[string]interface{}{
	{
		"name":         "event1",
		"timeUnixNano": int64(957142923000000000),
		"attributes":   []map[string]interface{}{},
	}, {
		"name":         "event2",
		"timeUnixNano": int64(956624523000000000),
		"attributes":   []map[string]interface{}{},
	},
}

func prepareSpanContexts() (trace.SpanContext, trace.SpanContext) {
	traceState := trace.TraceState{}
	traceState, _ = traceState.Insert("r", "5")
	traceState, _ = traceState.Insert("p", "2")

	spCtx1 := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    trace.TraceID{0x68, 0x65, 0x6c, 0x6c, 0x6f},
		SpanID:     trace.SpanID{0x05, 0x06, 0x07, 0x08},
		TraceState: traceState,
	})
	spCtx2 := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    trace.TraceID{0x77, 0x6f, 0x72, 0x6c, 0x64},
		SpanID:     trace.SpanID{0x05, 0x06, 0x07, 0x08},
		TraceState: traceState,
	})
	return spCtx1, spCtx2
}

func prepareLinks() []sdktrace.Link {
	spCtx1, spCtx2 := prepareSpanContexts()
	link1 := sdktrace.Link{
		SpanContext: spCtx1,
		Attributes:  []attribute.KeyValue{},
	}
	link2 := sdktrace.Link{
		SpanContext: spCtx2,
		Attributes:  []attribute.KeyValue{},
	}

	return []sdktrace.Link{link1, link2}
}

func prepareEvents() []sdktrace.Event {
	event1 := sdktrace.Event{
		Name:       "event1",
		Time:       time.Date(2000, time.May, 01, 1, 2, 3, 0, time.UTC),
		Attributes: []attribute.KeyValue{},
	}
	event2 := sdktrace.Event{
		Name:       "event2",
		Time:       time.Date(2000, time.April, 25, 1, 2, 3, 0, time.UTC),
		Attributes: []attribute.KeyValue{},
	}
	return []sdktrace.Event{event1, event2}
}

func TestAttributesToJSON(t *testing.T) {
	pe := &payloadEncoder{}
	attributes := []attribute.KeyValue{
		{
			Key: "key1", Value: attribute.BoolValue(true),
		}, {
			Key: "key2", Value: attribute.BoolSliceValue([]bool{true, false}),
		}, {
			Key: "key3", Value: attribute.Float64Value(3.14),
		}, {
			Key: "key4", Value: attribute.Float64SliceValue([]float64{3.14, 2.71}),
		}, {
			Key: "key5", Value: attribute.Int64Value(4),
		}, {
			Key: "key6", Value: attribute.Int64SliceValue([]int64{4, 5}),
		}, {
			Key: "key7", Value: attribute.StringValue("value1"),
		}, {
			Key: "key8", Value: attribute.StringSliceValue([]string{"1", "2"}),
		},
	}
	encodedAttributes := pe.attributesToSlice(attributes)

	attributesJSON := `[{"key":"key1","value":{"boolValue":true}},{"key":"key2","value":{"arrayValue":{"values":[{"boolValue":true},{"boolValue":false}]}}},{"key":"key3","value":{"doubleValue":3.14}},{"key":"key4","value":{"arrayValue":{"values":[{"doubleValue":3.14},{"doubleValue":2.71}]}}},{"key":"key5","value":{"intValue":4}},{"key":"key6","value":{"arrayValue":{"values":[{"intValue":4},{"intValue":5}]}}},{"key":"key7","value":{"stringValue":"value1"}},{"key":"key8","value":{"arrayValue":{"values":[{"stringValue":"1"},{"stringValue":"2"}]}}}]`

	payload, err := json.Marshal(encodedAttributes)
	if err != nil {
		t.Fatalf("Error encoding attributes: %v\n", err)
	}

	if string(payload) != attributesJSON {
		t.Fatalf("Expected %s, got %s", attributesJSON, string(payload))
	}
}

func TestLinksListEncoding(t *testing.T) {
	pe := &payloadEncoder{}
	output := pe.linksToSlice(prepareLinks())

	if !reflect.DeepEqual(encodedTestLinks, output) {
		t.Fatalf("Expected %#+v, got %#+v", encodedTestLinks, output)
	}
}

func TestEventListEncoding(t *testing.T) {
	pe := &payloadEncoder{}
	output := pe.eventsToSlice(prepareEvents())

	if !reflect.DeepEqual(output, encodedTestEvents) {
		t.Fatalf("Expected %#+v, got %#+v", encodedTestEvents, output)
	}
}
