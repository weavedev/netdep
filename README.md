## Usage
To initiate dependency scanning `depScan` command has to be used in a command-line interface.
There's an option to include directories for the target project and its services. Flags have both full and shorthand
versions and usage looks like this:
- `go run main.go depScan` without flags assumes a default project directory path `./` and default service directory path `./svc`
- `go run main.go depScan --project-directory "./some/project/dir" --service-directory "./some/service/directory"` sets the project directory and service directory variables and **validates** the paths
  - Invalid or non-existing directories will result in an error and terminate the tool
- `go run main.go depScan -p "./some/project/dir" -s "./some/service/dir"` shorthand versions


