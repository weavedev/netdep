/*
Package cmd contains all the application command definitions
Copyright © 2022 TW Group 13C, Weave BV, TU Delft
*/
package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"

	"lab.weave.nl/internships/tud-2022/netDep/stages/discovery"
	"lab.weave.nl/internships/tud-2022/netDep/stages/discovery/callanalyzer"
	"lab.weave.nl/internships/tud-2022/netDep/stages/discovery/natsanalyzer"
	"lab.weave.nl/internships/tud-2022/netDep/stages/matching"
	"lab.weave.nl/internships/tud-2022/netDep/stages/output"
	"lab.weave.nl/internships/tud-2022/netDep/stages/preprocessing"
	"lab.weave.nl/internships/tud-2022/netDep/structures"
)

// RunConfig defines the parameters for a depScan command run
type RunConfig struct {
	ProjectDir string
	ServiceDir string
	EnvFile    string
	Verbose    bool
}

// RootCmd creates and returns a depScan command object
func RootCmd() *cobra.Command {
	// Variables that are supplied as command-line args
	var (
		projectDir     string
		serviceDir     string
		envVars        string
		outputFilename string
		verbose        bool
	)

	cmd := &cobra.Command{
		Use:   "netDep",
		Short: "Scan and report dependencies between microservices",
		Long: `Outputs network-communication-based dependencies of services within a microservice architecture Golang project.
Output is an adjacency list of service dependencies in a JSON format`,

		RunE: func(cmd *cobra.Command, args []string) error {
			ok, err := areInputPathsValid(projectDir, serviceDir, envVars, outputFilename)
			if !ok {
				return err
			}

			config := RunConfig{
				ProjectDir: projectDir,
				ServiceDir: serviceDir,
				Verbose:    verbose,
				EnvFile:    envVars,
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

			err = printOutput(outputFilename, jsonString)
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
	return cmd
}

// printOutput writes the output to the target file (btw stdout is also a file on UNIX)
func printOutput(targetFileName, jsonString string) error {
	if targetFileName != "" {
		const filePerm = 0o600
		err := os.WriteFile(targetFileName, []byte(jsonString), filePerm)
		if err == nil {
			fmt.Printf("Successfully analysed, the dependencies have been output to %v\n", targetFileName)
		} else {
			// Could not write to file, output to stdout
			fmt.Println(jsonString)
			return err
		}
	} else {
		fmt.Println("Successfully analysed, here is the list of dependencies:")
		fmt.Println(jsonString)
	}
	return nil
}

// areInputPathsValid verifies that all the specified directories exist before running the main logic
func areInputPathsValid(projectDir, serviceDir, envVars, outputFilename string) (bool, error) {
	if !pathOk(projectDir) {
		return false, fmt.Errorf("invalid project directory specified: %s", projectDir)
	}

	if !pathOk(serviceDir) {
		return false, fmt.Errorf("invalid service directory: %s", serviceDir)
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
	// Given a correct project directory en service directory,
	// apply our discovery algorithm to find all interesting calls
	if ex, err := pathExists(config.EnvFile); !ex && config.EnvFile != "" || err != nil {
		return nil, fmt.Errorf("invalid environment variable file specified: %s", config.EnvFile)
	}

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
func processEachService(services *[]string, config *RunConfig, analyserConfig *callanalyzer.AnalyserConfig) ([]*callanalyzer.CallTarget, []*callanalyzer.CallTarget, map[string]map[preprocessing.Position]string, error) {
	allClientTargets := make([]*callanalyzer.CallTarget, 0)
	allServerTargets := make([]*callanalyzer.CallTarget, 0)
	annotations := make(map[string]map[preprocessing.Position]string)

	analyserConfig.SetAnnotations(annotations)

	packageCount := 0

	for _, serviceDir := range *services {
		if config.Verbose {
			fmt.Println("Analysing service " + serviceDir)
		}

		// load packages
		packagesInService, err := preprocessing.LoadAndBuildPackages(config.ProjectDir, serviceDir)
		if err != nil {
			return nil, nil, nil, err
		}
		packageCount += len(packagesInService)

		serviceName := strings.Split(serviceDir, string(os.PathSeparator))[len(strings.Split(serviceDir, string(os.PathSeparator)))-1]
		err = preprocessing.LoadAnnotations(serviceDir, serviceName, annotations)
		if err != nil {
			return nil, nil, nil, err
		}

		// discover calls
		clientCalls, serverCalls, err := discovery.DiscoverAll(packagesInService, analyserConfig)
		if err != nil {
			return nil, nil, nil, err
		}

		// append
		allClientTargets = append(allClientTargets, clientCalls...)
		allServerTargets = append(allServerTargets, serverCalls...)
	}

	if packageCount == 0 {
		return nil, nil, nil, fmt.Errorf("no service to analyse were found")
	}
	return allClientTargets, allServerTargets, annotations, nil
}
