package bugsnagperformance

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestProbabilityManagerSetProbability(t *testing.T) {
	resetEnv()
	testSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(samplingResponseHeader, "0.1234")
	}))
	defer testSrv.Close()
	Config.Endpoint = testSrv.URL

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pMgr := createProbabilityManager(ctx, createDelivery())

	// default probability is 1.0
	if pMgr.getProbability() != 1.0 {
		t.Errorf("Expected probability to be 1.0, got %f", pMgr.getProbability())
	}

	// wait for the first fetch to complete
	time.Sleep(1 * time.Second)
	if pMgr.getProbability() != 0.1234 {
		t.Errorf("Expected probability to be 0.1234, got %f", pMgr.getProbability())
	}

	time.Sleep(2 * time.Second)
	// manual setup
	pMgr.setProbability(0.5)
	if pMgr.getProbability() != 0.5 {
		t.Errorf("Expected probability to be 0.5, got %f", pMgr.getProbability())
	}
}
