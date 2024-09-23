package bugsnagperformance

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type response struct {
	statusCode         int
	samplingProbablity *float64
}

func newParsedResponse(rawResponse http.Response) response {
	probability := parseSamplingProbability(rawResponse)

	return response{
		statusCode:         rawResponse.StatusCode,
		samplingProbablity: probability,
	}
}

func parseSamplingProbability(rawResponse http.Response) *float64 {
	var probability *float64

	probabilityHeader := rawResponse.Header.Get(SAMPLING_PROBABILITY_HEADER)
	if probabilityHeader != "" {
		value, err := strconv.ParseFloat(probabilityHeader, 64)
		if err == nil {
			if value <= 1.0 && value >= 0.0 {
				probability = &value
			} else {
				fmt.Printf("Invalid sampling probability: %v\n", value)
			}
		}
	}

	return probability
}

type delivery struct {
	uri     string
	headers map[string]string
}

func (d *delivery) sendPayload(payload []byte) (*http.Response, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	body := bytes.NewBuffer(payload)
	req, err := http.NewRequest("POST", d.uri, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	for key, val := range d.headers {
		req.Header.Set(key, val)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending payload: %v", err)
	}

	return resp, nil
}

func createDelivery(uri, apiKey string) *delivery {
	headers := map[string]string{
		"Bugsnag-Api-Key": apiKey,
		"Content-Type":    "application/json",
		"User-Agent":      fmt.Sprintf("Go Bugsnag Performance SDK v%v", Version),
	}

	return &delivery{
		uri:     uri,
		headers: headers,
	}
}

func (d *delivery) send(headers map[string]string, payload []byte) (*http.Response, error) {
	d.headers["Bugsnag-Sent-At"] = time.Now().Format(time.RFC3339)
	for k, v := range headers {
		d.headers[k] = v
	}

	resp, err := d.sendPayload(payload)
	if err != nil {
		return nil, fmt.Errorf("error sending payload: %v", err)
	}
	return resp, nil
}
