package preprocessing

import (
	"go/ast"
	"go/parser"
	"go/token"
	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/discovery/callanalyzer"
	"strconv"
	"strings"
)

// parseComments parses the given file with a parser.ParseComments mode, filters out
// the comments which don't contain a substring "netdep:client" or "netdep:endpoint", generates an Annotation for
// every remaining comment and returns a list of them.
func parseComments(path string, serviceName string, annotations map[string]map[callanalyzer.Position]string) error {
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, path, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	for _, commentGroup := range f.Comments {
		for _, comment := range commentGroup.List {
			if strings.HasPrefix(comment.Text, "//netdep:") {
				tokenPos := fs.Position(comment.Slash)
				pos := callanalyzer.Position{
					Filename: tokenPos.Filename,
					Line:     tokenPos.Line,
				}
				value := strings.Join(strings.Split(comment.Text, "netdep:")[1:], "")

				annotations[serviceName][pos] = value
			}
		}
	}

	return nil
}

// parseInterfaces parses the given file, finds and collects all methods defined in interfaces
func parseInterfaces(path string, serviceName string, serviceCalls map[IntCall]string, serverTargets *[]*callanalyzer.CallTarget) {
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

// parseMethods parses the given file, finds and collects all interesting methods (interesting = found in the servicecalls package).
func parseMethods(path string, calls map[IntCall]string, serviceName string) (*[]*callanalyzer.CallTarget, error) {
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	clientTargets := make([]*callanalyzer.CallTarget, 0)

	ast.Inspect(f, func(n ast.Node) bool {
		// Find Function Call Statements
		funcCall, ok := n.(*ast.CallExpr)
		if ok {
			mthd, ok := funcCall.Fun.(*ast.SelectorExpr)
			if ok {
				sel, ok := mthd.X.(*ast.SelectorExpr)
				if !ok || (ok && !strings.HasSuffix(sel.Sel.Name, "DB")) {
					intCall := IntCall{
						Name:      mthd.Sel.Name,
						NumParams: len(funcCall.Args),
					}
					_, ex := calls[intCall]
					if ex {
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

						clientTargets = append(clientTargets, clientTarget)
					}
				}
			}
		}
		return true
	})

	return &clientTargets, nil
}
