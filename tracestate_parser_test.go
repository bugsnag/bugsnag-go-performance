package bugsnagperformance

import (
	"testing"

	"go.opentelemetry.io/otel/trace"
)

func TestEmptyTracestate(t *testing.T) {
	resetEnv()
	parser := &tracestateParser{}
	state, _ := trace.ParseTraceState("")
	parsedState := parser.parse(state)

	if parsedState.isValid() != false {
		t.Fatalf("Expected parsed state to be invalid")
	}
	if parsedState.version != nil {
		t.Fatalf("Expected version to be nil, got %v", *parsedState.version)
	}
	if parsedState.rValue32 != nil {
		t.Fatalf("Expected rValue32 to be nil, got %v", *parsedState.rValue32)
	}
	if parsedState.rValue64 != nil {
		t.Fatalf("Expected rValue64 to be nil, got %v", *parsedState.rValue64)
	}
}

func TestTracestateNoSmartbearValues(t *testing.T) {
	resetEnv()
	parser := &tracestateParser{}
	state, _ := trace.ParseTraceState("ab=c:1,xyz=lmn:op")
	parsedState := parser.parse(state)

	if parsedState.isValid() != false {
		t.Fatalf("Expected parsed state to be invalid")
	}
	if parsedState.version != nil {
		t.Fatalf("Expected version to be nil, got %v", *parsedState.version)
	}
	if parsedState.rValue32 != nil {
		t.Fatalf("Expected rValue32 to be nil, got %v", *parsedState.rValue32)
	}
	if parsedState.rValue64 != nil {
		t.Fatalf("Expected rValue64 to be nil, got %v", *parsedState.rValue64)
	}
}

func TestTracestateNoVersion64(t *testing.T) {
	resetEnv()
	parser := &tracestateParser{}
	state, _ := trace.ParseTraceState("ab=c:1,xyz=lmn:op,sb=r64:1234")
	parsedState := parser.parse(state)

	if parsedState.isValid() != false {
		t.Fatalf("Expected parsed state to be invalid")
	}
	if parsedState.version != nil {
		t.Fatalf("Expected version to be nil, got %v", *parsedState.version)
	}
	if parsedState.rValue32 != nil {
		t.Fatalf("Expected rValue32 to be nil, got %v", *parsedState.rValue32)
	}
	if *parsedState.rValue64 != uint64(1234) {
		t.Fatalf("Expected rValue64 to be %v, got %v", uint64(1234), *parsedState.rValue64)
	}
}

func TestTracestateNoVersion32(t *testing.T) {
	resetEnv()
	parser := &tracestateParser{}
	state, _ := trace.ParseTraceState("ab=c:1,xyz=lmn:op,sb=r32:1234")
	parsedState := parser.parse(state)

	if parsedState.isValid() != false {
		t.Fatalf("Expected parsed state to be invalid")
	}
	if parsedState.version != nil {
		t.Fatalf("Expected version to be nil, got %v", *parsedState.version)
	}
	if *parsedState.rValue32 != uint32(1234) {
		t.Fatalf("Expected rValue32 to be %v, got %v", uint32(1234), *parsedState.rValue32)
	}
	if parsedState.rValue64 != nil {
		t.Fatalf("Expected rValue64 to be nil, got %v", *parsedState.rValue64)
	}
}

func TestTracestateNoRValue(t *testing.T) {
	resetEnv()
	parser := &tracestateParser{}
	state, _ := trace.ParseTraceState("ab=c:1,xyz=lmn:op,sb=v:1")
	parsedState := parser.parse(state)

	if parsedState.isValid() != false {
		t.Fatalf("Expected parsed state to be invalid")
	}
	if *parsedState.version != "1" {
		t.Fatalf("Expected version to be '1', got %v", *parsedState.version)
	}
	if parsedState.rValue32 != nil {
		t.Fatalf("Expected rValue32 to be nil, got %v", *parsedState.rValue32)
	}
	if parsedState.rValue64 != nil {
		t.Fatalf("Expected rValue64 to be nil, got %v", *parsedState.rValue64)
	}
}

func TestTracestateFull64(t *testing.T) {
	resetEnv()
	parser := &tracestateParser{}
	state, _ := trace.ParseTraceState("ab=c:1,xyz=lmn:op,sb=v:2;r64:999")
	parsedState := parser.parse(state)

	if parsedState.isValid() == false {
		t.Fatalf("Expected parsed state to be valid")
	}
	if *parsedState.version != "2" {
		t.Fatalf("Expected version to be '2', got %v", *parsedState.version)
	}
	if parsedState.rValue32 != nil {
		t.Fatalf("Expected rValue32 to be nil, got %v", *parsedState.rValue32)
	}
	if *parsedState.rValue64 != uint64(999) {
		t.Fatalf("Expected rValue64 to be 999, got %v", *parsedState.rValue64)
	}
	if parsedState.isValue32() == true {
		t.Fatalf("Expected isValue32 to be false")
	}
	if parsedState.getRValue64() != uint64(999) {
		t.Fatalf("Expected getRValue64 to be 999, got %v", parsedState.getRValue64())
	}
}

func TestTracestateFull32(t *testing.T) {
	resetEnv()
	parser := &tracestateParser{}
	state, _ := trace.ParseTraceState("ab=c:1,xyz=lmn:op,sb=v:2;r32:999")
	parsedState := parser.parse(state)

	if parsedState.isValid() == false {
		t.Fatalf("Expected parsed state to be valid")
	}
	if *parsedState.version != "2" {
		t.Fatalf("Expected version to be '2', got %v", *parsedState.version)
	}
	if *parsedState.rValue32 != uint32(999) {
		t.Fatalf("Expected rValue32 to be 999, got %v", *parsedState.rValue32)
	}
	if parsedState.rValue64 != nil {
		t.Fatalf("Expected rValue64 to be nil, got %v", *parsedState.rValue64)
	}
	if parsedState.isValue32() == false {
		t.Fatalf("Expected isValue32 to be true")
	}
	if parsedState.getRValue32() != uint32(999) {
		t.Fatalf("Expected getRValue32 to be 999, got %v", parsedState.getRValue32())
	}
}
