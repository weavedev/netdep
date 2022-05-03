package main

import (
	"net/http"
)

const protocolSeperator string = ":"

// target: GET example.com
func main() {
	//nolint:errcheck
	protocol := "http"
	host := "example.com"
	endpoint := "/"
	url := protocol + protocolSeperator + host + endpoint
	http.Get(url)
}
