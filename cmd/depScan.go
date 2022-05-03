/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var projectDir string
var serviceDir string

func depScanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "depScan",
		Short: "Scan and report dependencies between microservices",
		Long: `Outputs network-communication-based dependencies of services within a microservice architecture Golang project.
Output is an adjacency list of service dependencies in a JSON format`,

		RunE: func(cmd *cobra.Command, args []string) error {
			// Path validation
			if ex, err := exists(projectDir); !ex || err != nil {
				return fmt.Errorf("invalid project directory specified: %s", projectDir)
			}
			if ex, err := exists(serviceDir); !ex || err != nil {
				return fmt.Errorf("invalid service directory specified: %s", serviceDir)
			}

			// CALL OUR MAIN FUNCTIONALITY LOGIC FROM HERE AND SUPPLY BOTH PROJECT DIR AND SERVICE DIR
			fmt.Println("depScan called")
			fmt.Println("Project directory: " + projectDir)
			fmt.Println("Service directory: " + serviceDir)

			return nil
		},
	}
	cmd.Flags().StringVarP(&projectDir, "project-directory", "p", "./", "project directory")
	cmd.Flags().StringVarP(&serviceDir, "service-directory", "s", "./svc", "service directory")
	return cmd
}

func init() {
	depScanCmd := depScanCmd()
	rootCmd.AddCommand(depScanCmd)
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
