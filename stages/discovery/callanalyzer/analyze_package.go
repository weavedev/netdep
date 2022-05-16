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

// locationIdx Stores Relevant Libraries
// their Relevant Methods and for each method
// a position of location in the Args of ssa.Call

//nolint
var locationIdx = map[string]map[string][]int{
	"net/http": {
		"Get":      []int{0, 1},
		"Post":     []int{0, 1},
		"Put":      []int{0, 1},
		"PostForm": []int{0, 1},
		"Head":     []int{0, 1},
		// "Do":                    []int{0},
		"NewRequest":            []int{1},
		"NewRequestWithContext": []int{2},
		// this is a bit different, as it uses http.Request
		// as an argument. This will be completed in the future.
	},
}

type Caller struct {
	RequestLocation string
	Library         string
	MethodName      string
	// TODO: Add package name, filename, code line
}

func getMainFunction(pkg *ssa.Package) *ssa.Function {
	mainMember, hasMain := pkg.Members["main"]
	if !hasMain {
		return nil
	}
	mainFunction, ok := mainMember.(*ssa.Function)
	if !ok {
		return nil
	}

	return mainFunction
}

func discoverCall(call *ssa.Call) *Caller {
	var caller *Caller

	//nolint
	switch call.Call.Value.(type) {
	case *ssa.Function:
		calledFunction, _ := call.Call.Value.(*ssa.Function)
		calledFunctionPackage := calledFunction.Pkg.Pkg.Path()

		relevantPackage, isRelevantPackage := locationIdx[calledFunctionPackage]
		if isRelevantPackage {
			indices, isRelevantFunction := relevantPackage[calledFunction.Name()]
			if call.Call.Args != nil && isRelevantFunction {
				arguments := resolveVariables(call.Call.Args, indices)
				caller = &Caller{
					RequestLocation: strings.Join(arguments, ""),
					Library:         calledFunctionPackage,
					MethodName:      calledFunction.Name(),
				}
			}
		}

		if calledFunction.Blocks != nil && caller == nil {
			discoverBlocks(calledFunction.Blocks)
		}
		return caller
	default:
		return nil
	}
}

func discoverBlock(block *ssa.BasicBlock) []*Caller {
	if block.Instrs == nil {
		return nil
	}

	var calls []*Caller

	for _, instr := range block.Instrs {
		//nolint // can't rewrite switch with 1 case into if,
		// because .(type) is not allowed outside switch.
		switch instruction := instr.(type) {
		// Every complex Instruction is split into several instructions
		// so even if the call is part of variable assignment
		// or a loop it will be stored as a separate ssa.Call instruction
		case *ssa.Call:
			calls = append(calls, discoverCall(instruction))
		}
	}

	return calls
}

func discoverBlocks(blocks []*ssa.BasicBlock) []*Caller {
	var calls []*Caller

	for _, block := range blocks {
		calls = append(calls, discoverBlock(block)...)
	}

	return calls
}

func AnalyzePackageCalls(pkg *ssa.Package) ([]*Caller, error) {
	mainFunction := getMainFunction(pkg)
	// TODO: Expand with endpoint searching

	if mainFunction == nil {
		return nil, fmt.Errorf("no main function found")
	}

	return discoverBlocks(mainFunction.Blocks), nil
}
