package main

import (
	"fmt"
	"os"
)

func validateSearchQueryArg(searchQuery string) {
	if len(searchQuery) == 0 {
		newError := NewInvalidArgumentError("Search Query", searchQuery, 1)
		newError.PrintErrorAndExit()
		os.Exit(1)
	}
}

// TODO need to look up error codes and print message
func Check(e error) {
	if e != nil {
		panic(e)
	}
}

/***************************** INVALID ARGUMENT ******************************/

// InvalidArgumentError is an error that occurs when a particular parameter
// defies expectations.
type InvalidArgumentError struct {
	ArgumentName  string
	ArgumentValue interface{}
	ErrorCode     int
}

// NewInvalidParameterError creates a new InvalidParameterError.
func NewInvalidArgumentError(
	argumentName string,
	argumentValue interface{},
	errorCode int,
) InvalidArgumentError {
	return InvalidArgumentError{
		ArgumentName:  argumentName,
		ArgumentValue: argumentValue,
		ErrorCode:     errorCode,
	}
}

func (err InvalidArgumentError) Error() string {
	return fmt.Sprintf(
		`Invalid value "%v" specified for argument "%s".`,
		err.ArgumentValue,
		err.ArgumentName,
	)
}

func (err InvalidArgumentError) String() string {
	return err.Error()
}

func (err InvalidArgumentError) PrintErrorAndExit() {
	fmt.Println(err.Error())
	os.Exit(1)
}
