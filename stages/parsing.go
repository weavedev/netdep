package stages

import (
	"go/ast"

	"golang.org/x/tools/go/callgraph"
)

/*
Copyright © 2022 Team 1, Weave BV, TU Delft

The Parsing stages involves construction of necessary data structures for the analysis of the target project.
This stages may include constructing the ASTs or some CallGraph structs, which help find usages of
wrapped HTTP types and methods (see https://pkg.go.dev/golang.org/x/tools/go/callgraph).
Refer to the Project plan, chapter 5 for more information.
*/

// CreateCallGraph
/**
A placeholder Call Graph creation method
*/
// TODO: Remove the following line when implementing this method
//goland:noinspection GoUnusedParameter
func CreateCallGraph(astInst []*ast.File) callgraph.Graph {
	return callgraph.Graph{}
}
