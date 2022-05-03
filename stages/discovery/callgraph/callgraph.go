package callgraph

import (
	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/callgraph/static"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa/ssautil"
	"os"
)

func DoCallGraph(dir string, pkgsArr []string) (*callgraph.Graph, error) {
	var cg *callgraph.Graph

	cfg := &packages.Config{
		Dir:  "",
		Mode: packages.LoadAllSyntax,
		Env:  append(os.Environ(), "GOOS=plan9", "GOARCH=386"),
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
	prog.Build()
	var gr = static.CallGraph(prog)
	print(discarded)
	print(gr)
	return gr, nil
}
