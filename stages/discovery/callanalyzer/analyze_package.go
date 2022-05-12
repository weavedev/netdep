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

var blackList = map[string]bool{
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

// interestingCalls Stores Relevant Libraries
// their Relevant Methods and for each method
// a position of location in the Args of ssa.Call
//nolint
var (
	interestingCalls = map[string][]int{
		"(*net/http.Client).Do": {0},
		"os.Getenv":             {0},
		//	"Post":     []int{0},
		//	"Put":      []int{0},
		//	"PostForm": []int{0},
		//	"Do":       []int{0}, // this is a bit different, as it uses http.Request
		//	// Where 2nd argument of NewRequest is a URL.
		//},
		//"github.com/gin-gonic/gin": {
		//	"GET": []int{0, 1},
		//	"Any": []int{0, 1},
		//},
	}
)

type Target struct {
	requestLocation string
	library         string
	methodName      string
	// TODO: Add package name, filename, code line
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

// Finds
func recurseOnTheTarget(call *ssa.Call, fr Frame) {
	switch fn := call.Call.Value.(type) {

	case *ssa.Function:
		signature := fn.RelString(nil)
		rootPackage := fn.Pkg.Pkg.Path()

		_, isInteresting := interestingCalls[signature]
		if isInteresting {

			arguments := resolveVariables(call.Call.Args, fr)
			fmt.Println("Arguments: " + strings.Join(arguments, ", "))

			caller := &Target{
				requestLocation: strings.Join(arguments, "/"),
				library:         rootPackage,
				methodName:      signature,
			}

			fmt.Println("Found call to function " + signature)
			//TODO Handle error
			_ = append(*fr.targets, caller)
			return
		}

		_, isBlacklisted := blackList[rootPackage]

		if isBlacklisted {
			return
		}

		//fr.mappings = make(map[string]ssa.Value)
		fmt.Println("Called function " + signature + " " + rootPackage)

		for i, param := range fn.Params {
			fr.Mappings[param.Name()] = call.Call.Args[i]
		}

		if fn.Blocks != nil {
			discoverBlocks(fn.Blocks, fr)
		}
	}
}

//
func discoverBlock(block *ssa.BasicBlock, fr Frame) {
	if block.Instrs == nil {
		return
	}

	for _, instr := range block.Instrs {
		//nolint // can't rewrite switch with 1 case into if,
		// because .(type) is not allowed outside switch.
		switch instruction := instr.(type) {
		// Every complex Instruction is split into several instructions
		// so even if the call is part of variable assignment
		// or a loop it will be stored as a separate ssa.Call instruction
		case *ssa.Call:
			recurseOnTheTarget(instruction, fr)
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

	// List of stuff this package calls
	targets := make([]*Target, 0)

	baseFrame := Frame{
		visited:  make([]*ssa.BasicBlock, 0),
		Mappings: make(map[string]ssa.Value),
		// Reference to the final list of all targets of the entire package
		targets: &targets,
	}

	discoverBlocks(mainFunction.Blocks, baseFrame)

	return targets, nil
}
