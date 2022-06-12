// Package preprocessing defines preprocessing of a given Go project directory
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft
package preprocessing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolvingInvalid(t *testing.T) {
	_, err := IndexEnvironmentVariables("../test/example/svc/node-basic-http/values.yaml")
	assert.NotNil(t, err)
	assert.Equal(t, "the file cannot be parsed", err.Error())
}

func TestMapEnvVarFile(t *testing.T) {
	res, _ := IndexEnvironmentVariables("../../test/example/svc/node-basic-http/env")
	expected := make(map[string]map[string]string)
	expected["service1"] = make(map[string]string)
	expected["service1"]["var1"] = "value1"
	expected["service1"]["var2"] = "value2"
	expected["service2"] = make(map[string]string)
	expected["service2"]["var3"] = "value3"
	assert.Equal(t, expected, res, "Expected to return a map of env vars")
}
