package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/fatih/color"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Define Constants
const bufferSize = 7

func main() {
	app := cli.NewApp()
	app.Name = "gophr"
	app.Usage = "A good go package manager"
	// TODO create dep sub action
	app.Action = func(c *cli.Context) {
		// Retrieve the file names of all the go files in the current dir
		fls, err := filepath.Glob("*.go")
		check(err)
		// Readfile
		readFiles(fls)
	}
	app.Run(os.Args)
}

/*
Print Deps Functions
*/

func readFiles(goFiles []string) {
	for _, goFile := range goFiles {
		readFile(goFile)
	}
}

func readFile(goFilePath string) {
	// create a file reference and open it
	fileRef, err := os.Open(goFilePath)
	check(err)

	importStringbuffer := make([]string, 7)
	depsBuffer := make([]rune, 1)
	var foundImportStatement bool = false
	var isInImport bool = false
	var importCheckCount int = 0

	for {
		// create 7 byte read buffer
		readBuffer := make([]byte, bufferSize)
		_, err := fileRef.Read(readBuffer)

		// TODO move this into function
		// TODO check the end :
		if err != nil {
			if err == io.EOF {
				fileRef.Close()
				break
			}
			fmt.Println(err)
		}

		// For each read buffer, parse and put it into
		for _, token := range readBuffer {
			token := rune(token)

			if foundImportStatement {
				if isInImport {
					if token != ')' {
						depsBuffer = append(depsBuffer, token)
					} else {
						isInImport = false
						foundImportStatement = false
						printDeps(string(depsBuffer[:len(depsBuffer)]), goFilePath)
					}
				} else {
					if importCheckCount < 2 {
						if token == '"' || token == '(' {
							isInImport = true
						}
						importCheckCount++
					} else {
						foundImportStatement = false
					}
				}
			} else {
				// check if importStringBuffer == 'import'
				if strings.Join(importStringbuffer[:], "") == "import" {
					foundImportStatement = true
				}
			}
			importStringbuffer = append(importStringbuffer[:1], importStringbuffer[1+1:]...)
			importStringbuffer = append(importStringbuffer, string(token))
		}
	}

	// Once found parse for deps
	// printDeps
}

func printDeps(depsArray string, goFilePath string) {
	depsArray = strings.Trim(depsArray, "\n\t\x00 ")
	importPackages := strings.Split(depsArray, "\n")
	// include file name when listing dependencies
	fmt.Print("Go Dependecies for ")
	color.Blue(goFilePath)
	for i := 0; i < len(importPackages); i++ {
		depName := importPackages[i]
		depName = strings.Replace(depName, string('"'), " ", 2)
		depName = strings.Replace(depName, "\t", "", 10)

		if i == (len(importPackages) - 1) {
			// This is the last dependency
			color.Green("└──" + depName)

		} else {
			color.Green("├─┬" + depName)
		}
	}
	fmt.Println("")
}

/*
Helper Functions
*/

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// NOTE might need this for later functionality
func getRequest(string) string {
	resp, err := http.Get("http://github.com/shikkic")
	if err != nil {
		return "nil"
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body)
}
