// Package stages
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft

package stages

import (
	"go/ast"
	"go/token"
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

	_, _, err := LoadServices(projDir, svcDir)

	assert.Equal(t, "no service to analyse were found", err.Error())
}

func TestLoadServices(t *testing.T) {
	_, thisFilePath, _, _ := runtime.Caller(0)
	thisFileParent := path.Dir(thisFilePath)

	projDir := path.Dir(thisFileParent)
	svcDir := path.Join(path.Dir(thisFileParent), "stages")

	services, _, _ := LoadServices(projDir, svcDir)

	assert.Equal(t, "discovery", services[0].Pkg.Name())
	assert.Equal(t, "output", services[1].Pkg.Name())
}

func TestLoadServicesError(t *testing.T) {
	_, thisFilePath, _, _ := runtime.Caller(0)
	thisFileParent := path.Dir(thisFilePath)

	projDir := path.Dir(thisFileParent)
	svcDir := path.Join(path.Dir(thisFileParent), "test", "example", "svc")

	_, _, err := LoadServices(projDir, svcDir)

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

func TestLoadAnnotations(t *testing.T) {
	_, thisFilePath, _, _ := runtime.Caller(0)
	thisFileParent := path.Dir(thisFilePath)
	svcDir := path.Join(path.Dir(thisFileParent), path.Join("test/sample", path.Join("http", "basic_call")))
	ann, _ := LoadAnnotations(svcDir, "basic_call")
	expected := &Annotation{
		ServiceName: "basic_call",
		Position: token.Position{
			Filename: path.Join(svcDir, "basic_call.go"),
			Offset:   89,
			Line:     10,
			Column:   2,
		},
		Value: "//netdep:caller -s targetService",
	}
	assert.Equal(t, expected, ann[0])
}

func TestLoadAnnotationsInvalidPath(t *testing.T) {
	ann, err := LoadAnnotations("invalidPath", "serviceName")
	assert.Nil(t, ann)
	assert.NotNil(t, err)
}
