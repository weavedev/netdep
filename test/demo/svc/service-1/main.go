package main

import (
	"log"
	"net/http"
)

func main() {
	// perform calls to service-gin-server
	// GET request
	http.Get("http://service-gin-server:80/endpoint/get")

	// POST request
	http.Post("http://service-gin-server:80/endpoint/post", "application/json", nil)

	// Start service
	log.Fatal(http.ListenAndServe(":80", nil))
}
