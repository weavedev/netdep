package analyze

import (
	"fmt"
	"golang.org/x/tools/go/ssa"
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

func discoverCall(call *ssa.Call) {
	calledFunction, _ := call.Call.Value.(*ssa.Function)
	calledFunctionPackage := calledFunction.Pkg.Pkg.Path()

	fmt.Println("Called function " + calledFunctionPackage + "->" + calledFunction.Name())
	// args

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
