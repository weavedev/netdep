package callanalyzer

import (
	"go/constant"
	"go/token"
	"golang.org/x/tools/go/ssa"
)

func resolveVariable(value ssa.Value) string {
	switch val := value.(type) {
	case *ssa.Parameter:
		return "[[Unknown]]"
	case *ssa.Slice:
		return val.Type().String()
	case *ssa.BinOp:
		switch val.Op {
		case token.ADD:
			return resolveVariable(val.X) + resolveVariable(val.Y)
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

func resolveVariables(parameters []ssa.Value, positions []int) []string {
	stringParameters := make([]string, len(positions))
	for i, idx := range positions {
		stringParameters[i] = resolveVariable(parameters[idx])
	}

	return stringParameters
}

// resolveGinAddrSlice is a hardcoded solution to resolve the port address of a Run command from the "github.com/gin-gonic/gin" library
// TODO: implement a general way for resolving variables in slices
func resolveGinAddrSlice(value ssa.Value) []string {
	switch val := value.(type) {
	case *ssa.Slice:
		switch val1 := val.X.(type) {
		case *ssa.Alloc:
			block := val1.Block()
			for i := range block.Instrs {
				switch instruction := block.Instrs[len(block.Instrs)-1-i].(type) {
				case *ssa.Store:
					switch storeVal := instruction.Val.(type) {
					case *ssa.Const:
						switch storeVal.Value.Kind() {
						case constant.String:
							return []string{constant.StringVal(storeVal.Value)}
						}
					}
				}
			}
		}
	case *ssa.Const:
		return []string{":8080"}
	}
	return []string{"var(" + value.Name() + ") = ??"}
}
