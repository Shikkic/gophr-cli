package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/codegangsta/cli"
)

func RunDepsCommand(c *cli.Context) error {
	fileNameArg := c.Args().First()

	if FileNameArgIsEmpty(fileNameArg) {
		err := PrintDepsFromCurrentDirectory()
		return err
	}

	err := PrintDepsFromFileName(fileNameArg)
	Check(err)
	return nil
}

// DepsCommand Helpers
func FileNameArgIsEmpty(fileNameArg string) bool {
	if len(fileNameArg) == 0 {
		return true
	}

	return false
}

// TODO finishing building this
func validateFileNameIsValid(fileNameArg string) error {
	// need to return special error here
	return nil
}

func PrintDepsFromCurrentDirectory() error {
	goFilesInCurrentDir, err := filepath.Glob("./*.go")
	if err != nil {
		return err
	}

	for _, goFile := range goFilesInCurrentDir {
		err := PrintDepsFromFileName(goFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func PrintDepsFromFileName(fileName string) error {
	err := validateFileNameIsValid(fileName)
	Check(err)

	file, err := OpenASTFilePointerFromFileName(fileName)
	if err != nil {
		return err
	}

	fileDepURLs := ParseDepURLsFromFile(file)
	PrintFileDepURLsAndFileName(fileDepURLs, fileName)

	return nil
}

func PrintFileDepURLsAndFileName(depsArray []string, goFileName string) {
	fmt.Printf("\n%s\n", Blue(goFileName))

	for index, depName := range depsArray {
		if index == (len(depsArray) - 1) {
			// TODO create get function for gophr domain
			if strings.Contains(depName, "github") || strings.Contains(depName, "gophr.dev") {
				fmt.Printf("└── ⚠ %s\n\n", Yellow(depName))
			} else {
				fmt.Printf("└── %s\n", depName)
			}
		} else {
			// TODO create get function for gophr domain
			if strings.Contains(depName, "github") || strings.Contains(depName, "gophr.dev") {
				fmt.Printf("├─┬ ⚠ %s\n", Yellow(depName))
			} else if strings.Contains(depName, "gophr.dev") {
				fmt.Printf("├─┬ ✓%s\n", Green(depName))
			} else {
				fmt.Println("├─┬ " + depName)
			}
		}
	}
}

func OpenASTFilePointerFromFileName(fileName string) (*ast.File, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fileName, nil, parser.ImportsOnly)

	return f, err
}
