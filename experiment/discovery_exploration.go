package main

// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"fmt"
	"go/build"
	"go/types"
	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/pointer"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/interp"
	"golang.org/x/tools/go/ssa/ssautil"
	"lab.weave.nl/internships/tud-2022/static-analysis-project/experiment/analyze"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strings"
)

// flags
var (
	mode = ssa.BuilderMode(0)
	//rootDir    = "/../nid-core"
	//projectDir = "./svc/autopseudo/"
	rootDir    = "/"
	projectDir = "./test/sample/http/multiple_calls/"
	//rootDir    = "/test/example"
	//projectDir = "./svc/node-basic-http/"
	CpuProfile = ""
	args       = []string{}
	shouldRun  = false
)

func main() {
	if err := doMain(); err != nil {
		fmt.Fprintf(os.Stderr, "ssadump: %s\n", err)
		os.Exit(1)
	}
}

func getSizes() *types.StdSizes {

	var wordSize int64 = 8

	switch build.Default.GOARCH {
	case "386", "arm":
		wordSize = 4
	}

	sizes := &types.StdSizes{
		MaxAlign: 8,
		WordSize: wordSize,
	}

	return sizes
}

func doMain() error {

	path, cwdErr := os.Getwd()

	fmt.Println(path)

	if cwdErr != nil {
		log.Println(cwdErr)
	}

	config := &packages.Config{
		Dir: filepath.Clean(path + rootDir),
		//Mode:  packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedImports | packages.NeedTypes | packages.NeedTypesSizes | packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedDeps,
		Mode:  packages.LoadAllSyntax,
		Tests: false,
	}

	sizes := getSizes()

	var interpretMode interp.Mode
	//interpretMode |= interp.EnableTracing
	//interpretMode |= interp.DisableRecover

	// Profiling support.
	if CpuProfile != "" {
		f, err := os.Create(CpuProfile)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	initial, err := packages.Load(config, projectDir)

	if err != nil {
		return err
	}

	if len(initial) == 0 {
		return fmt.Errorf("no packages")
	}
	if packages.PrintErrors(initial) > 0 {
		return fmt.Errorf("packages contain errors")
	}

	// Create SSA-form program representation.
	prog, pkgs := ssautil.AllPackages(initial, mode)
	prog.Build()

	for i, p := range pkgs {
		if p == nil {
			return fmt.Errorf("cannot build SSA for package %s", initial[i])
		}
	}

	mains := ssautil.MainPackages(pkgs)

	ptConfig := &pointer.Config{
		Mains:          mains,
		BuildCallGraph: true,
	}

	ptares, err := pointer.Analyze(ptConfig)
	if err != nil {
		return err // internal error in pointer analysis
	}

	cg := ptares.CallGraph

	callgraph.GraphVisitEdges(cg, func(edge *callgraph.Edge) error {
		if edge.Callee.Func.Pkg != nil && edge.Callee.Func.Pkg.Pkg.Path() == "net/http" {
			funcName := edge.Callee.Func.RelString(nil)
			if strings.Contains(funcName, "Do") {
				//fmt.Println(edge)
			}
		}
		return nil
	})

	if !shouldRun {
		// Build and display only the initial packages
		// (and synthetic wrappers).
		for _, p := range pkgs {
			p.Build()
			analyze.AnalyzePackage(p)
		}
	} else {
		// Run the interpreter.
		// Build SSA for all packages.
		prog.Build()

		// The interpreter needs the runtime package.
		// It is a limitation of go/packages that
		// we cannot add "runtime" to its initial set,
		// we can only check that it is present.
		if prog.ImportedPackage("runtime") == nil {
			return fmt.Errorf("-run: program does not depend on runtime")
		}

		if runtime.GOARCH != build.Default.GOARCH {
			return fmt.Errorf("cross-interpretation is not supported (target has GOARCH %s, interpreter has %s)",
				build.Default.GOARCH, runtime.GOARCH)
		}

		// Run first main package.
		for _, main := range ssautil.MainPackages(pkgs) {
			fmt.Fprintf(os.Stderr, "Running: %s\n", main.Pkg.Path())
			interp.Interpret(main, interpretMode, sizes, main.Pkg.Path(), args)
			return nil
		}
		return fmt.Errorf("no main package")
	}
	return nil
}
