//nolint
package main

import (
	"net/http"
)

func wrappedGetCall(url string) {
	location := url + "endpoint"
	http.Get(location)
}

// This recursion should NOT be fully resolved, only 1 pass!
func recurse(url string, depth int) {
	if depth > 0 {
		recurse(url, depth-1)
	} else {
		wrappedGetCall(url)
	}
}

const unresolved = "UNKNOWN"

func main() {
	recurse(unresolved, 25)
}
