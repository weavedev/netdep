//nolint
package main

import (
	"net/http"
)

// target: GET example.com
func main() {
	http.Get("http://basic_handle:8080/count")
}
