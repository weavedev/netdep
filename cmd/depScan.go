/*
Package cmd contains all the application command definitions
Copyright © 2022 TW Group 13C, Weave BV, TU Delft
*/
package cmd

import (
	"fmt"
	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/discovery"
	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/discovery/callanalyzer"
	"os"
	"path/filepath"
	"strings"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/preprocessing"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/matching"
	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages/output"

	"github.com/spf13/cobra"

	"lab.weave.nl/internships/tud-2022/static-analysis-project/stages"
)

// RunConfig defines the parameters for a depScan command run
type RunConfig struct {
	ProjectDir       string
	ServiceDir       string
	EnvFile          string
	Verbose          bool
	ServiceCallsDir  string
	OnlyServiceCalls bool
}

// depScanCmd creates and returns a depScan command object
func depScanCmd() *cobra.Command {
	// Variables that are supplied as command-line args
	var (
		projectDir       string
		serviceDir       string
		envVars          string
		outputFilename   string
		verbose          bool
		serviceCallsDir  string
		onlyServiceCalls bool
	)

	cmd := &cobra.Command{
		Use:   "depScan",
		Short: "Scan and report dependencies between microservices",
		Long: `Outputs network-communication-based dependencies of services within a microservice architecture Golang project.
Output is an adjacency list of service dependencies in a JSON format`,

		RunE: func(cmd *cobra.Command, args []string) error {
			ok, err := areInputPathsValid(projectDir, serviceDir, serviceCallsDir, envVars, outputFilename)
			if !ok {
				return err
			}

			config := RunConfig{
				ProjectDir:       projectDir,
				ServiceDir:       serviceDir,
				Verbose:          verbose,
				EnvFile:          envVars,
				ServiceCallsDir:  serviceCallsDir,
				OnlyServiceCalls: onlyServiceCalls,
			}

			// CALL OUR MAIN FUNCTIONALITY LOGIC FROM HERE AND SUPPLY BOTH PROJECT DIR AND SERVICE DIR
			clientCalls, serverCalls, err := discoverAllCalls(config)
			if err != nil {
				return err
			}

			// generate output
			graph := matching.CreateDependencyGraph(clientCalls, serverCalls)
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
	cmd.Flags().StringVarP(&serviceCallsDir, "servicecalls-directory", "c", "", "servicecalls package directory")
	cmd.Flags().BoolVar(&onlyServiceCalls, "only-servicecalls", false, "only scan the servicecalls")

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
func areInputPathsValid(projectDir, serviceDir, serviceCallsDir, envVars, outputFilename string) (bool, error) {
	if !pathOk(projectDir) {
		return false, fmt.Errorf("invalid project directory specified: %s", projectDir)
	}

	if !pathOk(serviceDir) {
		return false, fmt.Errorf("invalid service directory specified: %s", serviceDir)
	}

	if !pathOk(serviceCallsDir) && serviceCallsDir != "" {
		return false, fmt.Errorf("invalid servicecalls directory specified: %s", serviceCallsDir)
	}

	if !pathOk(envVars) && envVars != "" {
		return false, fmt.Errorf("invalid environment variable file specified: %s", envVars)
	}
	jsonParentDir := filepath.Dir(outputFilename)

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
	// Filtering
	services, err := preprocessing.FindServices(config.ServiceDir)
	fmt.Printf("Starting to analyse %d services.\n", len(services))

	if err != nil {
		return nil, nil, err
	}

	internalCalls, serverTargets, err := preprocessing.FindServiceCalls(config.ServiceCallsDir)
	if err != nil {
		return nil, nil, err
	}

	// resolve environment values
	// TODO: Integrate the envVariables into discovery
	envVariables, err := resolveEnvironmentValues(config.EnvFile)
	if err != nil {
		return nil, nil, err
	}

	analyserConfig := callanalyzer.DefaultConfigForFindingHTTPCalls()
	analyserConfig.SetVerbose(config.Verbose)
	analyserConfig.SetEnv(envVariables)

	allClientTargets, allServerTargets, annotations, err := processEachService(&services, &config, &analyserConfig, internalCalls, serverTargets)
	if err != nil {
		return nil, nil, err
	}

	if config.Verbose {
		fmt.Println("Discovered annotations:")
		anyHits := false
		for k1, serMap := range annotations {
			for k2, val := range serMap {
				anyHits = true
				fmt.Println("Service name: " + k1)
				fmt.Print("Position: " + k2.Filename + ":")
				fmt.Println(k2.Line)
				fmt.Println("Value: " + val)
			}
		}
		if !anyHits {
			fmt.Println("[Discovered none]")
		}
	}

	return allClientTargets, allServerTargets, err
}

// processEachService preprocesses and analyses each of the services using RunConfig and callanalyzer.AnalyserConfig
func processEachService(services *[]string, config *RunConfig, analyserConfig *callanalyzer.AnalyserConfig, internalCalls map[preprocessing.IntCall]string, serverTargets *[]*callanalyzer.CallTarget) ([]*callanalyzer.CallTarget, []*callanalyzer.CallTarget, map[string]map[callanalyzer.Position]string, error) {
	allClientTargets := make([]*callanalyzer.CallTarget, 0)
	allServerTargets := make([]*callanalyzer.CallTarget, 0)
	annotations := make(map[string]map[callanalyzer.Position]string)

	analyserConfig.SetAnnotations(annotations)

	packageCount := 0

	allServerTargets = append(allServerTargets, *serverTargets...)
	internalClientTargets := make([]*callanalyzer.CallTarget, 0)

	for _, serviceDir := range *services {
		serviceName := strings.Split(serviceDir, string(os.PathSeparator))[len(strings.Split(serviceDir, string(os.PathSeparator)))-1]

		if config.Verbose {
			fmt.Println("Analysing service " + serviceDir)
		}

		err := preprocessing.LoadAnnotations(serviceDir, serviceName, annotations)
		if err != nil {
			return nil, nil, nil, err
		}

		// There are some interesting internal calls so the tool should parse all methods
		if len(internalCalls) != 0 {
			err = preprocessing.LoadServiceCalls(serviceDir, serviceName, internalCalls, &internalClientTargets)
			if err != nil {
				return nil, nil, nil, err
			}
		}

		// Load and build packages and proceed with discovery if the user
		// Didn't ask for scanning of only servicecalls package
		if !config.OnlyServiceCalls {
			// load packages
			packagesInService, err := preprocessing.LoadAndBuildPackages(config.ProjectDir, serviceDir)
			if err != nil {
				return nil, nil, nil, err
			}
			packageCount += len(packagesInService)

			// discover calls
			clientCalls, serverCalls, err := discovery.DiscoverAll(packagesInService, analyserConfig)
			if err != nil {
				return nil, nil, nil, err
			}

			// append
			allClientTargets = append(allClientTargets, clientCalls...)
			allServerTargets = append(allServerTargets, serverCalls...)
		}
	}
	allClientTargets = append(allClientTargets, internalClientTargets...)

	if !config.OnlyServiceCalls && packageCount == 0 {
		return nil, nil, nil, fmt.Errorf("no service to analyse were found")
	}
	return allClientTargets, allServerTargets, annotations, nil
}
