//nolint
package main

import (
	"net/http"
)

func wrappedWrappedGetClient(client *http.Client, url string) {
	location := url + "endpoint"
	wrappedGetClient(client, location)
}

func wrappedGetClient(client *http.Client, url string) {
	client.Get(url)
}

// target: GET example.com
func main() {
	wrappedWrappedGetClient(&http.Client{}, "http://example.com/")
}
