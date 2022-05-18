// Package stages
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft

package discovery

import (
	"path"
	"runtime"
	"testing"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages"

	"github.com/stretchr/testify/assert"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/discovery"
)

/*
A test for the sample implementation of the resolution method
*/
func TestDiscovery(t *testing.T) {
	_, thisFilePath, _, _ := runtime.Caller(0)
	thisFileParent := path.Dir(thisFilePath)

	projDir := path.Dir(path.Dir(path.Dir(thisFileParent)))
	svcDir := path.Join(path.Dir(path.Dir(thisFileParent)), "sample", "http")

	initial, _ := stages.LoadServices(projDir, svcDir)
	res, _ := discovery.Discover(initial)
	assert.Equal(t, 12, len(res), "Expect 12 interesting call")
	assert.Equal(t, "(*net/http.Client).Get", res[0].MethodName, "Expect net/http.Client+Do to be called")
}
