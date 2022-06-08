// Package preprocessing defines preprocessing of a given Go project directory
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft
package preprocessing

import (
	"path"
	"testing"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/helpers"

	"github.com/stretchr/testify/assert"
)

func TestLoadAnnotations(t *testing.T) {
	svcDir := path.Join(helpers.RootDir, path.Join("test/sample", path.Join("http", "object_call")))
	ann := make(map[string]map[Position]string)
	LoadAnnotations(svcDir, "object_call", ann)
	expected := make(map[string]map[Position]string)
	expected["object_call"] = make(map[Position]string)
	pos := Position{
		Filename: path.Join(helpers.RootDir, "test", "sample", "http", "object_call", "object_call.go"),
		Line:     14,
	}
	expected["object_call"][pos] = "client http://example.com/"

	assert.Equal(t, expected, ann)
}

func TestLoadAnnotationsInvalidPath(t *testing.T) {
	m := make(map[string]map[Position]string)
	err := LoadAnnotations("invalidPath", "serviceName", m)
	assert.NotNil(t, err)
}

func TestLoadHostAnnotations(t *testing.T) {
	svcDir := path.Join(helpers.RootDir, path.Join("test/sample", path.Join("http", "basic_handle")))
	ann := make(map[string]map[Position]string)
	LoadAnnotations(svcDir, "basic_handle", ann)
	expected := make(map[string]map[Position]string)
	expected["basic_handle"] = make(map[Position]string)
	pos := Position{
		Filename: path.Join(helpers.RootDir, "test", "sample", "http", "basic_handle", "basic_handle.go"),
		Line:     25,
	}
	expected["basic_handle"][pos] = "host http://basic_handle:8080"

	assert.Equal(t, expected, ann)
}
