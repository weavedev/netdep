/*
Package cmd contains all the application command definitions
Copyright © 2022 TW Group 13C, Weave BV, TU Delft
*/
package cmd

import (
	"fmt"
	"os"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/matching"
	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/output"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/discovery"
	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/discovery/callanalyzer"

	"github.com/spf13/cobra"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages"
)

var (
	projectDir string
	serviceDir string
	verbose    bool
	envVars    string
)

// RunConfig defines the parameters for a depScan command run
type RunConfig struct {
	ProjectDir string
	ServiceDir string
	EnvFile    string
	Verbose    bool
}

// depScanCmd creates and returns a depScan command object
func depScanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "depScan",
		Short: "Scan and report dependencies between microservices",
		Long: `Outputs network-communication-based dependencies of services within a microservice architecture Golang project.
Output is an adjacency list of service dependencies in a JSON format`,

		RunE: func(cmd *cobra.Command, args []string) error {
			// Path validation
			if ex, err := pathExists(projectDir); !ex || err != nil {
				return fmt.Errorf("invalid project directory specified: %s", projectDir)
			}
			if ex, err := pathExists(serviceDir); !ex || err != nil {
				return fmt.Errorf("invalid service directory specified: %s", serviceDir)
			}

			config := RunConfig{
				ProjectDir: projectDir,
				ServiceDir: serviceDir,
				Verbose:    verbose,
				EnvFile:    envVars,
			}

			// CALL OUR MAIN FUNCTIONALITY LOGIC FROM HERE AND SUPPLY BOTH PROJECT DIR AND SERVICE DIR
			clientCalls, serverCalls, err := discoverAllCalls(config)
			if err != nil {
				return err
			}

			fmt.Println("Successfully analysed, here is a list of dependencies:")

			// generate output
			graph := matching.CreateDependencyGraph(clientCalls, serverCalls)
			adjacencyList := output.ConstructAdjacencyList(graph)
			JSON, err := output.SerializeAdjacencyList(adjacencyList, true)
			if err != nil {
				return err
			}

			// print output
			// TODO: output to file
			fmt.Println(JSON)

			return nil
		},
	}
	cmd.Flags().StringVarP(&projectDir, "project-directory", "p", "./", "project directory")
	cmd.Flags().StringVarP(&serviceDir, "service-directory", "s", "./svc", "service directory")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "toggle logging trace of unknown variables")
	cmd.Flags().StringVarP(&envVars, "environment-variables", "e", "", "environment variable file")
	return cmd
}

// init initialises the depScan command and adds it as a subcommand of the root
func init() {
	depScanCmd := depScanCmd()
	rootCmd.AddCommand(depScanCmd)
}

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
	return stages.MapEnvVars(path)
}

// discoverAllCalls calls the correct stages for loading, building,
// filtering and discovering all client and server calls.
func discoverAllCalls(config RunConfig) ([]*callanalyzer.CallTarget, []*callanalyzer.CallTarget, error) {
	// Given a correct project directory en service directory,
	// apply our discovery algorithm to find all interesting calls
	if ex, err := pathExists(config.EnvFile); !ex && config.EnvFile != "" || err != nil {
		return nil, nil, fmt.Errorf("invalid environment variable file specified: %s", config.EnvFile)
	}

	// Filtering
	services, err := stages.FindServices(config.ServiceDir)
	fmt.Printf("Starting to analyse %d services.\n", len(services))

	if err != nil {
		return nil, nil, err
	}

	// resolve environment values
	// TODO: Integrate the envVariables into discovery
	envVariables, err := resolveEnvironmentValues(config.EnvFile)
	if err != nil {
		return nil, nil, err
	}

	analyseConfig := callanalyzer.DefaultConfigForFindingHTTPCalls(envVariables)
	analyseConfig.SetVerbose(config.Verbose)

	allClientTargets := make([]*callanalyzer.CallTarget, 0)
	allServerTargets := make([]*callanalyzer.CallTarget, 0)

	packageCount := 0

	for _, serviceDir := range services {
		// load packages
		packagesInService, err := stages.LoadAndBuildPackages(projectDir, serviceDir)
		packageCount += len(packagesInService)

		if err != nil {
			return nil, nil, err
		}

		// discover calls
		clientCalls, serverCalls, err := discovery.DiscoverAll(packagesInService, &analyseConfig)
		if err != nil {
			return nil, nil, err
		}

		// append
		allClientTargets = append(allClientTargets, clientCalls...)
		allServerTargets = append(allServerTargets, serverCalls...)
	}

	if packageCount == 0 {
		return nil, nil, fmt.Errorf("no service to analyse were found")
	}

	if err != nil {
		return nil, nil, err
	}

	return allClientTargets, allServerTargets, err
}
