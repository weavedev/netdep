//nolint
package main

import (
	"net/http"
)

func wrappedGetCall(url string) {
	location := url + "/endpoint"
	http.Get(location)
}

var globalVariable = "https://example.com"

func main() {
	wrappedGetCall(globalVariable)
}
