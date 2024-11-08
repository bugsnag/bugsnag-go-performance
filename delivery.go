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

	probabilityHeader := rawResponse.Header.Get(samplingResponseHeader)
	if probabilityHeader != "" {
		value, err := strconv.ParseFloat(probabilityHeader, 64)
		if err == nil {
			if value <= 1.0 && value >= 0.0 {
				probability = &value
			} else {
				Config.Logger.Printf("Invalid sampling probability: %v\n", value)
			}
		}
	}

	return probability
}

type delivery struct {
	uri     string
	headers map[string]string
}

func (d *delivery) sendPayload(headers map[string]string, payload []byte) (*http.Response, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	body := bytes.NewBuffer(payload)
	req, err := http.NewRequest("POST", d.uri, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	for key, val := range headers {
		req.Header.Set(key, val)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending payload: %v", err)
	}

	return resp, nil
}

func createDelivery() *delivery {
	headers := map[string]string{
		"Bugsnag-Api-Key": Config.APIKey,
		"Content-Type":    "application/json",
		"User-Agent":      fmt.Sprintf("%v v%v", sdkName, Version),
	}

	return &delivery{
		uri:     Config.Endpoint,
		headers: headers,
	}
}

func (d *delivery) send(headers map[string]string, payload []byte) (*http.Response, error) {
	newHeaders := map[string]string{}
	newHeaders["Bugsnag-Sent-At"] = time.Now().Format(time.RFC3339)
	// merge constant headers with the headers passed in
	for k, v := range headers {
		newHeaders[k] = v
	}
	for k, v := range d.headers {
		newHeaders[k] = v
	}

	resp, err := d.sendPayload(newHeaders, payload)
	if err != nil {
		return nil, fmt.Errorf("error sending payload: %v", err)
	}
	return resp, nil
}
