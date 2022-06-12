// Package preprocessing defines preprocessing of a given Go project directory
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft
package preprocessing

import (
	"path/filepath"
	"testing"

	"lab.weave.nl/internships/tud-2022/netDep/helpers"

	"github.com/stretchr/testify/assert"
)

func TestLoadServicesEmpty(t *testing.T) {
	svcDir := filepath.Join(helpers.RootDir, "test", "empty", "empty")
	_, err := FindServices(svcDir)
	assert.Equal(t, "no service to analyse were found", err.Error())
}

func TestLoadServices(t *testing.T) {
	svcDir := filepath.Join(helpers.RootDir, "stages")

	services, _ := FindServices(svcDir)

	assert.Equal(t, 4, len(services))

	assert.Equal(t, filepath.Join(helpers.RootDir, "stages", "discovery"), services[0])
	assert.Equal(t, filepath.Join(helpers.RootDir, "stages", "matching"), services[1])
	assert.Equal(t, filepath.Join(helpers.RootDir, "stages", "output"), services[2])
	assert.Equal(t, filepath.Join(helpers.RootDir, "stages", "preprocessing"), services[3])
}

func TestFindServiceCalls(t *testing.T) {
	serviceCallsDir := filepath.Join(helpers.RootDir, "test", "sample", "servicecalls")

	internalCalls, _, _ := FindServiceCalls(serviceCallsDir)

	posOne := IntCall{
		Name:      "FirstMethod",
		NumParams: 3,
	}
	posTwo := IntCall{
		Name:      "SecondMethod",
		NumParams: 3,
	}
	posThree := IntCall{
		Name:      "ThirdMethod",
		NumParams: 3,
	}

	assert.Equal(t, 3, len(internalCalls))

	assert.Equal(t, "test", internalCalls[posOne])
	assert.Equal(t, "test", internalCalls[posTwo])
	assert.Equal(t, "test", internalCalls[posThree])
}

func TestFindServiceCallsEmptyDir(t *testing.T) {
	serviceCallsDir := ""
	internalCalls, serverTargets, _ := FindServiceCalls(serviceCallsDir)

	assert.Equal(t, 0, len(internalCalls))
	assert.Equal(t, 0, len(*serverTargets))
}

func TestFindServiceCallsInvalidDir(t *testing.T) {
	serviceCallsDir := "invalidDir"
	internalCalls, serverTargets, err := FindServiceCalls(serviceCallsDir)

	assert.Equal(t, 0, len(internalCalls))
	assert.Equal(t, 0, len(*serverTargets))
	assert.NotNil(t, err)
}
