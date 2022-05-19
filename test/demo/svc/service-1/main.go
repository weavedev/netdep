package main

import (
	"log"
	"net/http"
)

func main() {
	http.Get("http://service-gin-server:80/endpoint/get")
	http.Post("http://service-gin-server:80/endpoint/post", "application/json", nil)

	log.Fatal(http.ListenAndServe(":80", nil))
}
