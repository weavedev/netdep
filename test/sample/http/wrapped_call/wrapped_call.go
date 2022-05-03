package main

import (
	"net/http"
)

func wrappedGetCall(url string) {
	http.Get(url)
}

// target: GET example.com
func main() {
	wrappedGetCall("http://example.com/")
}
