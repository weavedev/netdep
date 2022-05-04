// Package stages
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft

package stages

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages"
)

/*
A test for the sample implementation of the resolution method
*/
func TestDiscovery(t *testing.T) {
	res := stages.FindCallersForEndpoint("testService", "testEndpoint", "https://example.com")
	assert.Equal(t, []interface{}(nil), res, "Expect nil as the output of the sample FindCallersForEndpoint method")
}
