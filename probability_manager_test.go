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
	pMgr := createProbabilityManager(ctx, 50*time.Second, 300*time.Millisecond)

	time.Sleep(2 * time.Second)
	pMgr.setProbability(0.5)
	if pMgr.getProbability() != 0.5 {
		t.Errorf("Expected probability to be 0.5, got %f", pMgr.getProbability())
	}
}
