package callanalyzer

import (
	"go/constant"
	"golang.org/x/tools/go/ssa"
)

func resolveVariable(value ssa.Value) string {
	switch val := value.(type) {
	case *ssa.Const:
		switch val.Value.Kind() {
		case constant.String:
			return constant.StringVal(val.Value)
		default:
			return "non-string constant"
		}
	default:
		return "not a constant"
	}
}

func resolveVariables(parameters []ssa.Value, positions []int) []string {
	stringParameters := make([]string, len(positions))
	for i, idx := range positions {
		stringParameters[i] = resolveVariable(parameters[idx])
	}

	return stringParameters
}
