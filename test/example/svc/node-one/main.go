package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-resty/resty/v2"

	httpUtils "lab.weave.nl/nid/nid-core/pkg/utilities/http"
)

// http call using http request object
func performRequest(client httpUtils.Client, host string) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://"+host+"/hello", nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)

		return
	}

	newUrl, err := url.Parse("https://" + host + "/hello2")
	req.URL = newUrl

	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("Error: %v\n", err)

		return
	}
}

func main() {
	// resty client get request
	client := resty.New()

	resp, err := client.R().Get("http://httpbin.org/get")
}
