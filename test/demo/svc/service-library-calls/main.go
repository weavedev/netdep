package main

/*
This "Service" performs calls to http endpoints using different libraries.
*/

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/go-retryablehttp"
	"log"
	"net/http"
	"time"
)

func checkRetry(ctx context.Context, resp *http.Response, err error) (bool, error) {
	if ctx.Err() != nil {
		return false, ctx.Err()
	}

	return resp.StatusCode != http.StatusAccepted, nil
}

// NewCallbackHandler creates a new callback handler with a http retry client.
func NewCallbackHandler(retryMax int) *http.Client {
	retryClient := retryablehttp.NewClient()
	retryClient.CheckRetry = checkRetry
	retryClient.RetryMax = retryMax
	retryClient.RetryWaitMax = 4 * time.Hour

	httpClient := retryClient.StandardClient()

	return httpClient
}

var url = "http://service-gin-server:80/endpoint/get"

func main() {
	// resty client get request
	client := resty.New()

	// perform resty post request
	_, err := client.R().Post("http://service-gin-server:80/endpoint/post")
	if err != nil {
		fmt.Printf("Error: %v\n", err)

		return
	}

	// perform resty post request again
	_, err = client.R().Post("http://service-gin-server:80/endpoint/post")
	if err != nil {
		fmt.Printf("Error: %v\n", err)

		return
	}

	// initialise retryablehttp client
	retryClient := NewCallbackHandler(2)

	// create a new request object and perform request using client.Do
	req, err := http.NewRequestWithContext(context.Background(), "GET", "http://service-gin-server:80/endpoint/get", nil)
	retryClient.Do(req)

	// perform get request using a global url
	http.Get(url)

	log.Fatal(http.ListenAndServe(":80", nil))
}
