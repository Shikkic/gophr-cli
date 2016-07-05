# gophr-cli

:package: Gophr-cli is the cli tool for the [Gophr](https://github.com/skeswa/gophr) end-to-end package management solution for Go.

```
NAME:
   gophr - An end-to-end package management solution for Go

USAGE:
   gophr [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
    search	Search go packages on gophr.pm
    deps	List go packages of a specified go file or folder
    install	Install dependency
    uninstall	Uninstall dependency
    init	initialize new project
    lock	Lock a file(s) github go packages to use gophr.pm/<REPO_NAME>

GLOBAL OPTIONS:
   --dev		enable developer mode on commands
   --help, -h		show help
   --version, -v	print the version
```

### Prerequisites
- Must have [go](https://golang.org/) installed

Run to find out if you have go installed:
```sh
$ go && echo "Go is installed"
```
If not installation instructions can be found [here](https://golang.org/dl/)

- Go **MUST** be installed correctly
  - `$GOBIN` and `$GOPATH` must exist in your `$PATH` env

Run to find out if your go has been properly setup:
```sh
$ echo $GOBIN && echo $GOPATH && echo "Go is properly setup"
```
If not setup instructions can be found [here]()

- [Gophr](https://github.com/skeswa/gophr) dev environmnet MUST be running
Installation and setup instructions can be found [here](https://github.com/skeswa/gophr)


### Compiling Gophr-cli from source

Clone the repo:
```
$ git clone git@github.com:Shikkic/gophr-cli.git
```

Navigate to the $GOPHR_REPO:
```sh
$ cd $GOPHR_REPO/gophr
```

Build and install the go files:
```sh
$ go build && go install
```

You should not be able to call `gophr` like so:
```
$ gophr --help
```

### Developer Mode

Gophr-cli's default mode is production, but as of (05/12/2016) there is no PROD server that can handle requests. You **MUST** be running a `gophr` dev server, and you must use the `--dev` flag on all commands or they will not work
