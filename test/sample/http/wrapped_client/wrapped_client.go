//nolint
package main

import (
	"net/http"
)

type Client interface {
	Do(*http.Request) (*http.Response, error)
}

func wrappedWrappedGetClient(client http.Client, url string) {
	location := url + "endpoint"
	wrappedGetClient(client, location)
}

func wrappedGetClient(client http.Client, url string) {
	client.Get(url)
}

func wrappedGetCustomClient(client Client, url string) {
	request, _ := http.NewRequest("GET", url, nil)
	client.Do(request)
}

// target: GET example.com
func main() {
	wrappedWrappedGetClient(http.Client{}, "http://example.com/")
	wrappedGetCustomClient(&http.Client{}, "http://example.com/")
}
