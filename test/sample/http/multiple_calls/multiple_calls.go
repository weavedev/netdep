package main

import (
	"net/http"
	"net/url"
)

// target: GET example.com and POST example2.com
func main() {
	http.Get("http://example.com/")
	http.PostForm("http://example2.com/form", url.Values{"key": {"Value"}, "id": {"123"}})
}
