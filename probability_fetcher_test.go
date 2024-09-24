package bugsnagperformance

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFetchCorrectProbability(t *testing.T) {
	resetEnv()
	testSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(samplingResponseHeader, "0.1234")
	}))
	defer testSrv.Close()
	Config.Endpoint = testSrv.URL

	testCallback := func(probability float64) {
		if probability != 0.1234 {
			t.Errorf("Expected probability to be 0.1234, got %f", probability)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pf := createProbabilityFetcherInternal(ctx, 50*time.Second, 300*time.Millisecond, createDelivery(), testCallback)
	pf.fetchProbability()
}

func TestRetriesForError(t *testing.T) {
	resetEnv()
	count := 0
	testSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if count == 5 {
			w.Header().Set(samplingResponseHeader, "0.1234")
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Printf("Nope! Counted %d\n", count)
		}
		count++
	}))
	defer testSrv.Close()
	Config.Endpoint = testSrv.URL

	testCallback := func(probability float64) {
		if probability != 0.1234 {
			t.Errorf("Expected probability to be 0.1234, got %f", probability)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pf := createProbabilityFetcherInternal(ctx, 50*time.Second, 300*time.Millisecond, createDelivery(), testCallback)
	pf.fetchProbability()
}

func TestRetriesOnIncorrectValue(t *testing.T) {
	resetEnv()
	count := 0
	testSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch count {
		case 0:
			w.Header().Set(samplingResponseHeader, "-0.1234")
			w.WriteHeader(http.StatusOK)
			fmt.Printf("Nope! Counted %d\n", count)
		case 1:
			w.Header().Set(samplingResponseHeader, "1.1234")
			w.WriteHeader(http.StatusOK)
			fmt.Printf("Nope! Counted %d\n", count)
		case 2:
			w.Header().Set(samplingResponseHeader, "0.1234")
			w.WriteHeader(http.StatusOK)
		}
		count++
	}))
	defer testSrv.Close()
	Config.Endpoint = testSrv.URL

	testCallback := func(probability float64) {
		if probability != 0.1234 {
			t.Errorf("Expected probability to be 0.1234, got %f", probability)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pf := createProbabilityFetcherInternal(ctx, 50*time.Second, 300*time.Millisecond, createDelivery(), testCallback)
	pf.fetchProbability()
}

func TestUpdateValueAfterRefresh(t *testing.T) {
	resetEnv()
	count := 0
	testSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch count {
		case 0:
			w.Header().Set(samplingResponseHeader, "0.1234")
			w.WriteHeader(http.StatusOK)
		case 1:
			w.Header().Set(samplingResponseHeader, "0.555")
			w.WriteHeader(http.StatusOK)
		}
		count++
	}))
	defer testSrv.Close()
	Config.Endpoint = testSrv.URL

	testCallback := func(probability float64) {
		if count == 1 && probability != 0.1234 {
			t.Errorf("Expected first probability to be 0.1234, got %f", probability)
		}
		if count == 2 && probability != 0.555 {
			t.Errorf("Expected second probability to be 0.555, got %f", probability)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pf := createProbabilityFetcherInternal(ctx, 2*time.Second, 300*time.Millisecond, createDelivery(), testCallback)
	pf.fetchProbability()
	// second fetch should be after 2 seconds
	time.Sleep(3 * time.Second)
}
