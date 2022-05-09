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

func getPackageFunction(pkg *ssa.Package, name string) *ssa.Function {
	mainMember, hasMain := pkg.Members[name]
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
	switch fn := call.Call.Value.(type) {
	case *ssa.Function:
		signature := fn.RelString(nil)
		rootPackage := fn.Pkg.Pkg.Path()

		_, isInteresting := interestingCalls[signature]
		if isInteresting {
			if call.Call.Args != nil {
				arguments := resolveVariables(call.Call.Args, fr)
				for _, arg := range call.Call.Args {
					var visited = make([]*ssa.Value, 0)
					findStringConstants(arg, visited)
				}
				fmt.Println("Arguments: " + strings.Join(arguments, ", "))
			}

			fmt.Println("Found call to function " + signature)
			return
		}

		_, isBlacklisted := blackList[rootPackage]

		if isBlacklisted {
			return
		}

		//fr.mappings = make(map[string]ssa.Value)
		//fmt.Println("Called function " + signature + " " + rootPackage)

		for i, param := range fn.Params {
			fr.mappings[param.Name()] = call.Call.Args[i]
		}

		if fn.Blocks != nil {
			discoverBlocks(fn.Blocks, fr)
		}
	}
}

func discoverStore(store *ssa.Store, fr frame) {
	switch val := store.Val.(type) {
	case *ssa.Call:
		fmt.Println(store.Addr.Name() + " = ")
		discoverCall(val, fr)
		break
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
		case *ssa.Store:
			discoverStore(instruction, fr)
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

func Package(pkg *ssa.Package) {
	mainFunction := getPackageFunction(pkg, "main")
	//initFunction := getPackageFunction(pkg, "init")

	if mainFunction == nil {
		fmt.Println("No main function found in package!")
		return
	}

	baseFrame := frame{
		visited:  make([]*ssa.BasicBlock, 0),
		mappings: make(map[string]ssa.Value),
	}

	discoverBlocks(mainFunction.Blocks, baseFrame)

}
