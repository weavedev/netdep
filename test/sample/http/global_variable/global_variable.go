//nolint
package main

import (
	"net/http"
)

func wrappedGetCall(url string) {
	location := url + "/endpoint"
	http.Get(location)
}

const unresolved = "https://example.com"

func main() {
	wrappedGetCall(unresolved)
}
