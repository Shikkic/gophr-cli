# gophr-cli

Gophr-cli is the cli tool for the [Gophr](https://github.com/skeswa/gophr) end-to-end package management solution for Go.

### Prerequisites
- Must have [go](https://golang.org/) installed

Run to find out if you have go installed:
```

```
If not installation instructions can be found [here](https://golang.org/dl/)

- Go **MUST** be installed correctly
  - `$GOBIN` and `$GOPATH` must exist in your `$PATH` env

Run to find out if your go has been properly setup:
```

```
If not setup instructions can be found [here]()

- [Gophr] dev environmnet MUST be running
Installation and setup instructions can be found [here](https://github.com/skeswa/gophr)


### Installation Gophr Development 

Clone the repo:
```
$ git clone git@github.com:Shikkic/gophr-cli.git
```

Navigate to the $GOPHR_REPO:
```
$ cd $GOPHR_REPO
```

Build and install the go files:
```
$ go build && go install
```

You should not be able to call `gophr` like so:
```
$ gophr
```
