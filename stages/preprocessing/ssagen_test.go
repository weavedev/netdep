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
	projDir := path.Join(helpers.RootDir, path.Join("test/sample", path.Join("http", "basic_call")))
	initial, _ := LoadPackages(projDir, projDir)

	assert.Equal(t, "main", initial[0].Pkg.Name())
}

func TestLoadPackagesError(t *testing.T) {
	projDir := path.Join(helpers.RootDir, path.Join("test/example", path.Join("svc")))
	_, err := LoadPackages(projDir, projDir)

	assert.Equal(t, "packages contain errors", err.Error())
}
