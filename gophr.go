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
	"unicode"
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

	for {
		// create 7 byte read buffer
		readBuffer := make([]byte, bufferSize)
		n1, err := fileRef.Read(readBuffer)
		var _ = n1

		// TODO move this into function
		if err != nil {
			if err == io.EOF {
				fmt.Println("EOF")
				fileRef.Close()
				break
			}
			fmt.Println(err)
		}

		// For each read buffer, parse and put it into
		for _, token := range readBuffer {
			token := rune(token)
			if token != '\n' && !unicode.IsSpace(token) {
				// check if importStringBuffer == 'import'
				if strings.Join(importStringbuffer[:], "") == "import" {
					fmt.Println("Found keyword 'import'")
				}

				importStringbuffer = append(importStringbuffer[:1], importStringbuffer[1+1:]...)
				importStringbuffer = append(importStringbuffer, string(token))

				// Once found parse for deps
				// TODO create an array of string dep arrays
			}
		}
	}
}

func printDeps(goFile string) {
	importPackages := strings.Split(string(goFile), "\n")
	fmt.Println("Go Dependecies for this package")
	for i := 0; i < len(importPackages); i++ {
		depName := importPackages[i]
		depName = strings.Replace(depName, string('"'), " ", 2)
		depName = strings.Replace(depName, " ", "", 10)

		if strings.Contains(depName, "github.com") {
			color.Green("├─┬" + depName)
		}
	}
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
