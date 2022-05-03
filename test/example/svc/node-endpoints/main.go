package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
)

func main() {

	// Void endpoint
	// @mark HTTP endpoint "/bar"
	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	// Endpoint that calls endpoints
	// @mark HTTP endpoint "/foo"
	http.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		// @mark HTTP call to "http://example.com/"
		http.Get("http://example.com/")

		fmt.Fprintf(w, "Hello Foo")
	})

	// http listenAndServe is used in nid-core/svc/connectinmesh/main.go:64
	log.Fatal(http.ListenAndServe(":8080", nil))
}
