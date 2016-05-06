package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

// TODO consider renaming to more specific
func ReadFiles(goFiles []string) {
	if len(goFiles) == 0 {
		path, err := os.Getwd()
		Check(err)
		fmt.Println(path)
		fmt.Println("└── (empty)")
		fmt.Println("")
		os.Exit(3)
	}

	for _, goFile := range goFiles {
		ReadFile(goFile)
	}
}

func ReadFile(goFilePath string) {
	depsArray := ParseDeps(goFilePath)
	// TODO Check to determine all github which packages are installed for
	// use map to distinguish
	printDeps(depsArray, goFilePath)
}

func printDeps(depsArray []string, goFileName string) {
	fmt.Print("Go Dependecies for ")
	color.Blue(goFileName)

	for index, depName := range depsArray {
		if index == (len(depsArray) - 1) {
			if strings.Contains(depName, "github") || strings.Contains(depName, "gophr.dev") {
				color.Green("└── " + depName + "\n")
			} else {
				fmt.Println("└── " + depName + "\n")
			}
		} else {
			if strings.Contains(depName, "github") || strings.Contains(depName, "gophr.dev") {
				color.Green("├─┬ " + depName)
			} else {
				fmt.Println("├─┬ " + depName)
			}
		}
	}
}

func appendDepsToBuffer(buffer []byte, depName []byte) []byte {
	for _, token := range depName {
		buffer = append(buffer, token)
	}

	return buffer
}
