package callanalyzer

import (
	"go/constant"
	"go/token"
	"golang.org/x/tools/go/ssa"
)

// Resolves the value of a variable with respect to a specific frame of execution
func resolveVariable(value ssa.Value, fr Frame) string {
	switch val := value.(type) {
	case *ssa.Parameter:
		paramValue, hasValue := fr.Mappings[val.Name()]
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

func resolveVariables(parameters []ssa.Value, fr Frame) []string {
	if parameters == nil {
		return []string{}
	}

	stringParameters := make([]string, len(parameters))
	for i, param := range parameters {
		stringParameters[i] = resolveVariable(param, fr)
	}

	return stringParameters
}
