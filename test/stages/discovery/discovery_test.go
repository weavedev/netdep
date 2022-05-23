// Package stages
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft

package discovery

import (
	"fmt"
	"os"
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
	resC, _, _ := discovery.Discover(initial)
	assert.Equal(t, 13, len(resC), "Expect 12 interesting call")
	assert.Equal(t, "net/http.Get", resC[0].MethodName, "Expect net/http.Get to be called")
}

func TestDiscoveryBasicCall(t *testing.T) {
	_, thisFilePath, _, _ := runtime.Caller(0)
	thisFileParent := path.Dir(thisFilePath)

	projDir := path.Join(path.Dir(path.Dir(thisFileParent)), path.Join("sample", path.Join("http", "basic_call")))
	initial, _ := stages.LoadPackages(projDir, projDir)
	resC, _, _ := discovery.Discover(initial)
	assert.Equal(t, 1, len(resC), "Expect 1 interesting call")
	assert.Equal(t, "net/http.Get", resC[0].MethodName, "Expect net/http.Get to be called")
}

func TestDiscoveryBasicHandle(t *testing.T) {
	_, thisFilePath, _, _ := runtime.Caller(0)
	thisFileParent := path.Dir(thisFilePath)

	projDir := path.Join(path.Dir(path.Dir(thisFileParent)), path.Join("sample", path.Join("http", "basic_handle")))
	initial, _ := stages.LoadPackages(projDir, projDir)
	_, resS, _ := discovery.Discover(initial)
	assert.Equal(t, 2, len(resS), "Expect 2 interesting calls")
	assert.Equal(t, "net/http.Handle", resS[0].MethodName, "Expect net/http.Handle to be called")
}

func TestDiscoveryBasicHandleFunc(t *testing.T) {
	_, thisFilePath, _, _ := runtime.Caller(0)
	thisFileParent := path.Dir(thisFilePath)

	projDir := path.Join(path.Dir(path.Dir(thisFileParent)), path.Join("sample", path.Join("http", "basic_handlefunc")))
	initial, _ := stages.LoadPackages(projDir, projDir)
	_, resS, _ := discovery.Discover(initial)
	assert.Equal(t, 2, len(resS), "Expect 2 interesting calls")
	assert.Equal(t, "net/http.HandleFunc", resS[0].MethodName, "Expect net/http.HandleFunc to be called")
}

func TestDiscoveryGinHandle(t *testing.T) {
	_, thisFilePath, _, _ := runtime.Caller(0)
	thisFileParent := path.Dir(thisFilePath)

	projDir := path.Join(path.Dir(path.Dir(thisFileParent)), path.Join("sample", path.Join("http", "gin_handle")))
	initial, _ := stages.LoadPackages(projDir, projDir)
	_, resS, _ := discovery.Discover(initial)
	assert.Equal(t, 2, len(resS), "Expect 2 interesting calls")
	assert.Equal(t, "(*github.com/gin-gonic/gin.RouterGroup).GET", resS[0].MethodName, "Expect (*github.com/gin-gonic/gin.RouterGroup).GET to be called")
}

func TestCallInfo(t *testing.T) {
	_, thisFilePath, _, _ := runtime.Caller(0)
	thisFileParent := path.Dir(thisFilePath)

	projDir := path.Dir(path.Dir(path.Dir(thisFileParent)))
	svcDir := path.Join(path.Dir(path.Dir(thisFileParent)), "sample", "http")

	initial, _ := stages.LoadServices(projDir, svcDir)
	res, _, _ := discovery.Discover(initial)
	assert.Equal(t, "multiple_calls", res[5].ServiceName, "Expected service name multiple_calls.go")
	assert.Equal(t, "25", res[7].PositionInFile, "Expected line number 27")
	assert.Equal(t, "multiple_calls"+string(os.PathSeparator)+"multiple_calls.go", res[7].FileName, "Expected file name multiple_calls/multiple_calls.go")
}

func TestWrappedClientCall(t *testing.T) {
	_, thisFilePath, _, _ := runtime.Caller(0)
	thisFileParent := path.Dir(thisFilePath)

	projDir := path.Dir(path.Dir(path.Dir(thisFileParent)))
	svcDir := path.Join(path.Dir(path.Dir(thisFileParent)), "sample", "http", "wrapped_client")

	initial, _ := stages.LoadPackages(projDir, svcDir)
	res, _, _ := discovery.Discover(initial)
	fmt.Println(res)
	//assert.Equal(t, "multiple_calls", res[5].ServiceName, "Expected service name multiple_calls.go")
	//assert.Equal(t, "25", res[7].PositionInFile, "Expected line number 27")
	//assert.Equal(t, "multiple_calls"+string(os.PathSeparator)+"multiple_calls.go", res[7].FileName, "Expected file name multiple_calls/multiple_calls.go")
}
