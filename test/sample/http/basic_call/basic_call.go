//nolint
package main

import (
	"net/http"
)

// target: GET example.com
func main() {
	//netdep:client https://example.com/
	http.Get("https://example.com/")
}
