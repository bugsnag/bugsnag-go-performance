package bugsnagperformance

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

type delivery struct {
	uri     string
	headers map[string]string
}

func (d *delivery) sendPayload(payload []byte) error {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	body := bytes.NewBuffer(payload)
	req, err := http.NewRequest("POST", d.uri, body)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	for key, val := range d.headers {
		req.Header.Set(key, val)
	}

	_, err = client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending payload: %v", err)
	}

	return nil
}

func createDelivery(uri, apiKey string) *delivery {
	// TODO - sampling header hardcoded
	headers := map[string]string{
		"Bugsnag-Api-Key":       apiKey,
		"Content-Type":          "application/json",
		"Bugsnag-Span-Sampling": "1.0:0",
	}

	return &delivery{
		uri:     uri,
		headers: headers,
	}
}

func (d *delivery) send(headers map[string]string, payload []byte) error {
	d.headers["Bugsnag-Sent-At"] = time.Now().Format(time.RFC3339)
	for k, v := range headers {
		d.headers[k] = v
	}

	err := d.sendPayload(payload)
	if err != nil {
		return fmt.Errorf("error sending payload: %v", err)
	}
	return nil
}
