package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/test", handler)
}

func handler(w http.ResponseWriter, _ *http.Request) {
	// This should be recognized as dependency
	http.Get("https://example.com/")
}
