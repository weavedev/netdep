//nolint
package main

import (
	"net/http"
)

func wrappedGetCall(url string) {
	newUrl := url + "endpoint"
	http.Get(newUrl)
}

// target: GET example.com
func main() {
	wrappedGetCall("http://example.com/")
}
