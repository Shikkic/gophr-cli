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
const readBufferSize = 7

func main() {
	app := cli.NewApp()
	app.Name = "gophr"
	app.Usage = "A good go package manager"
	// TODO Will need flags later
	/*
		app.Flags = []cli.Flag{
			cli.StringFlag{
				Name:  "deps",
				Value: "list dependencies",
				Usage: "list go dependencies in file(s)",
			},
		}
	*/
	app.Commands = []cli.Command{
		{
			Name:    "deps",
			Aliases: []string{"dependencies"},
			Usage:   "List dependencies of a go file or folder",
			Action: func(c *cli.Context) {
				fileName := c.Args().First()
				fmt.Println("")

				switch {
				case len(fileName) != 0:
					readFile(fileName)
				default:
					fls, err := filepath.Glob("*.go")
					check(err)
					readFiles(fls)
				}
			},
		},
		{
			Name:    "install",
			Aliases: []string{"install deps"},
			Usage:   "Install dependency",
			Action: func(c *cli.Context) {
				// Determine if argument is passed in
				// ...
				//fileName := c.Args().First()
			},
		},
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
		readBuffer := make([]byte, readBufferSize)
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
				if strings.Join(importStringbuffer[:], "") == "import" {
					foundImportStatement = true
				}
			}
			importStringbuffer = append(importStringbuffer[:1], importStringbuffer[1+1:]...)
			importStringbuffer = append(importStringbuffer, string(token))
		}
	}

	// printDeps
	// TODO Proably include some functions to clean up the depsBuffer for printing, to seperate concern
	printDeps(string(depsBuffer[:len(depsBuffer)]), goFilePath)
}

func printDeps(depsArray string, goFileName string) {
	depsArray = strings.Trim(depsArray, "\n\t\x00 ")
	importPackages := strings.Split(depsArray, "\n")
	// Clean up strings and remove non-github
	fmt.Print("Go Dependecies for ")
	color.Blue(goFileName)
	
	for i := 0; i < len(importPackages); i++ {
		depName := importPackages[i]
		depName = strings.Replace(depName, string('"'), " ", 2)
		depName = strings.Replace(depName, "\t", "", 10)

		if i == (len(importPackages) - 1) {
			// This is the last dependency
			color.Green("└──" + depName)

		} else {
			if strings.Contains(depName, "github") {
				color.Green("├─┬" + depName)
			} else {
				fmt.Println("├─┬" + depName)
			}
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
