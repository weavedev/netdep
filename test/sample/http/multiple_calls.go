package main

import (
	"net/http"
	"net/url"
)

func main() {
	http.Get("http://example.com/")
	http.PostForm("http://example.com/form", url.Values{"key": {"Value"}, "id": {"123"}})
}
