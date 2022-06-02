/*
Package cmd contains all the application command definitions
Copyright © 2022 TW Group 13C, Weave BV, TU Delft
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/preprocessing"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/matching"
	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/output"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/discovery"
	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/discovery/callanalyzer"

	"github.com/spf13/cobra"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages"
)

var (
	projectDir   string
	serviceDir   string
	envVars      string
	jsonFilename string
)

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

			// Given a correct project directory en service directory,
			// apply our discovery algorithm to find all interesting calls
			if ex, err := pathExists(envVars); !ex && envVars != "" || err != nil {
				return fmt.Errorf("invalid environment variable file specified: %s", envVars)
			}

			// CALL OUR MAIN FUNCTIONALITY LOGIC FROM HERE AND SUPPLY BOTH PROJECT DIR AND SERVICE DIR
			clientCalls, serverCalls, err := discoverAllCalls(serviceDir, projectDir, envVars)
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
	cmd.Flags().StringVarP(&envVars, "environment-variables", "e", "", "environment variable file")
	cmd.Flags().StringVarP(&envVars, "json-filename", "j", "", "json output filename")
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
func discoverAllCalls(svcDir string, projectDir string, envVars string) ([]*callanalyzer.CallTarget, []*callanalyzer.CallTarget, error) {
	// Filtering
	services, err := preprocessing.FindServices(svcDir)
	fmt.Printf("Starting to analyse %d services.\n", len(services))

	if err != nil {
		return nil, nil, err
	}

	// resolve environment values
	// TODO: Integrate the envVariables into discovery
	envVariables, err := resolveEnvironmentValues(envVars)
	if err != nil {
		return nil, nil, err
	}

	allClientTargets := make([]*callanalyzer.CallTarget, 0)
	allServerTargets := make([]*callanalyzer.CallTarget, 0)
	annotations := make(map[string]map[preprocessing.Position]string)

	packageCount := 0

	config := callanalyzer.DefaultConfigForFindingHTTPCalls(envVariables)

	for _, serviceDir := range services {
		// load packages
		packagesInService, err := preprocessing.LoadAndBuildPackages(projectDir, serviceDir)
		if err != nil {
			return nil, nil, err
		}
		packageCount += len(packagesInService)

		serviceName := strings.Split(serviceDir, "\\")[len(strings.Split(serviceDir, "\\"))-1]
		err = preprocessing.LoadAnnotations(serviceDir, serviceName, annotations)
		if err != nil {
			return nil, nil, err
		}

		// discover calls
		clientCalls, serverCalls, err := discovery.DiscoverAll(packagesInService, &config)
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

	// TODO: make use of annotations in the matching stage
	fmt.Println("Discovered annotations:")
	for k1, serMap := range annotations {
		for k2, val := range serMap {
			fmt.Println("Service name: " + k1)
			fmt.Print("Position: " + k2.Filename + ":")
			fmt.Println(k2.Line)
			fmt.Println("Value: " + val)
		}
	}

	return allClientTargets, allServerTargets, err
}
