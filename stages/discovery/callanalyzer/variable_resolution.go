package callanalyzer

import (
	"go/constant"
	"strings"

	"golang.org/x/tools/go/ssa"
)

func resolveVariable(value ssa.Value) string {
	switch val := value.(type) {
	case *ssa.Const:
		//nolint
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
		if idx < len(parameters) {
			variable := resolveVariable(parameters[idx])
			if !strings.HasPrefix(variable, "not a constant") {
				stringParameters[i] = variable
			}
		}
	}

	return stringParameters
}
