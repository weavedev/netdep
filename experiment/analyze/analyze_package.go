package analyze

import (
	"fmt"
	types "go/types"
	"golang.org/x/tools/go/ssa"
	"strings"
)

var (
	interestingCalls = map[string]map[string]bool{
		"net/http": {
			"Get":  true,
			"Post": true,
			"Put":  true,
			"Do":   true,
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

func resolveVariables(parameters []ssa.Value) []string {
	stringParameters := make([]string, len(parameters))
	for i, val := range parameters {
		parameter, ok := val.(*ssa.Parameter)
		if ok {
			stringParameters[i] = resolveParameterVariable(parameter)
		} else {
			stringParameters[i] = "[err]"
		}
	}

	return stringParameters
}

func resolveParameterVariable(parameter *ssa.Parameter) string {
	parameterName := parameter.Name()
	prog := parameter.Parent().Prog
	pkg := parameter.Parent().Pkg
	switch par := parameter.Object().(type) {
	case *types.Var:
		// TODO: find node?
		value, isAddr := prog.VarValue(par, pkg, nil)
		fmt.Println(parameterName, value, isAddr)
	}
	return "var(" + parameterName + ") = ??"
}

func discoverCall(call *ssa.Call) {
	calledFunction, _ := call.Call.Value.(*ssa.Function)
	calledFunctionPackage := calledFunction.Pkg.Pkg.Path()

	fmt.Println("Called function " + calledFunctionPackage + "->" + calledFunction.Name())

	interestingPackage, isInterestingPackage := interestingCalls[calledFunctionPackage]
	if isInterestingPackage {
		_, isInterestingFunction := interestingPackage[calledFunction.Name()]
		if isInterestingFunction {
			fmt.Println("Found call to function " + calledFunctionPackage + "." + calledFunction.Name() + "()")

			if call.Call.Args != nil {
				arguments := resolveVariables(call.Call.Args)
				fmt.Println("Arguments: " + strings.Join(arguments, ", "))
			}
		}
	}

	if calledFunction.Blocks != nil {
		discoverBlocks(calledFunction.Blocks)
	}
}

func discoverBlock(block *ssa.BasicBlock) {
	if block.Instrs == nil {
		return
	}

	for _, instr := range block.Instrs {
		switch instruction := instr.(type) {
		case *ssa.Call:
			discoverCall(instruction)
		}
	}
}

func discoverBlocks(blocks []*ssa.BasicBlock) {
	for _, block := range blocks {
		discoverBlock(block)
	}
}

func AnalyzePackage(pkg *ssa.Package) {
	mainFunction := getMainFunction(pkg)

	if mainFunction == nil {
		fmt.Println("No main function found!")
		return
	}

	discoverBlocks(mainFunction.Blocks)
}
