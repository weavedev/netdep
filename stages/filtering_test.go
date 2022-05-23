// Package stages
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft

package stages

import (
	"go/ast"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
A test for the sample implementation of the resolution method
*/
func TestFiltering(t *testing.T) {
	res := ScanAndFilter("test")
	assert.Equal(t, map[string][]*ast.File(nil), res, "Expect nil as the output of the ScanAndFilter method")
}

func TestLoadServicesEmpty(t *testing.T) {
	_, thisFilePath, _, _ := runtime.Caller(0)
	thisFileParent := path.Dir(thisFilePath)

	projDir := path.Dir(thisFileParent)
	svcDir := path.Join(path.Dir(thisFileParent), "test", "empty", "empty")

	_, err := LoadServices(projDir, svcDir)

	assert.Equal(t, "no service to analyse were found", err.Error())
}

func TestLoadServices(t *testing.T) {
	_, thisFilePath, _, _ := runtime.Caller(0)
	thisFileParent := path.Dir(thisFilePath)

	projDir := path.Dir(thisFileParent)
	svcDir := path.Join(path.Dir(thisFileParent), "stages")

	services, _ := LoadServices(projDir, svcDir)

	assert.Equal(t, "discovery", services[0].Pkg.Name())
	assert.Equal(t, "output", services[1].Pkg.Name())
}

func TestLoadServicesError(t *testing.T) {
	_, thisFilePath, _, _ := runtime.Caller(0)
	thisFileParent := path.Dir(thisFilePath)

	projDir := path.Dir(thisFileParent)
	svcDir := path.Join(path.Dir(thisFileParent), "test", "example", "svc")

	_, err := LoadServices(projDir, svcDir)

	assert.Equal(t, "packages contain errors", err.Error())
}

func TestLoadPackages(t *testing.T) {
	_, thisFilePath, _, _ := runtime.Caller(0)
	thisFileParent := path.Dir(thisFilePath)

	projDir := path.Join(path.Dir(thisFileParent), path.Join("test/sample", path.Join("http", "basic_call")))
	initial, _ := LoadPackages(projDir, projDir)

	assert.Equal(t, "main", initial[0].Pkg.Name())
}

func TestLoadPackagesError(t *testing.T) {
	_, thisFilePath, _, _ := runtime.Caller(0)
	thisFileParent := path.Dir(thisFilePath)

	projDir := path.Join(path.Dir(thisFileParent), path.Join("test/example", path.Join("svc")))
	_, err := LoadPackages(projDir, projDir)

	assert.Equal(t, "packages contain errors", err.Error())
}
