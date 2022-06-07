
# Microservice Architecture Analysis Tool

A command line utility for service network dependencies using static code analysis built in Go.


## Features

- Easy to use command line interface




## Installation

Clone the project and install using

```bash
  git clone ...
  go get -u ./...
```

**Important: For the tool to work, the project directory it is run against _has_ to compile and must have all its dependencies installed!**
## Usage

To initiate dependency scanning depScan command has to be used in a command-line interface.
ou can run the tool using the default settings using:

```sh
go run main.go depScan
```
or if you want more control you can use the options as defined below, for example:
```sh
go run main.go depScan -p "./some/project/dir" -s "./some/service/dir"
```

### Annotations

The tool supports code annotations. This is necessary, because it might fail to resolve some of the variables due to complexity of the code or lack of support.

#### Annotation format

User can add annotations as comments in their project before running the tool on it. Currently, the tool supports 2 types of annotations:
1) Annotations for client calls. Example:
```go
//netdep:client url=http://example.com/ targetSvc=service-2
c.Do(req)
```
Both `url` and `targetSvc` can be specified, but any one is also enough.
2) Annotations for endpoint definitions. Example:
```go
//netdep:endpoint url=http://example.com/ping
r.GET("/ping", func(c *gin.Context) {
    c.JSON(200, gin.H{
    "message": "pong",
    })
})
```

#### Annotation suggestions

An annotation suggestion will be printed for all unresolved targets. 

Suggestion output example:
```sh
service-1\main.go:24 couldn't be resolved. Add an annotation above it in the format "//netdep:client ..." or "//netdep:endpoint ..."
```


### Options

| Argument                  | Description                                                        | Default  |
|:--------------------------|:-------------------------------------------------------------------|:---------|
| `-h, --help`              | Print help                                                         |          |
| `-p, --project-directory` | The path to the project directory. Must be a valid path.           | `./`     |
| `-s --service-directory`  | The path to the services inside the project. Must be a valid path. | `./svc/` |
| `-v --verbose`            | Toggle printing stack traces of unknown variables.                 | `false`  |


## Documentation

Documentation can be found in the [project wiki](https://gitlab.ewi.tudelft.nl/cse2000-software-project/2021-2022-q4/cluster-13/microservice-architecture-analysis-tool/code/-/wikis/home).


## Roadmap

- CI/CD integration
- High(er) degree of accuracy
- Extensive test suite



## Contributing

See `CONTRIBUTING.md` for ways to get started.

Please adhere to this project's `code of conduct`.

## Authors

- Gints Kuļikovskis
- Martynas Krupskis
- Bianca Şerbănescu
- Lukas Ostrovskis
- Delano Flipse


## License


