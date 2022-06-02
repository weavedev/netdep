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
	_, err := FindServices(svcDir)
	assert.Equal(t, "no service to analyse were found", err.Error())
}

func TestLoadServices(t *testing.T) {
	svcDir := path.Join(helpers.RootDir, "stages")

	services, _ := FindServices(svcDir)

	assert.Equal(t, 4, len(services))

	assert.Equal(t, path.Join(helpers.RootDir, "stages", "discovery"), services[0])
	assert.Equal(t, path.Join(helpers.RootDir, "stages", "matching"), services[1])
	assert.Equal(t, path.Join(helpers.RootDir, "stages", "output"), services[2])
	assert.Equal(t, path.Join(helpers.RootDir, "stages", "preprocessing"), services[3])
}
