package main

import (
	"fmt"
	"net/http"
)

// target: GET example.com
func main() {
	//nolint:errcheck
	host := "example.com"
	endpoint := ""
	url := fmt.Sprintf("http://%s/%s", host, endpoint)
	http.Get(url)
}
