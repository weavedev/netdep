package CallAnalyzer

import (
	"go/constant"
	"go/token"
	"golang.org/x/tools/go/ssa"
)

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

func resolveVariables(parameters []ssa.Value, params map[string]ssa.Value, positions []int) []string {
	stringParameters := make([]string, len(positions))
	for i, idx := range positions {
		stringParameters[i] = resolveVariable(parameters[idx], params)
	}

	return stringParameters
}
