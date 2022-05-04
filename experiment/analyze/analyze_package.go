package analyze

import (
	"fmt"
	"go/constant"
	"go/token"
	"golang.org/x/tools/go/ssa"
	"strings"
)

var (
	interestingCalls = map[string]map[string]bool{
		"net/http": {
			"Get":      true,
			"Post":     true,
			"Put":      true,
			"PostForm": true,
			"Do":       true,
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

func resolveVariables(parameters []ssa.Value, params map[string]ssa.Value) []string {
	stringParameters := make([]string, len(parameters))
	for i, val := range parameters {
		stringParameters[i] = resolveVariable(val, params)
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
	calledFunction, _ := call.Call.Value.(*ssa.Function)
	calledFunctionPackage := calledFunction.Pkg.Pkg.Path()

	fmt.Println("Called function " + calledFunctionPackage + "->" + calledFunction.Name())

	interestingPackage, isInterestingPackage := interestingCalls[calledFunctionPackage]
	if isInterestingPackage {
		_, isInterestingFunction := interestingPackage[calledFunction.Name()]
		if isInterestingFunction {
			fmt.Println("Found call to function " + calledFunctionPackage + "." + calledFunction.Name() + "()")
		}

		if call.Call.Args != nil {
			arguments := resolveVariables(call.Call.Args, params)
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

func discoverBlock(block *ssa.BasicBlock, params map[string]ssa.Value) {
	if block.Instrs == nil {
		return
	}

	for _, instr := range block.Instrs {
		switch instruction := instr.(type) {
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
