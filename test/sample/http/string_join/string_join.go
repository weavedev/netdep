package main

import (
	"fmt"
	"net/http"
)

func main() {
	host := "example.com"
	endpoint := ""
	url := fmt.Sprintf("http://%s/%s", host, endpoint)
	http.Get(url)
}
