package callgraph

import (
	"fmt"
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
	fmt.Println(pkgsArr[0])
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
	var _, _ = ExtractMainPackages(discarded)
	prog.Build()
	var gr = static.CallGraph(prog)
	print(discarded)
	print(gr)
	return gr, nil
}
