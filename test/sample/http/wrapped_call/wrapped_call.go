//nolint
package main

import (
	"net/http"
)

func wrappedGetCall(url string) {
	location := url + "endpoint"
	http.Get(location)
}

// target: GET example.com
func main() {
	wrappedGetCall("http://example.com/")
}
