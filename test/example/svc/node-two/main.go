package main

import (
	"context"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/go-retryablehttp"
)

func main() {
	targetUrl := "http://httpbin.org/get"
	// resty client get request
	client := resty.New()
	resp, err := client.R().Get(targetUrl)

	// retryablehttp
	retryClient := retryablehttp.NewClient()
	req, err := http.NewRequestWithContext(context.Background(), "GET", targetUrl, nil)
	retryClient.Client.Do(req)
}
