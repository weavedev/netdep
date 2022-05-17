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

// CreateSSA constructs a default SSA
func CreateSSA(ssaConf SSAConfig) (*ssa.Program, []*ssa.Package, error) {
	pkgConf := &packages.Config{
		Dir: ssaConf.ProjDir,
		// Mode:  packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedImports | packages.NeedTypes | packages.NeedTypesSizes | packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedDeps,
		//nolint
		Mode:  packages.LoadAllSyntax,
		Tests: false,
	}

	initial, err := packages.Load(pkgConf, ssaConf.SvcDir)

	if err != nil {
		return nil, nil, err
	}

	if len(initial) == 0 {
		return nil, nil, fmt.Errorf("no packages")
	}
	if packages.PrintErrors(initial) > 0 {
		return nil, nil, fmt.Errorf("packages contain errors")
	}

	// Create SSA-form program representation.
	prog, pkgs := ssautil.AllPackages(initial, ssaConf.Mode)
	prog.Build()
	return prog, pkgs, nil
}
