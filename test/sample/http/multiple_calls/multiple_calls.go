//nolint
package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
)

type Client interface {
	Do(*http.Request) (*http.Response, error)
}

// target: GET example.com and POST example2.com
func main() {
	http.Get("http://example.com/")
	http.PostForm("http://example2.com/form", url.Values{"key": {"Value"}, "id": {"123"}})
	r, er := http.NewRequest(http.MethodGet, "https://example.com/hello3", nil)
	if er != nil {
		log.Fatal()
	}
	http.NewRequestWithContext(context.Background(), http.MethodGet, "https://example.com/hello", nil)
	client := &http.Client{}
	client.Get("https://example.com/hello")
	client.Head("https://example.com/hello")
	client.Do(r)
}
