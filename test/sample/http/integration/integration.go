//nolint
package main

import (
	"net/http"
	"os"
)

func getURL(endpoint string) string {
	if endpoint == "" {
		return "http://example.com"
	}

	return "http://example.com" + endpoint
}

func wrappedGetCall(url string) {
	location := url + "endpoint"
	http.Get(location)
	http.Post(location, "application/json", nil)
	req, _ := http.NewRequest("GET", url, nil)
	http.DefaultClient.Do(req)
}

var (
	endpoint = os.Getenv("BASE_ENDPOINT")
)

// target: GET example.com
func main() {
	url := getURL(endpoint)
	wrappedGetCall(url)
}
