package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

// http call using http request object with context
func performRequestWithContext(client http.Client, host string) {
	// initialise request, based on https://github.com/nID-sourcecode/nid-core/blob/main/svc/connectinmesh/main.go
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://"+host+"/hello", nil)

	// for testing purposes, mutate request object
	newUrl, err := url.Parse("https://" + host + "/hello2")
	req.URL = newUrl

	// @mark HTTP request to https://$OUTSIDE_MESH/hello2
	_, err = client.Do(req)

	if err != nil {
		fmt.Printf("Error: %v\n", err)

		return
	}
}

// http call using http request object
func performRequestWithoutContext(client http.Client, host string) {
	// initialise request, based on https://github.com/nID-sourcecode/nid-core/blob/main/svc/connectinmesh/main.go
	req, err := http.NewRequest(http.MethodGet, "https://"+host+"/hello3", nil)

	// @mark HTTP request to https://$OUTSIDE_MESH/hello3
	_, err = client.Do(req)

	if err != nil {
		fmt.Printf("Error: %v\n", err)

		return
	}
}

func main() {
	// initialise http client
	httpClient := http.Client{}
	// get host from env
	host := os.Getenv("OUTSIDE_MESH")

	// perform request
	performRequestWithContext(httpClient, host)
	performRequestWithoutContext(httpClient, host)
}
