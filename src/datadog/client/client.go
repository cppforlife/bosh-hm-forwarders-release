package client

import (
	"bytes"
	"fmt"
	"encoding/json"
	gourl "net/url"
	"io/ioutil"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
)

const DefaultAPIURL = ""

type Client struct {
	seriesURL string
	eventsURL string
	httpClient *retryablehttp.Client
}

func New(apiKey string, httpClient *retryablehttp.Client) Client {
	baseURL := "https://app.datadoghq.com"

	query := gourl.Values{}
	query.Add("api_key", apiKey)

	return Client{
		seriesURL: fmt.Sprintf("%s/api/v1/series?%s", baseURL, query.Encode()),
		eventsURL: fmt.Sprintf("%s/api/v1/events?%s", baseURL, query.Encode()),
		httpClient: httpClient,
	}
}

type MetricDataPoint [2]float64

type Metric struct {
	Metric string            `json:"metric"`
	Points []MetricDataPoint `json:"points"`
	Tags   []string          `json:"tags"`
}

type seriesReq struct {
	Series []Metric `json:"series"`
}

func (c Client) PostMetrics(metrics []Metric) error {
	reqBody, err := json.Marshal(seriesReq{metrics})
	if err != nil {
		return fmt.Errorf("serializing request: %s", err)
	}

	req, err := retryablehttp.NewRequest("POST", c.seriesURL, bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("building request: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("post metrics: %s", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("post metrics: response status '%s' body: '%s'", resp.Status, body)
	}

	return nil
}

type Event struct {
	Title       string  `json:"title"`
	Text        string  `json:"text"`
	Time        int     `json:"date_happened"`
	Priority    string  `json:"priority"`
	AlertType   string  `json:"alert_type"`
	Tags        []string `json:"tags"`
}

func (c Client) PostEvent(event Event) error {
	reqBody, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("serializing request: %s", err)
	}

	req, err := retryablehttp.NewRequest("POST", c.eventsURL, bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("building request: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("post event: %s", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("post event: response status '%s' body: '%s'", resp.Status, body)
	}

	return nil
}
