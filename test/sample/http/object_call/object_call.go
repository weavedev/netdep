//nolint
package main

import (
	"context"
	"lab.weave.nl/internships/tud-2022/netDep/test/sample/servicecalls"
	"net/http"
)

func main() {
	var testService servicecalls.TESTService
	testService.FirstMethod(1, 2, 3)

	c := http.Client{}

	// create a new request object and perform request using client.Do
	req, _ := http.NewRequestWithContext(context.Background(), "GET", "http://example.com/", nil)
	//netdep:client http://example.com/
	c.Do(req)
}
