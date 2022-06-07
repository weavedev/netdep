//nolint
package main

import "net/http"

func main() {
	http.Get(getEndpoint())
}

func getEndpoint() string {
	http.Get("https://example.com")
	return "/endpoint"
}
