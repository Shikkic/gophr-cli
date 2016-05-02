package main

import (
	"go/parser"
	"go/token"
	"strings"
)

// Define Dependency Struct
type Dependency struct {
	name, version string
	installed     bool
}

// Parse Dependencies from a .go file
func ParseDeps(fileName string) []string {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, fileName, nil, parser.ImportsOnly)
	Check(err)

	depsArray := make([]string, len(f.Imports))
	for index, s := range f.Imports {
		depName := strings.Replace(s.Path.Value, string('"'), " ", 2)
		depName = strings.Replace(depName, " ", "", 10)
		depsArray[index] = depName
	}

	return depsArray
}

func DepExistsInList(depName string, depArray []string) bool {
	for _, currDepName := range depArray {
		if currDepName == depName {
			return true
		}
	}

	return false
}

// Returns an array of built dependency structs from an array of dep names.
func buildDependencyStructs(depNames []string) {

}

// Return a map of dependencies that have the attributes installed or missing
func validateDepIsInstalled(depName string) {

}

func Check(e error) {
	if e != nil {
		panic(e)
	}
}
