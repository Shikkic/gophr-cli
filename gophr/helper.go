package main

import (
	"go/parser"
	"go/token"
	"strings"

	"github.com/fatih/color"
)

// Define Dependency Struct
type Dependency struct {
	name, version string
	installed     bool
}

// Parse Dependencies from a .go file
// TODO consider returning if the file name does not exist
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

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

/*
/////////////////////////////
//	 Text Color Functions  //
/////////////////////////////
# Use these color functions with printf() to make stdout colored
*/

func Magenta(text string) string {
	magenta := color.New(color.FgMagenta).SprintFunc()
	return magenta(text)
}

func Red(text string) string {
	red := color.New(color.FgRed).SprintFunc()
	return red(text)
}

func Green(text string) string {
	green := color.New(color.FgGreen).SprintFunc()
	return green(text)
}

func Blue(text string) string {
	blue := color.New(color.FgBlue).SprintFunc()
	return blue(text)
}
