package main

import (
	"net/http"
	"os"
)

// target: GET $FOO
func main() {
	//nolint:errcheck
	http.Get(os.Getenv("FOO"))
}
