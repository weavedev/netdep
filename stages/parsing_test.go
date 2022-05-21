// Package stages
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft

package stages

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"golang.org/x/tools/go/callgraph"
)

/*
A test for the sample implementation of the resolution method
*/
func TestParsing(t *testing.T) {
	res := CreateCallGraph(nil)
	assert.Equal(t, callgraph.Graph{}, res, "Expect empty graph as the output of the CreateCallGraph method")
}
