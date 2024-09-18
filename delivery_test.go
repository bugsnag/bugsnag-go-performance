package bugsnagperformance

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParsedEmptyResponse(t *testing.T) {
	rawResponse := http.Response{}
	response := newParsedResponse(rawResponse)

	if response.statusCode != 0 {
		t.Errorf("Expected status code to be 0, got %d", response.statusCode)
	}

	if response.samplingProbablity != nil {
		t.Errorf("Expected sampling probability to be nil, got %f", *response.samplingProbablity)
	}
}

func TestParsedIncorrectHeader(t *testing.T) {
	header := map[string][]string{
		SAMPLING_PROBABILITY_HEADER: {"invalid"},
	}
	rawResponse := http.Response{Header: header, StatusCode: 200}
	response := newParsedResponse(rawResponse)

	if response.statusCode != 200 {
		t.Errorf("Expected status code to be 200, got %d", response.statusCode)
	}

	if response.samplingProbablity != nil {
		t.Errorf("Expected sampling probability to be nil, got %f", *response.samplingProbablity)
	}
}

func TestParsedIncorrectProbability(t *testing.T) {
	header := map[string][]string{
		SAMPLING_PROBABILITY_HEADER: {"2.0"},
	}
	rawResponse := http.Response{Header: header, StatusCode: 200}
	response := newParsedResponse(rawResponse)

	if response.statusCode != 200 {
		t.Errorf("Expected status code to be 200, got %d", response.statusCode)
	}

	if response.samplingProbablity != nil {
		t.Errorf("Expected sampling probability to be nil, got %f", *response.samplingProbablity)
	}
}

func TestParsedCorrectProbability(t *testing.T) {
	header := map[string][]string{
		SAMPLING_PROBABILITY_HEADER: {"0.5"},
	}
	rawResponse := http.Response{Header: header, StatusCode: 200}
	response := newParsedResponse(rawResponse)

	if response.statusCode != 200 {
		t.Errorf("Expected status code to be 200, got %d", response.statusCode)
	}

	if *response.samplingProbablity != float64(0.5) {
		t.Errorf("Expected sampling probability to be 0.5, got %f", *response.samplingProbablity)
	}
}

func TestHeadersPresentAtSend(t *testing.T) {
	testAPIKey := "12356789"

	testSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("key1") != "value1" {
			t.Errorf("Expected header key1 to be value1, got %s", r.Header.Get("key1"))
		}
		if r.Header.Get("Bugsnag-Sent-At") == "" {
			t.Errorf("Expected header Bugsnag-Sent-At to be present")
		}
		if r.Header.Get("Bugsnag-Api-Key") != testAPIKey {
			t.Errorf("Expected header Bugsnag-Api-Key to be %s, got %s", testAPIKey, r.Header.Get("Bugsnag-Api-Key"))
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected header Content-Type to be application/json, got %s", r.Header.Get("Content-Type"))
		}
	}))
	defer testSrv.Close()

	delivery := createDelivery(testSrv.URL, testAPIKey)
	_, err := delivery.send(map[string]string{"key1": "value1"}, []byte("test"))
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
