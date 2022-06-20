/*
Package callanalyzer defines call scanning methods
Copyright © 2022 TW Group 13C, Weave BV, TU Delft
*/

package callanalyzer

import (
	"fmt"
	"go/constant"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/ssa"
)

// resolveParameter resolves a parameter in a frame
func resolveParameter(par *ssa.Parameter, fr *Frame) *ssa.Value {
	if fr == nil {
		return nil
	}

	if parameterValue, ok := fr.params[par]; ok {
		return parameterValue
	}

	for param, val := range fr.params {
		if param.Name() == par.Name() {
			return val
		}
	}

	return nil
}

// resolveRequestObject attempts to resolve a value for a request object
// TODO this implementation is very naive and will not cover many cases
func resolveRequestObject(value *ssa.Value, fr *Frame, substConf SubstitutionConfig) (string, bool) {
	if value == nil {
		return "", false
	}

	switch val := (*value).(type) {
	case *ssa.UnOp:
		return resolveRequestObject(&val.X, fr, substConf)
	case *ssa.Alloc:
		//	TODO: implement
		break
	case *ssa.Extract:
		return resolveRequestObject(&val.Tuple, fr, substConf)
	case *ssa.FieldAddr:
		return resolveRequestObject(&val.X, fr, substConf)
	case *ssa.Parameter:
		// (recursively) resolve a parameter to a value and return that value, if it is defined
		parameterValue := resolveParameter(val, fr.parent)

		if parameterValue != nil {
			return resolveRequestObject(parameterValue, fr.parent, substConf)
		}
	case *ssa.Call:
		call := val.Common()
		if call == nil || call.IsInvoke() {
			break
		}

		fn := call.StaticCallee()
		if fn == nil {
			break
		}

		signature, _ := getFunctionQualifiers(fn)

		if signature == "net/http.NewRequest" {
			return resolveRequestObject(&call.Args[1], fr, substConf)
		}

		if signature == "(*net/http.Request).WithContext" {
			return resolveRequestObject(&call.Args[0], fr, substConf)
		}

		if signature == "net/http.NewRequestWithContext" {
			return resolveRequestObject(&call.Args[2], fr, substConf)
		}

	// TODO: check for other calls
	// if packageName == "net/http" {
	// 	fmt.Println(signature)
	// }
	case *ssa.BinOp:
		return resolveValue(value, fr, substConf)
	case *ssa.Const:
		return resolveValue(value, fr, substConf)
	}
	return "", false
}

// IsRequestObjectType determines if the given type is either a pointer to, or a request object
func IsRequestObjectType(varType types.Type) bool {
	switch subType := varType.(type) {
	case *types.Pointer:
		return IsRequestObjectType(subType.Elem())
	case *types.Named:
		nameType := subType.String()
		return nameType == "net/http.Request"
	}

	return false
}

// resolveValue Resolves a supplied ssa.Value, only in the cases that are supported by the tool:
// - string concatenation (see BinOp),
// - string literal
// - call to os.GetEnv
// - other InterestingCalls with the action Substitute.
// It also returns a bool which indicates whether the variable was resolved.
func resolveValue(value *ssa.Value, fr *Frame, substConf SubstitutionConfig) (string, bool) {
	if value == nil {
		return "unknown: the give value is null", false
	}

	// if the value is a request object type, attempt to resolve it in some cases
	typ := (*value).Type()
	isRequestObject := IsRequestObjectType(typ)

	if isRequestObject {
		return resolveRequestObject(value, fr, substConf)
	}

	switch val := (*value).(type) {
	case *ssa.Parameter:
		// resolve a parameter to a value and return that value, if it is defined
		parameterValue := resolveParameter(val, fr.parent)

		if parameterValue != nil {
			return resolveValue(parameterValue, fr.parent, substConf)
		}

		return "unknown: the parameter was not resolved", false
	case *ssa.Global:
		if globalValue, ok := fr.globals[val]; ok {
			return resolveValue(globalValue, fr, substConf)
		}

		return "unknown: the global was not resolved", false

	case *ssa.UnOp:
		return resolveValue(&val.X, fr, substConf)

	case *ssa.BinOp:
		switch val.Op { //nolint:exhaustive
		case token.ADD:
			left, isLeftResolved := resolveValue(&val.X, fr, substConf)
			right, isRightResolved := resolveValue(&val.Y, fr, substConf)
			if isRightResolved && isLeftResolved {
				return left + right, true
			}

			return left + right, false
		default:
			return "unknown: only ADD binary operation is supported", false
		}
	case *ssa.Const:
		switch val.Value.Kind() { //nolint:exhaustive
		case constant.String:
			return constant.StringVal(val.Value), true
		default:
			return "unknown: not a string constant", false
		}
	case *ssa.Call:
		return handleSubstitutableCall(val, substConf)
	case *ssa.FieldAddr:
		return resolveValue(&val.X, fr, substConf)
	case *ssa.Alloc:
		// TODO: implement
		return "unknown: allocation not found", false
	default:
		return "unknown: the parameter was not resolved", false
	}
}

// handleSubstitutableCall handles substitution for calls that can't be easily resolved
// for example `os.getEnv()`
func handleSubstitutableCall(val *ssa.Call, substConf SubstitutionConfig) (string, bool) {
	unknownCallError := "unknown: substitutable call that is not supported"
	switch fnCallType := val.Call.Value.(type) {
	case *ssa.Function:
		{
			qualifiedFunctionNameOfTarget := fnCallType.RelString(nil)
			if substConf.substitutionCalls[qualifiedFunctionNameOfTarget].action == Substitute {
				switch argOfReplaceableCall := val.Call.Args[0].(type) {
				case *ssa.Const:
					{
						envVarName := constant.StringVal(argOfReplaceableCall.Value)
						envVarVal, ok := substConf.serviceEnv[envVarName]
						if ok {
							return envVarVal, true
						} else {
							return fmt.Sprintf("unknown: environment variable %s not set", envVarName), false
						}
					}
				default:
					return unknownCallError, false
				}
			}
		}
	default:
		return unknownCallError, false
	}
	return unknownCallError, false
}

// resolveParameters iterates over the parameters, resolving those where possible.
// It also keeps track of whether all variables could be resolved or not.
func resolveParameters(parameters []ssa.Value, positions []int, fr *Frame, serviceEnv SubstitutionConfig) ([]string, bool) {
	stringParameters := make([]string, len(positions))
	wasResolved := true

	for i, idx := range positions {
		if idx < len(parameters) {
			variable, isResolved := resolveValue(&parameters[idx], fr, serviceEnv)
			if isResolved {
				stringParameters[i] = variable
			} else {
				wasResolved = false
			}
		}
	}

	return stringParameters, wasResolved
}

// resolveGinAddrSlice is a hardcoded solution to resolve the port address of a Run command from the "github.com/gin-gonic/gin" library.
// Returns list of strings that represent the slice, and bool value indicating whether the variable was resolved.
// TODO: implement a general way for resolving variables in slices
func resolveGinAddrSlice(value ssa.Value) ([]string, bool) {
	unresolvedType := []string{fmt.Sprintf("unknown: var(%s) = ??", value.Name())}
	switch val := value.(type) {
	case *ssa.Slice:
		switch val1Type := val.X.(type) {
		case *ssa.Alloc:
			block := val1Type.Block()
			for i := range block.Instrs {
				// Iterate through instruction of the block in reverse order
				// In the case of the Gin library, the last Store instruction in the block contains the address value we're looking for
				switch instruction := block.Instrs[len(block.Instrs)-1-i].(type) {
				case *ssa.Store:
					switch storeValType := instruction.Val.(type) {
					case *ssa.Const:
						switch storeValType.Value.Kind() { //nolint:exhaustive
						case constant.String:
							return []string{constant.StringVal(storeValType.Value)}, true
						default:
							return unresolvedType, false
						}
					default:
						return unresolvedType, false
					}
				default:
					return unresolvedType, false
				}
			}
		default:
			return unresolvedType, false
		}
	case *ssa.Const:
		return []string{":8080"}, true
	}
	return unresolvedType, false
}
