/*
Package callanalyzer defines call scanning methods
Copyright © 2022 TW Group 13C, Weave BV, TU Delft
*/
package callanalyzer

import (
	"fmt"
	"golang.org/x/tools/go/ssa"
)

// ignoreList is a set of function names to not recurse into
var ignoreList = map[string]bool{
	"fmt":                  true,
	"reflect":              true,
	"net/url":              true,
	"strings":              true,
	"bytes":                true,
	"io":                   true,
	"errors":               true,
	"runtime":              true,
	"math/bits":            true,
	"internal/reflectlite": true,
	"(*sync.Once).Do":      true,
}

// DiscoveryAction indicates what to do when encountering
// a certain call. Used in interestingCalls
type DiscoveryAction int64

const (
	Output     DiscoveryAction = 0
	Substitute                 = 1
)

// interestingCalls is a map from target to action that is to be taken when encountering the target.
// Used internally to distinguish whether a call is to be:
// outputted as a party in a dependency (0) or substituted with a constant (1)
var (
	interestingCalls = map[string]DiscoveryAction{
		//"net/http.NewRequest":   Output,
		"(*net/http.Client).Do": Output,
		"os.Getenv":             Substitute,
	}
	// targets is the list of stuff this package calls
	targets []*Target
)

// Target holds information about the destination of a certain call
// executed by the main program; found in the SSA tree.
type Target struct {
	requestLocation string
	library         string
	MethodName      string
	packageName     string
	// TODO: Add filename and the position in code
}

// findFunctionInPackage finds the method by within the specified package
// Except it only looks for Exported functions
func findFunctionInPackage(pkg *ssa.Package, name string) *ssa.Function {
	member, hasSpecifiedMember := pkg.Members[name]
	if !hasSpecifiedMember {
		return nil
	}
	specifiedFunction, ok := member.(*ssa.Function)
	if !ok {
		// Not a function
		return nil
	}

	return specifiedFunction
}

// analyzeCall recursively traverses the SSA, with call being the starting point,
// and using the environment specified in the frame
// Variables are only resolved if the call is 'interesting'
// Recursion is only continued if the call is not in the 'ignoreList'
func analyzeCall(call *ssa.Call, frame *Frame) {
	// The fnCallType can be the function type, the anonymous function type, or something else.
	// See https://pkg.go.dev/golang.org/x/tools/go/ssa#Call
	switch fnCallType := call.Call.Value.(type) {

	// TODO: handle other kinds of call targets
	case *ssa.Function:
		// The full qualified name of the function, including its package
		qualifiedFunctionName := fnCallType.RelString(nil)
		rootPackage := fnCallType.Pkg.Pkg.Path()

		_, isInteresting := interestingCalls[qualifiedFunctionName]
		if isInteresting {
			//TODO: Resolve all arguments of the function call (because it is in the interesting list)

			callTarget := &Target{
				library:    rootPackage,
				MethodName: qualifiedFunctionName,
			}

			//fmt.Println("Found call to function " + qualifiedFunctionName)

			targets = append(targets, callTarget)
			return
		}

		_, isIgnored := ignoreList[rootPackage]

		if isIgnored {
			// Do not recurse into the library if it is ignored (for recursion)
			return
		}

		//Create a copy of the frame for the discovery of child call targets
		newFrame := *frame

		if fnCallType.Blocks != nil {
			visitBlocks(fnCallType.Blocks, &newFrame)
		}
	}
}

//analyzeBlock checks the type of block, and
func analyzeBlock(block *ssa.BasicBlock, fr *Frame) {
	if block.Instrs == nil {
		return
	}

	for _, instr := range block.Instrs {
		switch instruction := instr.(type) {
		// Every complex Instruction is split into several instructions
		// so even if the call is part of variable assignment
		// or a loop it will be stored as a separate ssa.Call instruction
		case *ssa.Call:
			analyzeCall(instruction, fr)
		default:
			continue
		}
	}
	return
}

// visitBlocks
func visitBlocks(blocks []*ssa.BasicBlock, fr *Frame) {
	if len(fr.visited) > 16 {
		//fmt.Println("Visited more than 16 times")
		return
	}

	for _, block := range blocks {
		if fr.hasVisited(block) || block == nil {
			continue
		}
		newFr := fr
		newFr.visited[block] = true
		analyzeBlock(block, newFr)
	}
}

// AnalyzePackageCalls takes a main package and finds all 'interesting' methods that are called
func AnalyzePackageCalls(pkg *ssa.Package) ([]*Target, error) {
	mainFunction := findFunctionInPackage(pkg, "main")
	//initFunction := findFunctionInPackage(pkg, "init")

	if mainFunction == nil {
		return nil, fmt.Errorf("no main function found in package %v", pkg)
	}

	baseFrame := Frame{
		visited: make(map[*ssa.BasicBlock]bool, 0),
		// Reference to the final list of all targets of the entire package
	}
	targets = make([]*Target, 0)

	visitBlocks(mainFunction.Blocks, &baseFrame)

	return targets, nil
}
