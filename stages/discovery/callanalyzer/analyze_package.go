/*
Package callanalyzer defines call scanning methods
Copyright © 2022 TW Group 13C, Weave BV, TU Delft
*/
package callanalyzer

import (
	"fmt"

	"golang.org/x/tools/go/ssa"
)

// CallTarget holds information about the destination of a certain call
// executed by the main program; found in the SSA tree.
type CallTarget struct {
	library    string
	MethodName string
	// packageName     string
	// requestLocation string
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
func analyzeCall(call *ssa.Call, frame *Frame, config *AnalyzerConfig, targets *[]*CallTarget) {
	// The fnCallType can be the function type, the anonymous function type, or something else.
	// See https://pkg.go.dev/golang.org/x/tools/go/ssa#Call
	switch fnCallType := call.Call.Value.(type) {
	// TODO: handle other kinds of call _targets
	case *ssa.Function:
		// The full qualified name of the function, including its package
		qualifiedFunctionName := fnCallType.RelString(nil)
		rootPackage := fnCallType.Pkg.Pkg.Path()

		_, isInteresting := config.interestingCalls[qualifiedFunctionName]
		if isInteresting {
			// TODO: Resolve all arguments of the function call (because it is in the interesting list)

			callTarget := &CallTarget{
				library:    rootPackage,
				MethodName: qualifiedFunctionName,
			}
			// fmt.Println("Found call to function " + qualifiedFunctionName)

			*targets = append(*targets, callTarget)
			return
		}

		_, isIgnored := config.ignoreList[rootPackage]

		if isIgnored {
			// Do not recurse into the library if it is ignored (for recursion)
			return
		}
		// Create a copy of the frame for the discovery of child call _targets
		newFrame := *frame

		if fnCallType.Blocks != nil {
			visitBlocks(fnCallType.Blocks, &newFrame, config, targets)
		}
	default:
		return
	}
}

// analyzeBlock checks the type of block, and
func analyzeBlock(block *ssa.BasicBlock, fr *Frame, config *AnalyzerConfig, targets *[]*CallTarget) {
	if block.Instrs == nil {
		return
	}

	for _, instr := range block.Instrs {
		switch instruction := instr.(type) {
		// Every complex Instruction is split into several instructions
		// so even if the call is part of variable assignment
		// or a loop it will be stored as a separate ssa.Call instruction
		case *ssa.Call:
			analyzeCall(instruction, fr, config, targets)
		default:
			continue
		}
	}
}

// visitBlocks
func visitBlocks(blocks []*ssa.BasicBlock, fr *Frame, config *AnalyzerConfig, targets *[]*CallTarget) {
	if len(fr.visited) > config.maxRecDepth {
		// fmt.Println("Visited more than 16 times")
		return
	}

	for _, block := range blocks {
		if fr.hasVisited(block) || block == nil {
			continue
		}
		newFr := fr
		newFr.visited[block] = true
		analyzeBlock(block, newFr, config, targets)
	}
}

// AnalyzePackageCalls takes a main package and finds all 'interesting' methods that are called
func AnalyzePackageCalls(pkg *ssa.Package, config *AnalyzerConfig) ([]*CallTarget, error) {
	mainFunction := findFunctionInPackage(pkg, "main")
	// initFunction := findFunctionInPackage(pkg, "init")

	if mainFunction == nil {
		return nil, fmt.Errorf("no main function found in package %v", pkg)
	}

	baseFrame := Frame{
		visited: make(map[*ssa.BasicBlock]bool, 0),
		// Reference to the final list of all _targets of the entire package
	}

	targets := make([]*CallTarget, 0)

	visitBlocks(mainFunction.Blocks, &baseFrame, config, &targets)

	return targets, nil
}
