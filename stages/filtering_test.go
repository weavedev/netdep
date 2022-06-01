// Package stages
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft

package stages

import (
	"go/ast"
	"path"
	"testing"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/helpers"

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
	svcDir := path.Join(helpers.RootDir, "test", "empty", "empty")
	_, err := LoadServices(svcDir)
	assert.Equal(t, "no service to analyse were found", err.Error())
}

func TestLoadServices(t *testing.T) {
	svcDir := path.Join(helpers.RootDir, "stages")

	services, _ := LoadServices(svcDir)

	assert.Equal(t, 3, len(services))

	assert.Equal(t, path.Join(helpers.RootDir, "stages", "discovery"), services[0])
	assert.Equal(t, path.Join(helpers.RootDir, "stages", "matching"), services[1])
	assert.Equal(t, path.Join(helpers.RootDir, "stages", "output"), services[2])
}

func TestLoadPackages(t *testing.T) {
	svcDir := path.Join(helpers.RootDir, "test", "sample", "http", "basic_call")
	initial, _ := LoadPackages(svcDir, svcDir)

	assert.Equal(t, "main", initial[0].Pkg.Name())
}

func TestLoadPackagesError(t *testing.T) {
	projDir := path.Join(helpers.RootDir, "test", "example", "svc")
	_, err := LoadPackages(projDir, projDir)

	assert.Equal(t, "packages contain errors", err.Error())
}
