/*
Package callanalyzer defines call scanning methods
Copyright © 2022 TW Group 13C, Weave BV, TU Delft
*/
package callanalyzer

import (
	"fmt"
	"strings"

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
		"(*net/http.Client).Do": Output,
		"os.Getenv":             Substitute,
	}
	// List of stuff this package calls
	Targets []*Target
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

// getPackageFunction finds the method by within the specified package
// Except it only looks for Exported functions
func getPackageFunction(pkg *ssa.Package, name string) *ssa.Function {
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

// recurseOnTheTarget recursively traverses the SSA, with call being the starting point,
// and using the environment specified in the frame
// Variables are only resolved if the call is 'interesting'
// Recursion is only continued if the call is not in the 'ignoreList'
func recurseOnTheTarget(call *ssa.Call, frame Frame) {
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
			// Resolve all arguments of the function call (because it is in the interesting list)
			arguments := resolveVariables(call.Call.Args, frame)
			//fmt.Println("Arguments: " + strings.Join(arguments, ", "))

			callTarget := &Target{
				requestLocation: strings.Join(arguments, "/"),
				library:         rootPackage,
				MethodName:      qualifiedFunctionName,
			}

			//fmt.Println("Found call to function " + qualifiedFunctionName)

			Targets = append(Targets, callTarget)
			return
		}

		_, isIgnored := ignoreList[rootPackage]

		if isIgnored {
			// Do not recurse into the library if it is ignored (for recursion)
			return
		}

		//frame.mappings = make(map[string]ssa.Value)
		//fmt.Println("Called function " + qualifiedFunctionName + " " + rootPackage)

		for i, param := range fnCallType.Params {
			frame.Mappings[param.Name()] = call.Call.Args[i]
		}

		if fnCallType.Blocks != nil {
			discoverBlocks(fnCallType.Blocks, frame)
		}
	}
}

//
func discoverBlock(block *ssa.BasicBlock, fr Frame) {
	if block.Instrs == nil {
		return
	}

	for _, instr := range block.Instrs {
		switch instruction := instr.(type) {
		// Every complex Instruction is split into several instructions
		// so even if the call is part of variable assignment
		// or a loop it will be stored as a separate ssa.Call instruction
		case *ssa.Call:
			recurseOnTheTarget(instruction, fr)
		default:
			continue
		}
	}
	return
}

// discoverBlocks
func discoverBlocks(blocks []*ssa.BasicBlock, fr Frame) []*Target {
	var calls []*Target

	for _, block := range blocks {
		discoverBlock(block, fr)
	}

	return calls
}

// AnalyzePackageCalls takes a main package and finds all 'interesting' methods that are called
func AnalyzePackageCalls(pkg *ssa.Package) ([]*Target, error) {
	mainFunction := getPackageFunction(pkg, "main")
	//initFunction := getPackageFunction(pkg, "init")

	if mainFunction == nil {
		return nil, fmt.Errorf("no main function found in package %v", pkg)
	}

	baseFrame := Frame{
		visited:  make([]*ssa.BasicBlock, 0),
		Mappings: make(map[string]ssa.Value),
		// Reference to the final list of all targets of the entire package
	}
	Targets = make([]*Target, 0)

	discoverBlocks(mainFunction.Blocks, baseFrame)

	return Targets, nil
}
