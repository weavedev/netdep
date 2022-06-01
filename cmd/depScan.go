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
	envVars    string
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

			if ex, err := pathExists(envVars); !ex && envVars != "" || err != nil {
				return fmt.Errorf("invalid environment variable file specified: %s", envVars)
			}

			// CALL OUR MAIN FUNCTIONALITY LOGIC FROM HERE AND SUPPLY BOTH PROJECT DIR AND SERVICE DIR
			clientCalls, serverCalls, err := buildDependencies(serviceDir, projectDir, envVars)
			if err != nil {
				return err
			}

			fmt.Println("Successfully analysed, here is a list of dependencies:")

			graph := matching.CreateDependencyGraph(clientCalls, serverCalls)
			adjacencyList := output.ConstructAdjacencyList(graph)
			JSON, err := output.SerializeAdjacencyList(adjacencyList, true)
			if err != nil {
				return err
			}

			fmt.Println(JSON)

			return nil
		},
	}
	cmd.Flags().StringVarP(&projectDir, "project-directory", "p", "./", "project directory")
	cmd.Flags().StringVarP(&serviceDir, "service-directory", "s", "./svc", "service directory")
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

func envMap(path string) (map[string]map[string]string, error) {
	if path == "" {
		return nil, nil
	}
	return stages.MapEnvVars(path)
}

// buildDependencies is responsible for integrating different stages
// of the program.
// TODO: the output should be changed to a list of string once the integration is done
func buildDependencies(svcDir string, projectDir string, envVars string) ([]*callanalyzer.CallTarget, []*callanalyzer.CallTarget, error) {
	// Filtering
	initial, err := stages.LoadServices(projectDir, svcDir)
	fmt.Printf("Starting to analyse %s\n", initial)
	if err != nil {
		return nil, nil, err
	}

	// var envVariables map[string]map[string]string = nil

	/*if envVars != "" {
		envVariables, err := stages.MapEnvVars(envVars)
		fmt.Println("env: ")
		fmt.Println(envVariables)
		if err != nil {
			return nil, nil, err
		}
	}*/

	envVariables, err := envMap(envVars)
	if err != nil {
		return nil, nil, err
	}

	fmt.Println("env: ")
	fmt.Println(envVariables)
	// TODO: Integrate the envVariables into discovery

	// TODO: Endpoint discovery
	// Client Call Discovery
	clientCalls, serverCalls, err := discovery.Discover(initial, nil, envVariables)
	if err != nil {
		return nil, nil, err
	}

	// For now this returns client calls,
	// as we don't have any other functionality in place.
	return clientCalls, serverCalls, err
}
