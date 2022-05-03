package main

import (
	"net/http"
	"os"
)

// target: GET $FOO
func main() {
	http.Get(os.Getenv("FOO"))
}
