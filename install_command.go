package main

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func RunInstallCommand(depName string, fileName string) {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()

	// Step 1 Determine if we need to download and install dependencies of the folder or for specified dependency
	// TODO instead use a flag to determine if it should install to all packages
	if len(depName) == 0 {
		fls, err := filepath.Glob("*.go")
		Check(err)

		if len(fls) == 0 {
			red := color.New(color.FgRed).SprintFunc()
			magenta := color.New(color.FgMagenta).SprintFunc()
			s.Stop()
			fmt.Printf("gophr %s %s not run in go a package\n", red("ERROR"), magenta("install"))
			os.Exit(3)
		}

		depName = "./"
	}

	// Step 2 run get command if
	if strings.Contains(depName, "github.com") {
		cmd := exec.Command("go", "get", depName)
		//var out bytes.Buffer
		//cmd.Stdout = &out
		err := cmd.Run()
		Check(err)
	}

	// Step 3 if command was successful, append to file
	if len(fileName) > 0 {
		// add to file
		file, err := ioutil.ReadFile(fileName)
		Check(err)
		augmentImportStatement(file, fileName, depName)
	}

	// Step 4 after adding it to import statement run go fmt on file
	cmd := exec.Command("go", "fmt", fileName)
	err := cmd.Run()
	Check(err)

	// Check if exits in file
	depsArray := ParseDeps(fileName)
	if DepExistsInList(depName, depsArray) == true {
		magenta := color.New(color.FgMagenta).SprintFunc()
		s.Stop()
		fmt.Printf("✓ %s was successfully installed into %s\n", magenta("'"+depName+"'"), magenta(fileName))
		os.Exit(3)
	} else {
		//PANIC ITS NOT THERE
		// TODO PANIC
	}

	s.Stop()
}

func augmentImportStatement(file []byte, fileName string, depName string) {
	formatedDepName := []byte("\n\t" + string('"') + depName + string('"'))
	importStringbuffer := make([]string, 7)
	newFileBuffer := make([]byte, 0)
	depsBuffer := make([]rune, 1)
	var foundImportStatement bool = false
	var isInImport bool = false
	var importCheckCount int = 0
	var addedImport bool = false

	for _, token := range file {
		newFileBuffer = append(newFileBuffer, token)
		token := rune(token)
		if addedImport == false {
			if foundImportStatement {
				if isInImport {
					if token != ')' {
						depsBuffer = append(depsBuffer, token)
					} else {
						isInImport = false
						foundImportStatement = false
						addedImport = true
					}
				} else {
					if importCheckCount < 2 {
						if token == '(' {
							isInImport = true
							newFileBuffer = appendDepsToBuffer(newFileBuffer, formatedDepName)
						}
						importCheckCount++
					} else {
						foundImportStatement = false
					}
				}
			} else {
				if strings.Join(importStringbuffer[:], "") == "import" {
					foundImportStatement = true
				}
			}

			importStringbuffer = append(importStringbuffer[:1], importStringbuffer[1+1:]...)
			importStringbuffer = append(importStringbuffer, string(token))
		}
	}

	err := ioutil.WriteFile(fileName, newFileBuffer, 0644)
	Check(err)
}
