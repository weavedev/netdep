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
var blackList = map[string]bool{
	"fmt":     true,
	"reflect": true,
}

var (
	interestingCalls = map[string][]int{
		"(*net/http.Client).Do": []int{0},
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

func resolveVariables(parameters []ssa.Value, fr frame) []string {
	stringParameters := make([]string, len(parameters))
	for i, param := range parameters {
		stringParameters[i] = resolveVariable(param, fr)
	}

	return stringParameters
}

func resolveVariable(value ssa.Value, fr frame) string {
	switch val := value.(type) {
	case *ssa.Parameter:
		paramValue, hasValue := fr.mappings[val.Name()]
		if hasValue {
			return resolveVariable(paramValue, fr)
		} else {
			return "[[Unknown]]"
		}
	case *ssa.BinOp:
		switch val.Op {
		case token.ADD:
			return resolveVariable(val.X, fr) + resolveVariable(val.Y, fr)
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

func discoverCall(call *ssa.Call, fr frame) {
	switch call.Call.Value.(type) {
	case *ssa.Function:
		calledFunction, _ := call.Call.Value.(*ssa.Function)
		signature := calledFunction.RelString(nil)
		rootPackage := calledFunction.Pkg.Pkg.Name()
		_, isBlacklisted := blackList[rootPackage]

		if isBlacklisted {
			return
		}

		fmt.Println("Called function " + signature + " " + rootPackage)

		_, isInteresting := interestingCalls[signature]
		if isInteresting {
			if call.Call.Args != nil {
				arguments := resolveVariables(call.Call.Args, fr)
				fmt.Println("Arguments: " + strings.Join(arguments, ", "))
			}

			fmt.Println("Found call to function " + signature)
			return
		}

		//fr.mappings = make(map[string]ssa.Value)

		for i, param := range calledFunction.Params {
			fr.mappings[param.Name()] = call.Call.Args[i]
		}

		if calledFunction.Blocks != nil {
			discoverBlocks(calledFunction.Blocks, fr)
		}
	}
}

func discoverBlock(block *ssa.BasicBlock, fr frame) {
	if block.Instrs == nil {
		return
	}

	if len(fr.visited) > 16 {
		//fmt.Println("Nested > 32")
		return
	}

	for _, instr := range block.Instrs {
		switch instruction := instr.(type) {
		// Every complex is split into several instructions
		// so even if the call is part of variable assignment
		// or a loop it will be stored as a separate ssa.Call instruction
		case *ssa.Call:
			discoverCall(instruction, fr)
		}
	}
}

func discoverBlocks(blocks []*ssa.BasicBlock, fr frame) {
	if len(fr.visited) > 16 {
		//fmt.Println("Nested > 32")
		return
	}

	for _, block := range blocks {
		if fr.hasVisited(block) || block == nil {
			continue
		}
		newFr := fr
		newFr.visited = append(fr.visited, block)
		discoverBlock(block, newFr)
	}
}

func AnalyzePackage(pkg *ssa.Package) {
	mainFunction := getMainFunction(pkg)

	if mainFunction == nil {
		fmt.Println("No main function found!")
		return
	}

	baseFrame := frame{
		visited:  make([]*ssa.BasicBlock, 0),
		mappings: make(map[string]ssa.Value),
	}
	discoverBlocks(mainFunction.Blocks, baseFrame)
}
