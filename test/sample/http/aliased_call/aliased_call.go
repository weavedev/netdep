package main

import (
	httpAlias "net/http"
)

// target: GET example.com
func main() {
	//nolint:errcheck
	httpAlias.Get("http://example.com/")
}
