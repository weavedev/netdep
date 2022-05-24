//nolint
package main

import (
	"net/http"
)

// target: GET example.com
func main() {
	//netdep:caller -s targetService
	http.Get("https://example.com/")
}
