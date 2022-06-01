package main

import "net/http"

func main() {
	http.HandleFunc("/test", func(writer http.ResponseWriter, request *http.Request) {
		http.Get("https://example.com/")
	})
}
