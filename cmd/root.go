/*
Package cmd contains all the application command definitions
Copyright © 2022 TW Group 13C, Weave BV, TU Delft
*/
package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/fatih/color"

	"github.com/spf13/cobra"

	"lab.weave.nl/internships/tud-2022/netDep/stages/discovery"
	"lab.weave.nl/internships/tud-2022/netDep/stages/discovery/callanalyzer"
	"lab.weave.nl/internships/tud-2022/netDep/stages/discovery/natsanalyzer"
	"lab.weave.nl/internships/tud-2022/netDep/stages/discovery/servicecallsanalyzer"
	"lab.weave.nl/internships/tud-2022/netDep/stages/matching"
	"lab.weave.nl/internships/tud-2022/netDep/stages/output"
	"lab.weave.nl/internships/tud-2022/netDep/stages/preprocessing"
	"lab.weave.nl/internships/tud-2022/netDep/structures"
)

// RunConfig defines the parameters for a depScan command run
type RunConfig struct {
	ProjectDir      string
	ServiceDir      string
	EnvFile         string
	Verbose         bool
	ServiceCallsDir string
	Shallow         bool
}

// RootCmd creates and returns a depScan command object
func RootCmd() *cobra.Command {
	// Variables that are supplied as command-line args
	var (
		projectDir      string
		serviceDir      string
		envVars         string
		outputFilename  string
		verbose         bool
		serviceCallsDir string
		shallow         bool
		noColor         bool
	)

	cmd := &cobra.Command{
		Use:   "netDep",
		Short: "Scan and report dependencies between microservices",
		Long: `Outputs network-communication-based dependencies of services within a microservice architecture Golang project.
Output is an adjacency list of service dependencies in a JSON format`,

		RunE: func(cmd *cobra.Command, args []string) error {
			color.NoColor = noColor // colourful terminal output

			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			projectDir = ensureAbsolutePath(cwd, projectDir)
			serviceDir = ensureAbsolutePath(cwd, serviceDir)

			ok, err := areInputPathsValid(projectDir, serviceDir, serviceCallsDir, envVars, outputFilename)
			if !ok {
				return err
			}

			config := RunConfig{
				ProjectDir:      projectDir,
				ServiceDir:      serviceDir,
				Verbose:         verbose,
				EnvFile:         envVars,
				ServiceCallsDir: serviceCallsDir,
				Shallow:         shallow,
			}

			// CALL OUR MAIN FUNCTIONALITY LOGIC FROM HERE AND SUPPLY BOTH PROJECT DIR AND SERVICE DIR
			dependencies, err := discoverAllCalls(config)
			if err != nil {
				return err
			}

			// generate output
			graph := matching.CreateDependencyGraph(dependencies)
			adjacencyList := output.ConstructAdjacencyList(graph)
			jsonString, err := output.SerializeAdjacencyList(adjacencyList, true)
			if err != nil {
				return err
			}

			allServices, err := preprocessing.FindServices(config.ServiceDir)
			if err != nil {
				return err
			}
			noReferenceToServices, noReferenceToAndFromServices := output.ConstructUnusedServicesLists(graph.Nodes, allServices)

			err = printOutput(outputFilename, jsonString, noReferenceToServices, noReferenceToAndFromServices)
			if err != nil {
				return err
			}

			return nil
		},
	}
	cmd.Flags().StringVarP(&projectDir, "project-directory", "p", "./", "project directory")
	cmd.Flags().StringVarP(&serviceDir, "service-directory", "s", "./svc", "service directory")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "toggle logging trace of unknown variables")
	cmd.Flags().StringVarP(&envVars, "environment-variables", "e", "", "environment variable file")
	cmd.Flags().StringVarP(&outputFilename, "output-filename", "o", "", "output filename such as ./deps.json")
	cmd.Flags().StringVarP(&serviceCallsDir, "servicecalls-directory", "c", "", "servicecalls package directory")
	cmd.Flags().BoolVarP(&noColor, "no-color", "n", false, "disable colourful terminal output")
	cmd.Flags().BoolVarP(&shallow, "shallow", "S", false, "toggle shallow scanning")
	return cmd
}

// ensureAbsolutePath makes sure the given path is absolute, or makes it absolute based on the current working directory
func ensureAbsolutePath(cwd, pth string) string {
	if !filepath.IsAbs(pth) {
		return filepath.Join(cwd, pth)
	}

	return pth
}

// printOutput writes the output to the target file (btw stdout is also a file on UNIX)
func printOutput(targetFileName, jsonString string, noReferenceToServices []string, noReferenceToAndFromServices []string) error {
	if targetFileName != "" {
		const filePerm = 0o600
		err := os.WriteFile(targetFileName, []byte(jsonString), filePerm)
		if err == nil {
			color.HiGreen("Successfully analysed, the dependencies have been output to %v\n", targetFileName)
		} else {
			color.Yellow("Could not write to file %s", targetFileName)
			color.HiGreen("Successfully analysed, here is the list of dependencies:")
			color.HiWhite(jsonString)
			output.PrintUnusedServices(noReferenceToServices, noReferenceToAndFromServices)
			return err
		}
	} else {
		color.HiGreen("Successfully analysed, here is the list of dependencies:")
		color.HiWhite(jsonString)
		output.PrintUnusedServices(noReferenceToServices, noReferenceToAndFromServices)
	}
	return nil
}

// areInputPathsValid verifies that all the specified directories exist before running the main logic
func areInputPathsValid(projectDir, serviceDir, serviceCallsDir, envVars, outputFilename string) (bool, error) {
	if !pathOk(projectDir) {
		return false, fmt.Errorf("invalid project directory specified: %s", projectDir)
	}

	if !pathOk(serviceDir) {
		return false, fmt.Errorf("invalid service directory: %s", serviceDir)
	}

	if !pathOk(serviceCallsDir) && serviceCallsDir != "" {
		return false, fmt.Errorf("invalid servicecalls directory specified: %s", serviceCallsDir)
	}

	if !pathOk(envVars) && envVars != "" {
		return false, fmt.Errorf("invalid environment variable file specified: %s", envVars)
	}
	jsonParentDir := path.Dir(outputFilename)

	if !pathOk(jsonParentDir) {
		return false, fmt.Errorf("parent directory of json path does not exist: %s", jsonParentDir)
	}

	return true, nil
}

// pathOk checks whether the specified directory dir exists
func pathOk(dir string) bool {
	if ex, err := pathExists(dir); !ex || err != nil {
		return false
	}
	return true
}

// init initialises the depScan command and adds it as a subcommand of the root
// func init() {
// netDepCmd := netDepCmd()
// rootCmd.AddCommand(netDepCmd)
// }

// pathExists determines whether a given path is valid by checking if it exists
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// resolveEnvironmentValues calls resolving stage if the path is not unspecified(""), returns nil otherwise
func resolveEnvironmentValues(path string) (map[string]map[string]string, error) {
	if path == "" {
		return nil, nil
	}
	return preprocessing.IndexEnvironmentVariables(path)
}

// discoverAllCalls calls the correct stages for loading, building,
// filtering and discovering all client and server calls.
func discoverAllCalls(config RunConfig) (*structures.Dependencies, error) {
	// Filtering
	services, err := preprocessing.FindServices(config.ServiceDir)
	fmt.Printf("Starting to analyse %d services.\n", len(services))

	if err != nil {
		return nil, err
	}

	// resolve environment values
	// TODO: Integrate the envVariables into discovery
	envVariables, err := resolveEnvironmentValues(config.EnvFile)
	if err != nil {
		return nil, err
	}

	analyserConfig := callanalyzer.DefaultConfigForFindingHTTPCalls()
	analyserConfig.SetVerbose(config.Verbose)
	analyserConfig.SetEnv(envVariables)

	allClientTargets, allServerTargets, annotations, err := processEachService(&services, &config, &analyserConfig)
	if err != nil {
		return nil, err
	}

	consumers, producers, err := natsanalyzer.FindNATSCalls(config.ServiceDir)
	if err != nil {
		return nil, err
	}

	if config.Verbose {
		output.PrintDiscoveredAnnotations(annotations)
	}

	dependencies := &structures.Dependencies{
		Calls:     allClientTargets,
		Endpoints: allServerTargets,
		Consumers: consumers,
		Producers: producers,
	}

	return dependencies, err
}

// processEachService preprocesses and analyses each of the services using RunConfig and callanalyzer.AnalyserConfig
func processEachService(services *[]string, config *RunConfig, analyserConfig *callanalyzer.AnalyserConfig) ([]*callanalyzer.CallTarget, []*callanalyzer.CallTarget, map[string]map[callanalyzer.Position]string, error) {
	allClientTargets := make([]*callanalyzer.CallTarget, 0)
	allServerTargets := make([]*callanalyzer.CallTarget, 0)
	annotations := make(map[string]map[callanalyzer.Position]string)

	analyserConfig.SetAnnotations(annotations)

	packageCount := 0

	internalCalls, serverTargets, err := servicecallsanalyzer.ParseServiceCallsPackage(config.ServiceCallsDir)
	if err != nil {
		return nil, nil, nil, err
	}

	allServerTargets = append(allServerTargets, *serverTargets...)
	internalClientTargets := make([]*callanalyzer.CallTarget, 0)
	noneInspected := true

	for _, serviceDir := range *services {
		serviceName := strings.Split(serviceDir, string(os.PathSeparator))[len(strings.Split(serviceDir, string(os.PathSeparator)))-1]

		if config.Verbose {
			fmt.Printf("Analysing service %s\n", serviceDir)
		}

		err := preprocessing.LoadAnnotations(serviceDir, serviceName, annotations)
		if err != nil {
			handleErr(fmt.Sprintf("Error while loading annotations of %s", serviceName), config.Verbose, err)
			continue
		}

		// There are some interesting internal calls so the tool should parse all methods
		if len(internalCalls) != 0 {
			err = servicecallsanalyzer.LoadServiceCalls(serviceDir, serviceName, internalCalls, &internalClientTargets)
			if err != nil {
				handleErr(fmt.Sprintf("Error while loading service calls of %s", serviceName), config.Verbose, err)
				continue
			}
		}

		// Load and build packages and proceed with discovery if the user
		// Didn't ask for shallow scanning
		if !config.Shallow {
			// load packages
			packagesInService, err := preprocessing.LoadAndBuildPackages(config.ProjectDir, serviceDir)
			if err != nil {
				handleErr(fmt.Sprintf("Error while loading packages of %s", serviceName), config.Verbose, err)
				continue
			}
			packageCount += len(packagesInService)

			// discover calls
			clientCalls, serverCalls, err := discovery.DiscoverAll(packagesInService, analyserConfig)
			if err != nil {
				handleErr(fmt.Sprintf("Error while trying to discover network calls of %s", serviceName), config.Verbose, err)
				continue
			}

			if config.Verbose {
				clientSum := len(allClientTargets)
				targetSum := len(clientCalls)

				color.Green("Found %d calls of which %d client call(s) and %d server call(s)", clientSum+targetSum, clientSum, targetSum)
			}

			// append
			allClientTargets = append(allClientTargets, clientCalls...)
			allServerTargets = append(allServerTargets, serverCalls...)
			noneInspected = false
		}
	}

	allClientTargets = append(allClientTargets, internalClientTargets...)

	if (!config.Shallow && packageCount == 0) || noneInspected {
		return nil, nil, nil, fmt.Errorf("found no services to analyse")
	}
	return allClientTargets, allServerTargets, annotations, nil
}

func handleErr(sprintf string, verbose bool, err error) {
	color.HiYellow(sprintf)
	if verbose {
		color.HiYellow("%s", err)
	}
}
