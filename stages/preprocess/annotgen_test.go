// Package preprocess defines preprocessing of a given Go project directory
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft
package preprocess

import (
	"go/token"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadAnnotations(t *testing.T) {
	_, thisFilePath, _, _ := runtime.Caller(0)
	thisFileParent := path.Dir(path.Dir(thisFilePath))
	svcDir := path.Join(path.Dir(thisFileParent), path.Join("test/sample", path.Join("http", "basic_call")))
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
	assert.Equal(t, expected, ann[0])
}

func TestLoadAnnotationsInvalidPath(t *testing.T) {
	ann, err := LoadAnnotations("invalidPath", "serviceName")
	assert.Nil(t, ann)
	assert.NotNil(t, err)
}
