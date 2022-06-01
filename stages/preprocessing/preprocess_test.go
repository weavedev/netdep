// Package preprocessing defines preprocessing of a given Go project directory
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft
package preprocessing

import (
	"path"
	"testing"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/helpers"

	"github.com/stretchr/testify/assert"
)

func TestLoadServicesEmpty(t *testing.T) {
	svcDir := path.Join(helpers.RootDir, "test", "empty", "empty")
	_, _, err := LoadServices(helpers.RootDir, svcDir)

	assert.Equal(t, "no service to analyse were found", err.Error())
}

func TestLoadServices(t *testing.T) {
	svcDir := path.Join(helpers.RootDir, "stages")
	services, _, _ := LoadServices(helpers.RootDir, svcDir)

	assert.Equal(t, "discovery", services[0].Pkg.Name())
	assert.Equal(t, "matching", services[1].Pkg.Name())
	assert.Equal(t, "output", services[2].Pkg.Name())
}

func TestLoadServicesError(t *testing.T) {
	svcDir := path.Join(helpers.RootDir, "test", "example", "svc")
	_, _, err := LoadServices(helpers.RootDir, svcDir)

	assert.Equal(t, "packages contain errors", err.Error())
}
