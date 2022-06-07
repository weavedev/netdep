package callanalyzer

import (
	"go/types"

	"golang.org/x/tools/go/ssa"
)

func getInvokedFunctionFromCall(call *ssa.CallCommon, frame *Frame) *ssa.Function {
	pkg := call.Method.Pkg()
	name := call.Method.Name()
	program := frame.pkg.Prog

	callValue := &call.Value

	// resolve if parameter
	if param, isParam := (*callValue).(*ssa.Parameter); isParam {
		parValue, _ := resolveParameter(param, frame)
		if parValue != nil {
			callValue = parValue
		}
	}

	// resolve type of value
	var methodType types.Type
	switch expr := (*callValue).(type) {
	case *ssa.MakeInterface:
		methodType = expr.X.Type()
	default:
		// methodType = expr.Type()
		methodType = call.Method.Type()
	}

	// lookup method in type
	methodSet := program.MethodSets.MethodSet(methodType)
	sel := methodSet.Lookup(pkg, name)

	if sel != nil {
		return program.MethodValue(sel)
	} else {
		return nil
	}
}

func getCallFunctionFromCall(call *ssa.CallCommon, frame *Frame) *ssa.Function {
	if param, isParam := call.Value.(*ssa.Parameter); isParam {
		parValue, _ := resolveParameter(param, frame)
		if paramFn, isFn := (*parValue).(*ssa.Function); isFn {
			// TODO: does this happen?
			return paramFn
		}
	}

	return call.StaticCallee()
}

func getFunctionFromCall(call *ssa.CallCommon, frame *Frame) *ssa.Function {
	if call.IsInvoke() {
		return getInvokedFunctionFromCall(call, frame)
	} else {
		return getCallFunctionFromCall(call, frame)
	}
}
