//nolint
package main

import (
	"net/http"
)

// target: GET example.com
func main() {
	http.Get("https://example.com/")
}
