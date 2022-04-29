/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	project_dir string
	service_dir string

	depScanCmd = &cobra.Command{
		Use:   "depScan",
		Short: "Scan and report dependencies between microservices",
		Long: `Outputs network-communication-based dependencies of services within a microservice architecture Golang project. 
Output is an adjacency list of service dependencies in a JSON format`,

		RunE: func(cmd *cobra.Command, args []string) error {
			// Path validation
			if ex, err := exists(project_dir); ex == false || err != nil {
				return fmt.Errorf("invalid project directory specified: %s", project_dir)
			}
			if ex, err := exists(service_dir); ex == false || err != nil {
				return fmt.Errorf("invalid service directory specified: %s", service_dir)
			}

			// CALL OUR MAIN FUNCTIONALITY LOGIC FROM HERE AND SUPPLY BOTH PROJECT DIR AND SERVICE DIR
			fmt.Println("depScan called")
			fmt.Println("Project directory: " + project_dir)
			fmt.Println("Service directory: " + service_dir)
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(depScanCmd)

	// Here you will define your flags and configuration settings.
	depScanCmd.PersistentFlags().StringVarP(&project_dir, "project-directory", "p", "./", "project directory")
	depScanCmd.PersistentFlags().StringVarP(&service_dir, "service-directory", "s", "./svc", "service directory")
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
