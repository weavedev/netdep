package main

import (
	httpAlias "net/http"
)

// target: GET example.com
func main() {
	httpAlias.Get("http://example.com/")
}
