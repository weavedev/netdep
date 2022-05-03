package main

import (
	"context"
	"net/http"
	//"github.com/go-resty/resty/v2"
	//"github.com/hashicorp/go-retryablehttp"
)

func main() {
	// @mark HTTP request to http://httpbin.org/get
	targetUrl := "http://httpbin.org/get"
	// resty client get request
	client := resty.New()
	_, err := client.R().Get(targetUrl)

	// retryablehttp
	retryClient := retryablehttp.NewClient()
	req, err := http.NewRequestWithContext(context.Background(), "GET", targetUrl, nil)
	retryClient.Client.Do(req)
}
