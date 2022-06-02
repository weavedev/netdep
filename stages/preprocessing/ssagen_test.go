// Package stages
// Package preprocessing defines preprocessing of a given Go project directory
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft
package preprocessing

import (
	"path"
	"testing"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/helpers"

	"github.com/stretchr/testify/assert"
)

func TestLoadPackages(t *testing.T) {
	svcDir := path.Join(helpers.RootDir, "test", "sample", "http", "basic_call")
	initial, _ := LoadAndBuildPackages(svcDir, svcDir)

	assert.Equal(t, "main", initial[0].Pkg.Name())
}

func TestLoadPackagesError(t *testing.T) {
	projDir := path.Join(helpers.RootDir, "test", "example", "svc")
	_, err := LoadAndBuildPackages(projDir, projDir)

	assert.Equal(t, "packages contain errors", err.Error())
}
