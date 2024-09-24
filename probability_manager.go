package bugsnagperformance

import (
	"context"
	"sync"
)

type probabilityManager struct {
	probability        float64
	probabilityFetcher *probabilityFetcher
	mtx                sync.Mutex
}

func createProbabilityManager(ctx context.Context, delivery *delivery) *probabilityManager {
	probMgr := &probabilityManager{
		probability: 1.0,
	}

	// Will fetch value from the server and update the probability on start
	probFetch := createProbabilityFetcher(ctx, delivery, probMgr.setProbability)
	probMgr.probabilityFetcher = probFetch
	go probFetch.fetchProbability()

	return probMgr
}

func (pm *probabilityManager) getProbability() float64 {
	pm.mtx.Lock()
	defer pm.mtx.Unlock()
	return pm.probability
}

func (pm *probabilityManager) setProbability(probability float64) {
	pm.mtx.Lock()
	defer pm.mtx.Unlock()
	pm.probability = probability
	if pm.probabilityFetcher != nil {
		pm.probabilityFetcher.resetInterval()
	}
}
