//nolint
package main

import (
	"net/http"
)

func wrappedGetCall(url string) {
	location := url + "/endpoint"
	http.Get(location)
}

var globalVariable string

func main() {
	globalVariable = "https://example.com"
	wrappedGetCall(globalVariable)

	globalVariable = "https://example2.com"
	wrappedGetCall(globalVariable)
}
