package callanalyzer

import (
	"fmt"
	"golang.org/x/tools/go/ssa"
	"strings"
)

// locationIdx Stores Relevant Libraries
// their Relevant Methods and for each method
// a position of location in the Args of ssa.Call
var (
	locationIdx = map[string]map[string][]int{
		"net/http": {
			"Get":      []int{0},
			"Post":     []int{0},
			"Put":      []int{0},
			"PostForm": []int{0},
			//"Do":       []int{0},  this is a bit different, as it uses http.Request
			// as an argument. This will be completed in the future.
		},
	}
)

type Caller struct {
	requestLocation string
	library         string
	methodName      string
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
					requestLocation: strings.Join(arguments, "/"),
					library:         calledFunctionPackage,
					methodName:      calledFunction.Name(),
				}
			}
		}

		if calledFunction.Blocks != nil {
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

func AnalyzePackageCalls(pkg *ssa.Package) []*Caller {
	mainFunction := getMainFunction(pkg)

	if mainFunction == nil {
		fmt.Println("No main function found!")
		return nil
	}

	return discoverBlocks(mainFunction.Blocks)
}
