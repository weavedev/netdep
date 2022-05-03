package callgraph

import (
	"fmt"
	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/packages"
	"os"
)

func DoCallGraph(dir string, pkgsArr []string) (*callgraph.Graph, error) {
	//var cg *callgraph.Graph

	cfg := &packages.Config{
		Dir:  "/Users/martynaskrupskis/Documents/code/",
		Mode: packages.NeedFiles,
	}
	fmt.Println(cfg.Env, os.Environ())

	return nil, nil
	//initial, err := packages.Load(cfg, pkgsArr...)
	//if err != nil {
	//	return cg, err
	//}
	//
	//if packages.PrintErrors(initial) > 0 {
	//	return cg, fmt.Errorf("Packages contain errors. Please resolve them before running again")
	//}
	//
	//prog, _ := ssautil.AllPackages(initial, 0)
	//prog.Build()
	//
	//return static.CallGraph(prog), nil
}
