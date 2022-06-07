//nolint
package main

import (
	"net/http"
)

type Client interface {
	Do(*http.Request) (*http.Response, error)
}

func wrappedGetCustomClient(client Client, url string) {
	request, _ := http.NewRequest("GET", url, nil)
	client.Do(request)
}

// target: GET example.com/endpoint
func main() {
	wrappedGetCustomClient(&http.Client{}, "http://example.com/endpoint")
}
