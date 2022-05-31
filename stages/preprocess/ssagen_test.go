// Package preprocess defines preprocessing of a given Go project directory
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft
package preprocess

import (
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadPackages(t *testing.T) {
	_, thisFilePath, _, _ := runtime.Caller(0)
	thisFileParent := path.Dir(path.Dir(thisFilePath))

	projDir := path.Join(path.Dir(thisFileParent), path.Join("test/sample", path.Join("http", "basic_call")))
	initial, _ := LoadPackages(projDir, projDir)

	assert.Equal(t, "main", initial[0].Pkg.Name())
}

func TestLoadPackagesError(t *testing.T) {
	_, thisFilePath, _, _ := runtime.Caller(0)
	thisFileParent := path.Dir(path.Dir(thisFilePath))

	projDir := path.Join(path.Dir(thisFileParent), path.Join("test/example", path.Join("svc")))
	_, err := LoadPackages(projDir, projDir)

	assert.Equal(t, "packages contain errors", err.Error())
}
