// Package preprocessing defines preprocessing of a given Go project directory
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft
package preprocessing

import (
	"go/token"
	"path"
	"testing"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/helpers"

	"github.com/stretchr/testify/assert"
)

func TestLoadAnnotations(t *testing.T) {
	svcDir := path.Join(helpers.RootDir, path.Join("test/sample", path.Join("http", "basic_call")))
	ann, _ := LoadAnnotations(svcDir, "basic_call")
	expected := &Annotation{
		ServiceName: "basic_call",
		Position: token.Position{
			Filename: path.Join(svcDir, "basic_call.go"),
			Offset:   89,
			Line:     10,
			Column:   2,
		},
		Value: "client https://example.com/",
	}

	assert.Equal(t, expected.ServiceName, ann[0].ServiceName)
	assert.Equal(t, expected.Value, ann[0].Value)
	assert.Equal(t, expected.Position.Filename, ann[0].Position.Filename)
	assert.Equal(t, expected.Position.Line, ann[0].Position.Line)
}

func TestLoadAnnotationsInvalidPath(t *testing.T) {
	ann, err := LoadAnnotations("invalidPath", "serviceName")
	assert.Nil(t, ann)
	assert.NotNil(t, err)
}
