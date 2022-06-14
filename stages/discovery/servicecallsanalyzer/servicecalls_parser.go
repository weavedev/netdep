/*
Package servicecallsanalyzer defines servicecalls package specific scanning methods
Copyright © 2022 TW Group 13C, Weave BV, TU Delft
*/
package servicecallsanalyzer

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"strings"

	"lab.weave.nl/internships/tud-2022/netDep/stages/discovery/callanalyzer"
)

// ParseInterfaces parses the given file, finds and collects all methods defined in interfaces
func ParseInterfaces(path string, serviceName string, serviceCalls map[IntCall]string, serverTargets *[]*callanalyzer.CallTarget) {
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, path, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	ast.Inspect(f, func(n ast.Node) bool {
		// Find Function Call Statements
		interf, ok := n.(*ast.InterfaceType)
		if ok {
			for _, method := range interf.Methods.List {
				funcType, ok := method.Type.(*ast.FuncType)
				if ok {
					numParams := len(funcType.Params.List)
					intCall := IntCall{
						Name:      method.Names[0].Name,
						NumParams: numParams,
					}
					serviceCalls[intCall] = serviceName

					target := &callanalyzer.CallTarget{
						PackageName:     "servicecalls",
						MethodName:      method.Names[0].Name,
						RequestLocation: method.Names[0].Name,
						IsResolved:      true,
						ServiceName:     serviceName,
						Trace:           nil,
					}

					*serverTargets = append(*serverTargets, target)
				}
			}
		}
		return true
	})
}

// ParseMethods parses the given file, finds and collects all interesting methods (interesting = found in the servicecalls package).
func ParseMethods(path string, calls map[IntCall]string, serviceName string) (*[]*callanalyzer.CallTarget, error) {
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	clientTargets := make([]*callanalyzer.CallTarget, 0)
	serviceCallsImport := false

	for _, im := range f.Imports {
		if strings.Contains(im.Path.Value, "servicecalls") {
			serviceCallsImport = true
		}
	}

	if !serviceCallsImport {
		return &clientTargets, nil
	}

	ast.Inspect(f, func(n ast.Node) bool {
		// Find Function Call Statements
		funcCall, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		mthd, ok := funcCall.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		// Retrieve the selector which appears before the function call
		// Ex: selector.FunctionCall()
		sel, ok1 := mthd.X.(*ast.SelectorExpr)

		// If there's a selector - ensure it doesn't end with DB, as that introduces
		// some false positives with very generic method names such as "Update", "Get", "Delete"
		if !ok1 || (ok1 && !strings.HasSuffix(sel.Sel.Name, "DB")) {
			checkInterestingCall(fs, mthd, funcCall, calls, serviceName, &clientTargets)
		}
		return true
	})
	return &clientTargets, nil
}

// checkInterestingCall checks whether the given call is one of the calls found in the serviceCalls package.
func checkInterestingCall(fs *token.FileSet, mthd *ast.SelectorExpr, funcCall *ast.CallExpr, calls map[IntCall]string, serviceName string, clientTargets *[]*callanalyzer.CallTarget) {
	intCall := IntCall{
		Name:      mthd.Sel.Name,
		NumParams: len(funcCall.Args),
	}
	_, ex := calls[intCall]
	// If the call is one of the servicecalls
	if ex {
		// Create a new client target
		clientTarget := &callanalyzer.CallTarget{
			PackageName:     "servicecalls",
			MethodName:      mthd.Sel.Name,
			RequestLocation: mthd.Sel.Name,
			IsResolved:      true,
			ServiceName:     serviceName,
			Trace: []callanalyzer.CallTargetTrace{
				{
					FileName:       fs.Position(mthd.Sel.NamePos).Filename,
					PositionInFile: strconv.Itoa(fs.Position(mthd.Sel.NamePos).Line),
				},
			},
		}

		*clientTargets = append(*clientTargets, clientTarget)
	}
}
