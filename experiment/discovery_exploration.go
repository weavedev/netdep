package main

// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"fmt"
	"go/build"
	"go/types"
	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/discovery/CallAnalyzer"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/interp"
	"golang.org/x/tools/go/ssa/ssautil"
)

// flags
var (
	mode       = ssa.BuilderMode(0)
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

	config := &packages.Config{
		//Dir:
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

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(path)
	initial, err := packages.Load(config, "./test/sample/http/wrapped_call/wrapped_call.go")

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

	for i, p := range pkgs {
		if p == nil {
			return fmt.Errorf("cannot build SSA for package %s", initial[i])
		}
	}

	if !shouldRun {
		// Build and display only the initial packages
		// (and synthetic wrappers).
		for _, p := range pkgs {
			p.Build()
			CallAnalyzer.AnalyzePackage(p)
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
