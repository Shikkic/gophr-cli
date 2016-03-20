package main

import (
	//"bytes"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/codegangsta/cli"
	"github.com/fatih/color"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
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
				// TODO check if deps are present in the go files AND if they're installed or not
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
				// TODO consider tabbing
				depName := c.Args().First()
				runGoGetCommand(depName)
			},
		},
	}
	app.Run(os.Args)
}

func runGoGetCommand(depName string) {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()

	if len(depName) == 0 {
		fls, err := filepath.Glob("*.go")
		check(err)

		if len(fls) == 0 {
			red := color.New(color.FgRed).SprintFunc()
			magenta := color.New(color.FgMagenta).SprintFunc()
			s.Stop()
			fmt.Printf("gophr %s %s not run in go a package\n", red("ERROR"), magenta("install"))
			os.Exit(3)
		} else {
			depName = "./"
		}
	}

	cmd := exec.Command("go", "get", depName)
	//var out bytes.Buffer
	//cmd.Stdout = &out
	err := cmd.Run()
	check(err)
	s.Stop()
	//fmt.Printf("%q", out.String())
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
