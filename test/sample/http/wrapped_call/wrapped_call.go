package main

import (
	"net/http"
)

func wrappedGetCall(url string) {
	http.Get(url)
}

func main() {
	wrappedGetCall("http://example.com/")
}
