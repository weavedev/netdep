// Package stages
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft

package discovery

import (
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/discovery"
)

/*
A test for the sample implementation of the resolution method
*/
func TestDiscovery(t *testing.T) {
	_, thisFilePath, _, _ := runtime.Caller(0)
	thisFileParent := path.Dir(thisFilePath)

	projDir := path.Join(path.Dir(path.Dir(thisFileParent)), path.Join("sample", path.Join("http", "basic_call")))
	res, _ := discovery.Discover(projDir, projDir)
	assert.Equal(t, 1, len(res), "Expect 1 interesting call")
	assert.Equal(t, "(*net/http.Client).Get", res[0].MethodName, "Expect net/http.Client+Do to be called")
}
