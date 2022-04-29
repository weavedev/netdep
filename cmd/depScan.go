/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

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

		Run: func(cmd *cobra.Command, args []string) {
			// CALL OUR MAIN FUNCTIONALITY LOGIC FROM HERE AND SUPPLY BOTH PROJECT DIR AND SERVICE DIR
			fmt.Println("depScan called")
			fmt.Println("Project directory: " + project_dir)
			fmt.Println("Service directory: " + service_dir)
		},
	}
)

func init() {
	rootCmd.AddCommand(depScanCmd)

	// Here you will define your flags and configuration settings.
	depScanCmd.PersistentFlags().StringVarP(&project_dir, "project-directory", "p", "./", "project directory")
	depScanCmd.PersistentFlags().StringVarP(&service_dir, "service-directory", "s", "./svc", "service directory")
}
