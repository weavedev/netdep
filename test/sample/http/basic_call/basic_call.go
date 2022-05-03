package main

import (
	"net/http"
)

// target: GET example.com
func main() {
	//nolint:errcheck
	http.Get("https://example.com/")
}
