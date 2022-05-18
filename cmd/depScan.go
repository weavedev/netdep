/*
Package cmd contains all the application command definitions
Copyright © 2022 TW Group 13C, Weave BV, TU Delft
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/discovery"
	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/discovery/callanalyzer"

	"github.com/spf13/cobra"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages"
)

var (
	projectDir string
	serviceDir string
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
			if ex, err := pathExists(path.Join(projectDir, serviceDir)); !ex || err != nil {
				return fmt.Errorf("invalid service directory specified: %s", serviceDir)
			}

			// CALL OUR MAIN FUNCTIONALITY LOGIC FROM HERE AND SUPPLY BOTH PROJECT DIR AND SERVICE DIR

			fmt.Println("Starting call scanning...")
			fmt.Println("Project directory: " + projectDir)
			fmt.Println("Service directory: " + serviceDir)
			dependencies, err := buildDependencies(serviceDir, projectDir)
			if err != nil {
				return err
			}

			fmt.Println("Successfully analysed, here is a list of dependencies:")
			for _, dependency := range dependencies {
				fmt.Println(json.Marshal(dependency))
			}

			return nil
		},
	}
	cmd.Flags().StringVarP(&projectDir, "project-directory", "p", "./", "project directory")
	cmd.Flags().StringVarP(&serviceDir, "service-directory", "s", "./svc", "service directory")
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

// buildDependencies is responsible for integrating different stages
// of the program.
// TODO: the output should be changed to a list of string once the integration is done
func buildDependencies(svcDir string, projectDir string) ([]*callanalyzer.CallTarget, error) {
	// Filtering
	initial, err := stages.LoadServices(projectDir, svcDir)
	fmt.Printf("Starting to analyse %s\n", initial)
	if err != nil {
		return nil, err
	}

	// TODO: Endpoint discovery
	// Client Call Discovery
	clientCalls, err := discovery.Discover(initial)
	if err != nil {
		return nil, err
	}

	// TODO: Matching

	// For now this returns client calls,
	// as we don't have any other functionality in place.
	return clientCalls, nil
}
