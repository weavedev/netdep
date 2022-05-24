//nolint
package main

import (
	"net/http"
)

func wrappedGetCall(url string) {
	location := url + "endpoint"
	http.Get(location)
}

func recurse(url string, depth int) {
	if depth > 0 {
		recurse(url, depth-1)
	} else {
		wrappedGetCall(url)
	}
}

const unresolved = "UNKNOWN"

// target: GET example.com
func main() {
	recurse(unresolved, 25)
}
