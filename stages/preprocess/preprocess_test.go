// Package preprocess defines preprocessing of a given Go project directory
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft
package preprocess

import (
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadServicesEmpty(t *testing.T) {
	_, thisFilePath, _, _ := runtime.Caller(0)
	thisFileParent := path.Dir(path.Dir(thisFilePath))

	projDir := path.Dir(thisFileParent)
	svcDir := path.Join(path.Dir(thisFileParent), "test", "empty", "empty")

	_, _, err := LoadServices(projDir, svcDir)

	assert.Equal(t, "no service to analyse were found", err.Error())
}

func TestLoadServices(t *testing.T) {
	_, thisFilePath, _, _ := runtime.Caller(0)
	thisFileParent := path.Dir(path.Dir(thisFilePath))

	projDir := path.Dir(thisFileParent)
	svcDir := path.Join(path.Dir(thisFileParent), "stages")

	services, _, _ := LoadServices(projDir, svcDir)

	assert.Equal(t, "discovery", services[0].Pkg.Name())
	assert.Equal(t, "matching", services[1].Pkg.Name())
	assert.Equal(t, "output", services[2].Pkg.Name())
}

func TestLoadServicesError(t *testing.T) {
	_, thisFilePath, _, _ := runtime.Caller(0)
	thisFileParent := path.Dir(path.Dir(thisFilePath))

	projDir := path.Dir(thisFileParent)
	svcDir := path.Join(path.Dir(thisFileParent), "test", "example", "svc")

	_, _, err := LoadServices(projDir, svcDir)

	assert.Equal(t, "packages contain errors", err.Error())
}
