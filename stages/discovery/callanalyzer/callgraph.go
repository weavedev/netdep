package callanalyzer

import (
	"fmt"
	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/pointer"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
	"os"
	"strings"
)

func getCallTarget(trace []*ssa.CallCommon, pkg *ssa.Package) *CallTarget {
	rootCall := trace[0]
	rootFn := rootCall.StaticCallee()
	MethodName, packageName := getFunctionQualifiers(rootFn)

	target := CallTarget{
		packageName:     packageName,
		MethodName:      MethodName,
		ServiceName:     pkg.String()[strings.LastIndex(pkg.String(), "/")+1:],
		RequestLocation: "",
		IsResolved:      false,
		Trace:           nil,
	}

	// add trace
	for _, tracedCall := range trace {
		if tracedCall == nil {
			continue
		}

		filePath, position := getPositionFromPos(tracedCall.Pos(), pkg.Prog)

		newTrace := CallTargetTrace{
			// split package name and take the last item to get the service name
			FileName:       filePath[strings.LastIndex(filePath, string(os.PathSeparator))+1:],
			PositionInFile: position,
		}

		target.Trace = append(target.Trace, newTrace)
	}

	return &target
}

func AnalyzeUsingCallGraph(pkgs []*ssa.Package, config *AnalyserConfig) ([]*CallTarget, []*CallTarget, error) {
	mains := ssautil.MainPackages(pkgs)

	ptConfig := &pointer.Config{
		Mains:          mains,
		BuildCallGraph: true,
	}

	// query

	fmt.Println("Running pointer analysis...")
	ptares, err := pointer.Analyze(ptConfig)

	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	cg := ptares.CallGraph
	cg.DeleteSyntheticNodes()

	allClientTargets := make([]*CallTarget, 0)
	allServerTargets := make([]*CallTarget, 0)

	fmt.Printf("Finding nodes (%d)...\n", len(cg.Nodes))

	markedPaths := make(map[int]bool)
	markPaths(cg.Root, markedPaths)
	fmt.Printf("Marked %d nodes.\n", len(markedPaths))

	for _, node := range cg.Nodes {
		qualifiedFunctionNameOfTarget := node.Func.RelString(nil)

		_, isInterestingClient := config.interestingCallsClient[qualifiedFunctionNameOfTarget]
		_, isInterestingServer := config.interestingCallsServer[qualifiedFunctionNameOfTarget]

		if !isInterestingClient && !isInterestingServer {
			continue
		}
		fmt.Println("found: " + qualifiedFunctionNameOfTarget)
		traces := getUniquePaths(node, cg.Root, markedPaths)
		fmt.Println(len(traces))

		//callgraph.PathSearch()

		for _, trace := range traces {
			callRoot := trace[len(trace)-1]
			call := getCallTarget(trace, callRoot.StaticCallee().Pkg)

			if call == nil {
				continue
			}

			if isInterestingClient {
				allClientTargets = append(allClientTargets, call)
			} else {
				allServerTargets = append(allServerTargets, call)
			}
		}
	}

	return allClientTargets, allServerTargets, nil
}

func markPaths(start *callgraph.Node, visited map[int]bool) {
	visited[start.ID] = true

	for _, edge := range start.Out {
		if visited[edge.Callee.ID] {
			continue
		}

		markPaths(edge.Callee, visited)
	}

	return
}

func getUniquePaths(start *callgraph.Node, end *callgraph.Node, marked map[int]bool) [][]*ssa.CallCommon {
	visited := make(map[int]bool)
	visited[start.ID] = true

	foundPaths := uniquePaths(start, end, visited, marked)
	return foundPaths
}

func uniquePaths(node *callgraph.Node, root *callgraph.Node, visited map[int]bool, marked map[int]bool) [][]*ssa.CallCommon {
	output := make([][]*ssa.CallCommon, 0)

	for _, edge := range node.In {
		inNode := edge.Caller
		if visited[inNode.ID] || !marked[inNode.ID] {
			continue
		}

		var call *ssa.CallCommon = nil
		if edge.Site != nil {
			call = edge.Site.Common()
		}

		if inNode.ID == root.ID {
			newList := make([]*ssa.CallCommon, 0, 1)
			output = append(output, newList)
			continue
		}

		newVisited := make(map[int]bool, len(visited)+1)
		for k, _ := range visited {
			newVisited[k] = true
		}

		newVisited[inNode.ID] = true

		newPaths := uniquePaths(inNode, root, newVisited, marked)

		if newPaths == nil {
			continue
		}

		for _, pth := range newPaths {
			newList := make([]*ssa.CallCommon, 0)
			newList = append(newList, call)
			newList = append(newList, pth...)
			output = append(output, newList)
		}
	}

	if len(output) == 0 {
		return nil
	}

	return output
}
