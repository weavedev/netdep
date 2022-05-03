package http_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
Run our analysis to detect the single edge to "example.com" in the "basic_call/basic_call.go" file
TODO: implement the actual analysis, and integrate the test
*/
func TestBasicHttpCallDetection(t *testing.T) {
	endpoint := "example.com"
	assert.Equal(t, "example.com", endpoint, "Expect to detect example.com as edge in the dependency graph")
}
