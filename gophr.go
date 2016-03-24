package main

import (
	"fmt"
	"github.com/briandowns/spinner"
	//"bytes"
	"github.com/codegangsta/cli"
	"github.com/fatih/color"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	//"reflect"
	"go/parser"
	"go/token"
	"strings"
	"time"
)

// Define Constants
const readBufferSize = 7

func main() {
	app := cli.NewApp()
	app.Name = "gophr"
	app.Usage = "A good go package manager"
	// TODO Will need flags later
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "deps",
			Value: "list dependencies",
			Usage: "list go dependencies in file(s)",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "deps",
			Aliases: []string{"dependencies"},
			Usage:   "List dependencies of a go file or folder",
			Action: func(c *cli.Context) {
				fileName := c.Args().First()
				parseDeps(fileName)
				// TODO FUNCTIONALITY check if deps are present in the go files AND if they're installed or not
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
				var depName string
				var fileName string

				// TODO Consider using -a or --all flag to re-install all dependencies
				if c.NArg() == 0 {
					red := color.New(color.FgRed).SprintFunc()
					magenta := color.New(color.FgMagenta).SprintFunc()
					fmt.Printf("gophr %s %s not run with a package name\n", red("ERROR"), magenta("install"))
					fmt.Printf("run %s for more help\n", magenta("gophr install -h"))
					os.Exit(3)
				}

				// TODO check if type string with reflect
				if c.NArg() < 2 {
					red := color.New(color.FgRed).SprintFunc()
					magenta := color.New(color.FgMagenta).SprintFunc()
					fmt.Printf("gophr %s %s not run with a file name\n", red("ERROR"), magenta("install"))
					fmt.Printf("run %s for more help\n", magenta("gophr install -h"))
					os.Exit(3)
				}

				if c.NArg() > 0 {
					depName = c.Args()[0]
				}

				// TODO consider tabbing for arg if not present
				if c.NArg() > 1 {
					fileName = c.Args()[1]
				}

				runGoGetCommand(depName, fileName)
			},
		},
	}
	app.Run(os.Args)
}

func parseDeps(fileName string) {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, fileName, nil, parser.ImportsOnly)
	check(err)

	fmt.Println(len(f.Imports))
	depsArray := make([]string, len(f.Imports))

	for index, s := range f.Imports {
		depsArray[index] = s.Path.Value
	}

	return depsArray
}

func runGoGetCommand(depName string, fileName string) {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()

	// Step 1 Determine if we need to download and install dependencies of the folder or for specified dependency
	// TODO instead use a flag to determine if it should install to all packages
	if len(depName) == 0 {
		fls, err := filepath.Glob("*.go")
		check(err)

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
		check(err)
	}

	// Step 3 if command was successful, append to file
	if len(fileName) > 0 {
		// add to file
		file, err := ioutil.ReadFile(fileName)
		check(err)
		augmentImportStatement(file, fileName, depName)
	}

	// Step 4 after adding it to import statement run go fmt on file
	cmd := exec.Command("go", "fmt", fileName)
	err := cmd.Run()
	check(err)

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
	check(err)
}

func appendDepsToBuffer(buffer []byte, depName []byte) []byte {
	for _, token := range depName {
		buffer = append(buffer, token)
	}

	return buffer
}

/*
 Deps Command Functions
*/

func readFiles(goFiles []string) {
	if len(goFiles) == 0 {
		path, err := os.Getwd()
		check(err)
		fmt.Println(path)
		fmt.Println("└── (empty)\n")
		os.Exit(3)
	}

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
			if strings.Contains(depName, "github") {
				color.Green("└──" + depName)
			} else {
				fmt.Println("└──" + depName)
			}
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
	Install Command Functions
*/

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
