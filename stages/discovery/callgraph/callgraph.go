package callgraph

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/callgraph/static"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

// ExtractMainPackages looks for packages with a name "main"
// and a function that is called main. This will be used to find
// entry points of execution
func ExtractMainPackages(pkgs []*ssa.Package) ([]*ssa.Package, error) {
	var mains []*ssa.Package
	for _, p := range pkgs {
		if p != nil && p.Pkg.Name() == "main" && p.Func("main") != nil {
			fmt.Println(p.Func("main"))
			mains = append(mains, p)
		}
	}

	if len(mains) == 0 {
		return nil, fmt.Errorf("no main packages")
	}

	return mains, nil
}

func DoCallGraph(dir string, pkgsArr []string) (*callgraph.Graph, error) {
	var cg *callgraph.Graph
	//fmt.Println(pkgsArr[0])
	cfg := &packages.Config{
		Dir:  "",
		Mode: packages.LoadAllSyntax,
	}
	//fmt.Println(cfg.Env, os.Environ())

	//return nil, nil
	initial, err := packages.Load(cfg, pkgsArr...)
	if err != nil {
		return cg, err
	}

	//if packages.PrintErrors(initial) > 0 {
	//	return cg, fmt.Errorf("Packages contain errors. Please resolve them before running again")
	//}

	prog, discarded := ssautil.AllPackages(initial, 0)
	for _, prog1 := range prog.AllPackages() {
		//fmt.Println(prog1.Pkg.Name())
		if prog1.Pkg.Name() == "stages" {
			//callgraph.GraphVisitEdges(prog1.Members,visitCallback)
			//for _, member := range prog1.Members {
			//member
			//}

			for _, member := range prog1.Members {
				switch memT := member.(type) {
				case *ssa.Function:
					switch syntaxT := memT.Syntax().(type) {
					case ast.Node:
						ast.Inspect(syntaxT, astInspectHandler)
					default:
						continue
					}
					// The following was an attempt at printing all calls manually. The above line is an attempt to do it automatically.
					//for _, stmt := range memT.Syntax().(*ast.FuncDecl).Body.List { //ConstructOutput
					//	switch stmtT := stmt.(type) {
					//	case *ast.AssignStmt:
					//		for _, rh := range stmtT.Rhs {
					//			if rh.(*ast.CallExpr).Fun.(*ast.SelectorExpr).X.(*ast.Ident).Name == "http" {
					//				fmt.Println("Found a call to an http-related method!")
					//			}
					//		}
					//	}
					//}
				default:
					//Not a function
					continue
				}
				//member.(*ssa.Function).
			}

			//Creating a callgraph instance seems to be unnecessary at this stage; we look at syntax anyway
			//var cgFun = callgraph.New(prog1.Members["ConstructOutput"].(*ssa.Function))
			//fmt.Println(cgFun)
		}
	}
	var _, _ = ExtractMainPackages(discarded)
	prog.Build()
	var gr = static.CallGraph(prog)
	//print(discarded)
	//print(gr)
	return gr, nil
}

func astInspectHandler(n ast.Node) bool {
	//fmt.Printf("Callback was called with node %+v\n", n)
	switch T := n.(type) {
	case *ast.CallExpr:
		fmt.Printf("%v", T)
		//switch n.(*ast.CallExpr).Fun
		var serialized, _ = json.MarshalIndent(n.(*ast.CallExpr).Fun, "", "\t")
		var serialized2, _ = json.MarshalIndent(n.(*ast.CallExpr).Args, "", "\t")

		fmt.Printf("It was a callexpression that calls the function: %+v \n", string(serialized))
		fmt.Printf("With arguments: %+v\n", string(serialized2))

	}
	return true
}
