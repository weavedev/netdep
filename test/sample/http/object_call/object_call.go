//nolint
package main

import (
	"context"
	"net/http"
)

func main() {
	c := http.Client{}

	// create a new request object and perform request using client.Do
	req, _ := http.NewRequestWithContext(context.Background(), "GET", "http://example.com/", nil)
	//netdep:client http://example.com/
	c.Do(req)
}
