package analyze

import (
	"fmt"
	"go/constant"
	"go/token"
	"golang.org/x/tools/go/ssa"
	"strings"
)

// interestingCalls Stores Relevant Libraries
// their Relevant Methods and for each method
// a position of location in the Args of ssa.Call
var (
	interestingCalls = map[string]map[string][]int{
		"net/http": {
			"Get":      []int{0},
			"Post":     []int{0},
			"Put":      []int{0},
			"PostForm": []int{0},
			"Do":       []int{0}, // this is a bit different, as it uses http.Request
			// Where 2nd argument of NewRequest is a URL.
		},
	}
)

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

func resolveVariables(parameters []ssa.Value, params map[string]ssa.Value, positions []int) []string {
	stringParameters := make([]string, len(positions))
	for i, idx := range positions {
		stringParameters[i] = resolveVariable(parameters[idx], params)
	}

	return stringParameters
}

func resolveVariable(value ssa.Value, params map[string]ssa.Value) string {
	switch val := value.(type) {
	case *ssa.Parameter:
		paramValue, hasValue := params[val.Name()]
		if hasValue {
			return resolveVariable(paramValue, params)
		} else {
			return "[[Unknown]]"
		}
	case *ssa.BinOp:
		switch val.Op {
		case token.ADD:
			return resolveVariable(val.X, params) + resolveVariable(val.Y, params)
		}
		return "[[OP]]"
	case *ssa.Const:
		switch val.Value.Kind() {
		case constant.String:
			return constant.StringVal(val.Value)
		}
		return "[[CONST]]"
	}

	return "var(" + value.Name() + ") = ??"
}

func discoverCall(call *ssa.Call, params map[string]ssa.Value) {
	switch call.Call.Value.(type) {
	case *ssa.Function:
		calledFunction, _ := call.Call.Value.(*ssa.Function)
		calledFunctionPackage := calledFunction.Pkg.Pkg.Path()

		fmt.Println("Called function " + calledFunctionPackage + "->" + calledFunction.Name())

		interestingPackage, isInterestingPackage := interestingCalls[calledFunctionPackage]
		if isInterestingPackage {
			positions, isInterestingFunction := interestingPackage[calledFunction.Name()]
			if isInterestingFunction {
				fmt.Println("Found call to function " + calledFunctionPackage + "." + calledFunction.Name() + "()")
			}
			if call.Call.Args != nil {
				arguments := resolveVariables(call.Call.Args, params, positions)
				fmt.Println("Arguments: " + strings.Join(arguments, ", "))
			}
		}

		paramMap := make(map[string]ssa.Value)
		for i, param := range calledFunction.Params {
			paramMap[param.Name()] = call.Call.Args[i]
		}

		if calledFunction.Blocks != nil {
			discoverBlocks(calledFunction.Blocks, paramMap)
		}
	}
}

func discoverBlock(block *ssa.BasicBlock, params map[string]ssa.Value) {
	if block.Instrs == nil {
		return
	}

	for _, instr := range block.Instrs {
		switch instruction := instr.(type) {
		// Every complex is split into several instructions
		// so even if the call is part of variable assignment
		// or a loop it will be stored as a separate ssa.Call instruction
		case *ssa.Call:
			discoverCall(instruction, params)
		}
	}
}

func discoverBlocks(blocks []*ssa.BasicBlock, params map[string]ssa.Value) {
	for _, block := range blocks {
		discoverBlock(block, params)
	}
}

func AnalyzePackage(pkg *ssa.Package) {
	mainFunction := getMainFunction(pkg)

	if mainFunction == nil {
		fmt.Println("No main function found!")
		return
	}

	discoverBlocks(mainFunction.Blocks, nil)
}
