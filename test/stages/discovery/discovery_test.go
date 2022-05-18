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
func TestDiscoveryBasicCall(t *testing.T) {
	_, thisFilePath, _, _ := runtime.Caller(0)
	thisFileParent := path.Dir(thisFilePath)

	projDir := path.Join(path.Dir(path.Dir(thisFileParent)), path.Join("sample", path.Join("http", "basic_call")))
	resC, _, _ := discovery.Discover(projDir, projDir)
	assert.Equal(t, 1, len(resC), "Expect 1 interesting call")
	assert.Equal(t, "net/http.Get", resC[0].MethodName, "Expect net/http.Get to be called")
}

func TestDiscoveryBasicHandle(t *testing.T) {
	_, thisFilePath, _, _ := runtime.Caller(0)
	thisFileParent := path.Dir(thisFilePath)

	projDir := path.Join(path.Dir(path.Dir(thisFileParent)), path.Join("sample", path.Join("http", "basic_handle")))
	_, resS, _ := discovery.Discover(projDir, projDir)
	assert.Equal(t, 2, len(resS), "Expect 2 interesting calls")
	assert.Equal(t, "net/http.Handle", resS[0].MethodName, "Expect net/http.Handle to be called")
}

func TestDiscoveryBasicHandleFunc(t *testing.T) {
	_, thisFilePath, _, _ := runtime.Caller(0)
	thisFileParent := path.Dir(thisFilePath)

	projDir := path.Join(path.Dir(path.Dir(thisFileParent)), path.Join("sample", path.Join("http", "basic_handlefunc")))
	_, resS, _ := discovery.Discover(projDir, projDir)
	assert.Equal(t, 2, len(resS), "Expect 2 interesting calls")
	assert.Equal(t, "net/http.HandleFunc", resS[0].MethodName, "Expect net/http.HandleFunc to be called")
}

func TestDiscoveryGinHandle(t *testing.T) {
	_, thisFilePath, _, _ := runtime.Caller(0)
	thisFileParent := path.Dir(thisFilePath)

	projDir := path.Join(path.Dir(path.Dir(thisFileParent)), path.Join("sample", path.Join("http", "gin_handle")))
	_, resS, _ := discovery.Discover(projDir, projDir)
	assert.Equal(t, 2, len(resS), "Expect 2 interesting calls")
	assert.Equal(t, "(*github.com/gin-gonic/gin.RouterGroup).GET", resS[0].MethodName, "Expect (*github.com/gin-gonic/gin.RouterGroup).GET to be called")
}
