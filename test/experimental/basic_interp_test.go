// This is a heavily simplified version of the file found at https://cs.opensource.google/go/x/tools/+/refs/tags/v0.1.10:go/ssa/interp/interp_test.go
// Team 13C IS NOT the original author of this file. For testing purposes only, with absolutely no warranty.
// IMPORTANT: The following changes were made to the original file:
// - Hardcoded paths for testing
// - Replaced GOOS from Linux to Win
// - Use imported package instead of created one
// - Hardcode actual goroot
// - Rename package
// - Remove TestTestdataFiles method

// Package experimental_test Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package experimental_test

// This test runs the SSA interpreter over sample Go programs.
// Because the interpreter requires intrinsics for assembly
// functions and many low-level runtime routines, it is inherently
// not robust to evolutionary change in the standard library.
// Therefore the test cases are restricted to programs that
// use a fake standard library in testdata/src containing a tiny
// subset of simple functions useful for writing assertions.
//
// We no longer attempt to interpret any real standard packages such as
// fmt or testing, as it proved too fragile.

import (
	"bytes"
	"fmt"
	"go/build"
	"go/types"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"golang.org/x/tools/go/loader"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/interp"
	"golang.org/x/tools/go/ssa/ssautil"
)

// Specific GOARCH to use for a test case in go.tools/go/ssa/interp/testdata/.
// Defaults to amd64 otherwise.
var testdataArchs = map[string]string{
	"width32.go": "386",
}

func run(t *testing.T, input string) bool {
	// The recover2 test case is broken on Go 1.14+. See golang/go#34089.
	// TODO(matloob): Fix this.
	if filepath.Base(input) == "recover2.go" {
		t.Skip("The recover2.go test is broken in go1.14+. See golang.org/issue/34089.")
	}

	t.Logf("Input: %s\n", input)

	start := time.Now()

	ctx := build.Default // copy
	//ctx.GOROOT = "C:\\Program Files\\Go" // fake goroot
	ctx.GOOS = "windows"
	ctx.GOARCH = "amd64"
	if arch, ok := testdataArchs[filepath.Base(input)]; ok {
		ctx.GOARCH = arch
	}

	conf := loader.Config{Build: &ctx}
	if _, err := conf.FromArgs([]string{input}, true); err != nil {
		t.Errorf("FromArgs(%s) failed: %s", input, err)
		return false
	}

	conf.Import("runtime")

	// Print a helpful hint if we don't make it to the end.
	var hint string
	defer func() {
		if hint != "" {
			fmt.Println("FAIL")
			fmt.Println(hint)
		} else {
			fmt.Println("PASS")
		}

		interp.CapturedOutput = nil
	}()
	hint = fmt.Sprintf("To dump SSA representation, run:\n%% go build golang.org/x/tools/cmd/ssadump && ./ssadump -test -build=CFP %s\n", input)

	iprog, err := conf.Load()
	if err != nil {
		t.Errorf("conf.Load(%s) failed: %s", input, err)
		return false
	}

	bmode := ssa.SanityCheckFunctions
	// bmode |= ssa.PrintFunctions // enable for debugging
	prog := ssautil.CreateProgram(iprog, bmode)
	prog.Build()
	var mainPkg *ssa.Package
	for _, info := range iprog.Imported {
		if strings.Contains(info.Pkg.Name(), "main") {
			mainPkg = prog.Package(info.Pkg)
			break
		}
	}

	//mainPkg = prog.Package(iprog.Created[0].Pkg)
	if mainPkg == nil {
		t.Fatalf("not a main package: %s", input)
	}

	interp.CapturedOutput = new(bytes.Buffer)

	sizes := types.SizesFor("gc", ctx.GOARCH)
	hint = fmt.Sprintf("To trace execution, run:\n%% go build golang.org/x/tools/cmd/ssadump && ./ssadump -build=C -test -run --interp=T %s\n", input)
	var imode interp.Mode // default mode
	// imode |= interp.DisableRecover // enable for debugging
	imode |= interp.EnableTracing // enable for debugging
	exitCode := interp.Interpret(mainPkg, imode, sizes, input, []string{})
	if exitCode != 0 {
		t.Fatalf("interpreting %s: exit code was %d", input, exitCode)
	}
	// $GOROOT/test tests use this convention:
	if strings.Contains(interp.CapturedOutput.String(), "BUG") {
		t.Fatalf("interpreting %s: exited zero but output contained 'BUG'", input)
	}

	hint = "" // call off the hounds

	if false {
		t.Log(input, time.Since(start)) // test profiling
	}

	return true
}

func printFailures(failures []string) {
	if failures != nil {
		fmt.Println("The following tests failed:")
		for _, f := range failures {
			fmt.Printf("\t%s\n", f)
		}
	}
}

// TestInterpBasic runs the interpreter on the hardcoded file
func TestInterpBasic(t *testing.T) {
	var failures []string
	currentDir, _ := os.Getwd()
	input := currentDir + "/../sample/uri_discovery/basic_concat"
	if !run(t, input) {
		failures = append(failures, input)
	}
	printFailures(failures)
}
