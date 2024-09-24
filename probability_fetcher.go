package bugsnagperformance

import (
	"context"
	"fmt"
	"time"
)

var REQUEST_HEADERS = map[string]string{"Bugsnag-Span-Sampling": "1.0:0"}

type probabilityFetcher struct {
	refreshInterval    time.Duration
	retryInterval      time.Duration
	refreshTimer       *time.Timer
	delivery           *delivery
	wakeupChan         chan bool
	onNewValueCallback func(float64)
	mainProgramCtx     context.Context
}

func createProbabilityFetcher(ctx context.Context, delivery *delivery, callback func(float64)) *probabilityFetcher {
	return createProbabilityFetcherInternal(ctx, fetcherRefreshInterval, fetcherRetryInterval, delivery, callback)
}

func createProbabilityFetcherInternal(ctx context.Context, refreshInterval, retryInterval time.Duration, delivery *delivery, callback func(float64)) *probabilityFetcher {
	wakeupchan := make(chan bool)
	sleepFor := time.NewTimer(refreshInterval)

	probFetch := probabilityFetcher{
		delivery:           delivery,
		refreshInterval:    refreshInterval,
		retryInterval:      retryInterval,
		wakeupChan:         wakeupchan,
		onNewValueCallback: callback,
		refreshTimer:       sleepFor,
		mainProgramCtx:     ctx,
	}

	go probFetch.waitForUpdateTimer()

	return &probFetch
}

func (pf *probabilityFetcher) fetchProbability() {
	// TODO think about retry logic - when to stop trying
	found := false
	for !found {
		requestBody := []byte(fetcherRequestBody)
		resp, err := pf.delivery.send(REQUEST_HEADERS, requestBody)

		if err != nil || resp.StatusCode != 200 {
			fmt.Printf("Failed to retrieve a probability value from BugSnag. Retrying in %v.\n", pf.retryInterval.String())
		} else if resp != nil {
			parsedResp := newParsedResponse(*resp)
			// update probability value if it is in the range [0, 1]
			if parsedResp.samplingProbablity != nil {
				found = true
				fmt.Printf("New probability value: %f\n", *parsedResp.samplingProbablity)
				pf.onNewValueCallback(*parsedResp.samplingProbablity)
			}
		}

		time.Sleep(pf.retryInterval)
	}
}

func (pf *probabilityFetcher) resetInterval() {
	pf.wakeupChan <- true
}

func (pf *probabilityFetcher) waitForUpdateTimer() {
	for {
		select {
		case <-pf.mainProgramCtx.Done():
			fmt.Println("Exiting probability fetcher")
			return
		case <-pf.wakeupChan:
			// we received new value, reset the timer
			pf.refreshTimer.Reset(pf.refreshInterval)
		case <-pf.refreshTimer.C:
			fmt.Printf("Timer expired, get probability value again\n")
			// fetchProbability will reset the timer on success
			pf.fetchProbability()
		}
	}
}
