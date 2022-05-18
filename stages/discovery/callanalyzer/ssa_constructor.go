package callanalyzer

import (
	"fmt"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

type SSAConfig struct {
	Mode    ssa.BuilderMode
	ProjDir string
	SvcDir  string
}

// CreateSSA constructs an SSA using predefined SSA construction config.
func CreateSSA(ssaConf SSAConfig) (*ssa.Program, []*ssa.Package, error) {
	pkgConf := &packages.Config{
		// The project directory
		Dir: ssaConf.ProjDir,
		// Mode:  packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedImports | packages.NeedTypes | packages.NeedTypesSizes | packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedDeps,
		// Unfortunately, it seems that the LoadAllSyntax flag is the one we need, and it is equivalent to the disjunction of the above types ^
		// whose iota is 991 in the current library version
		//nolint
		Mode: packages.LoadAllSyntax,
		// Do not analyse test files
		Tests: false,
	}

	// Load packages of the service directory
	initial, err := packages.Load(pkgConf, ssaConf.SvcDir)
	if err != nil {
		return nil, nil, err
	}

	if len(initial) == 0 {
		return nil, nil, fmt.Errorf("no packages found")
	}

	if packages.PrintErrors(initial) > 0 {
		return nil, nil, fmt.Errorf("packages contain errors")
	}

	// Create an SSA representation of the program.
	prog, pkgs := ssautil.AllPackages(initial, ssaConf.Mode)
	prog.Build()
	return prog, pkgs, nil
}
