package main

/*
This "Service" does calls to http endpoints using different libraries.
*/

import (
	"context"
	"fmt"
	"net/http"
	"time"

	httpClient "example/pkg/http"
	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/go-retryablehttp"
)

const (
	retryMaxWait = 4 * time.Hour
)

type CallbackHandler struct {
	Client httpClient.Client
}

func checkRetry(ctx context.Context, resp *http.Response, err error) (bool, error) {
	if ctx.Err() != nil {
		return false, ctx.Err()
	}

	return resp.StatusCode != http.StatusAccepted, nil
}

// NewCallbackHandler creates a new callback handler with a http retry client.
func NewCallbackHandler(retryMax int) *CallbackHandler {
	retryClient := retryablehttp.NewClient()
	retryClient.CheckRetry = checkRetry
	retryClient.RetryMax = retryMax
	retryClient.RetryWaitMax = retryMaxWait

	httpClient := retryClient.StandardClient()

	return &CallbackHandler{
		Client: httpClient,
	}
}

func main() {
	// @mark HTTP request to https://httpbin.org/get
	targetUrl := "https://httpbin.org/get"

	// resty client get request
	client := resty.New()
	_, err := client.R().Get(targetUrl)

	if err != nil {
		fmt.Printf("Error: %v\n", err)

		return
	}

	// retryablehttp
	retryClient := NewCallbackHandler(2)
	// @mark HTTP request to https://httpbin.org/get
	req, err := http.NewRequestWithContext(context.Background(), "GET", targetUrl, nil)
	retryClient.Client.Do(req)
}
