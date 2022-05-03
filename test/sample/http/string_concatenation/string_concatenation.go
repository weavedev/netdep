package main

import (
	"net/http"
)

const protocolSeperator string = ":"

func main() {
	protocol := "http"
	host := "example.com"
	endpoint := "/"
	url := protocol + protocolSeperator + host + endpoint
	http.Get(url)
}
