// Package stages
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft

package stages

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
A test for the sample implementation of the resolution method
*/
func TestResolving(t *testing.T) {
	res := ResolveEnvVars("../test/example/svc")

	expected := make(map[string]map[string]interface{})

	expected["node-basic-http"] = make(map[string]interface{})
	expected["node-basic-http"]["scopes"] = map[string]interface{}{
		"public": []interface{}{
			map[string]interface{}{
				"endpoint": "/services/ServiceB",
				"name":     "ServiceB",
			},
		},
		"secure": []interface{}{
			map[string]interface{}{
				"endpoint": "/services/ServiceA",
				"name":     "ServiceA",
				"scope":    "service_a",
			},
		},
	}

	assert.Equal(t, expected, res, "Expected the resolution method to return mapped env variables")
}

func TestResolvingInvalid(t *testing.T) {
	res := ResolveEnvVars("../test/example/svc/node-gin-http")
	expected := make(map[string]map[string]interface{})
	assert.Equal(t, expected, res, "Expected the resolution method to return an empty map")
}
