// Package preprocessing defines preprocessing of a given Go project directory
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft
package preprocessing

import (
	"path/filepath"
	"testing"

	"lab.weave.nl/internships/tud-2022/netDep/stages/discovery/callanalyzer"

	"lab.weave.nl/internships/tud-2022/netDep/helpers"

	"github.com/stretchr/testify/assert"
)

func TestLoadAnnotations(t *testing.T) {
	svcDir := filepath.Join(helpers.RootDir, "test", "sample", "http", "object_call")
	ann := make(map[string]map[callanalyzer.Position]string)
	LoadAnnotations(svcDir, "object_call", ann)
	expected := make(map[string]map[callanalyzer.Position]string)
	expected["object_call"] = make(map[callanalyzer.Position]string)
	pos := callanalyzer.Position{
		Filename: filepath.Join("object_call", "object_call.go"),
		Line:     19,
	}
	expected["object_call"][pos] = "client http://example.com/"

	assert.Equal(t, expected, ann)
}

func TestLoadAnnotationsInvalidPath(t *testing.T) {
	m := make(map[string]map[callanalyzer.Position]string)
	err := LoadAnnotations("invalidPath", "serviceName", m)
	assert.NotNil(t, err)
}

func TestLoadHostAnnotations(t *testing.T) {
	svcDir := filepath.Join(helpers.RootDir, "test", "sample", "http", "basic_handle")
	ann := make(map[string]map[callanalyzer.Position]string)
	LoadAnnotations(svcDir, "basic_handle", ann)
	expected := make(map[string]map[callanalyzer.Position]string)
	expected["basic_handle"] = make(map[callanalyzer.Position]string)
	pos := callanalyzer.Position{
		Filename: filepath.Join("basic_handle", "basic_handle.go"),
		Line:     25,
	}
	expected["basic_handle"][pos] = "host http://basic_handle:8080"

	assert.Equal(t, expected, ann)
}
