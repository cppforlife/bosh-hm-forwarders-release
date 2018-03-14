package client

import (
	"net/http"
	"time"
	"log"
	"io/ioutil"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
)

func NewHTTPClient(apiKey string) Client {
	httpClient := retryablehttp.NewClient()
	httpClient.HTTPClient = &http.Client{Timeout: 30*time.Second} // todo better config?
	httpClient.RetryWaitMin = 30 / 7
	httpClient.RetryWaitMax = 30 / 2
	httpClient.RetryMax = 3
	httpClient.Logger = log.New(ioutil.Discard, "", log.LstdFlags)
	return New(apiKey, httpClient)
}
