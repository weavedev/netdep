package main

import (
	"net/http"
	"os"
)

func main() {
	http.Get(os.Getenv("FOO"))
}
