package main

import (
	"bufio"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/codegangsta/cli"
	"github.com/fatih/color"
	//"io"
	"io/ioutil"
	"log"
	"os"
	//"os/exec"
	"path/filepath"
	//"reflect"
	"strings"
	"time"
)

// Define Constants
// TODO move this to helper library
// doesn't need to be constant
const readBufferSize int = 7

// TODO move this to helper library
// Define Dependency Struct
type Dependency struct {
	name, version string
	installed     bool
}

// TODO Consider breaking up each command into seperate go file
func main() {
	app := cli.NewApp()
	app.Name = "gophr"
	app.Usage = "A good go package manager"
	app.Commands = []cli.Command{
		{
			Name:    "deps",
			Aliases: []string{"d"},
			Usage:   "List dependencies of a go file or folder",
			Action: func(c *cli.Context) {
				fileName := c.Args().First()
				switch {
				case len(fileName) != 0:
					// TODO Rename this
					ReadFile(fileName)
				default:
					// TODO Rename this
					fls, err := filepath.Glob("*.go")
					Check(err)
					ReadFiles(fls)
				}
			},
		},
		{
			Name:    "install",
			Aliases: []string{"install deps"},
			Usage:   "Install dependency",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "a",
					Value: "all",
					Usage: "install dependency to all files",
				},
			},
			Action: func(c *cli.Context) {
				var depName string
				var fileName string

				// TODO Consider using -a or --all flag to re-install all dependencies
				if c.NArg() == 0 {
					// TODO move these into functions
					red := color.New(color.FgRed).SprintFunc()
					magenta := color.New(color.FgMagenta).SprintFunc()
					fmt.Printf("%s gophr %s %s not run with a package name\n", red("✗"), red("ERROR"), magenta("install"))
					fmt.Printf("run %s for more help\n", magenta("gophr install -h"))
					os.Exit(3)
				}

				// TODO check if type string with reflect!
				if c.NArg() < 2 {
					// TODO move these into functions
					red := color.New(color.FgRed).SprintFunc()
					magenta := color.New(color.FgMagenta).SprintFunc()
					fmt.Printf("%s gophr %s %s not run with a file name\n", red("✗"), red("ERROR"), magenta("install"))
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

				RunInstallCommand(depName, fileName)
			},
		},
		{
			Name:    "uninstall",
			Aliases: []string{"uninstall dep"},
			Usage:   "Uninstall dependency",
			Action: func(c *cli.Context) {
				var depName string
				var fileName string

				// TODO Consider using -a or --all flag to re-install all dependencies
				if c.NArg() == 0 {
					// TODO move these into functions
					red := color.New(color.FgRed).SprintFunc()
					magenta := color.New(color.FgMagenta).SprintFunc()
					fmt.Printf("%s gophr %s %s not run with a package name\n", red("✗"), red("ERROR"), magenta("uninstall"))
					fmt.Printf("run %s for more help\n", magenta("gophr uninstall -h"))
					os.Exit(3)
				}

				// TODO check if type string with reflect
				if c.NArg() < 2 {
					// TODO move these into functions
					red := color.New(color.FgRed).SprintFunc()
					magenta := color.New(color.FgMagenta).SprintFunc()
					fmt.Printf("%s gophr %s %s not run with a file name\n", red("✗"), red("ERROR"), magenta("uninstall"))
					fmt.Printf("run %s for more help\n", magenta("gophr uninstall -h"))
					os.Exit(3)
				}

				if c.NArg() > 0 {
					depName = c.Args()[0]
				}

				// TODO consider tabbing for arg if not present
				if c.NArg() > 1 {
					fileName = c.Args()[1]
				}

				runUninstallCommand(depName, fileName)
			},
		},
		{
			Name:    "init",
			Aliases: []string{"new"},
			Usage:   "initialize new project",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "lib",
					Value: "library",
					Usage: "create a basic library",
				},
			},
			Action: func(c *cli.Context) {
				var repoAuthor string
				var projectName string

				// First check if GOPATH is set, err if not
				goPath := os.Getenv("GOPATH")
				if len(goPath) < 0 {
					red := color.New(color.FgRed).SprintFunc()
					magenta := color.New(color.FgMagenta).SprintFunc()
					fmt.Printf("%s gophr %s %s $GOPATH not set\n", red("✗"), red("ERROR"), magenta("init"))
					os.Exit(3)
				}

				// TODO consider tabbing for arg if not present
				if c.NArg() == 0 {
					reader := bufio.NewReader(os.Stdin)
					fmt.Print("Repo Author: ")
					repoAuthorInput, _ := reader.ReadString('\n')
					repoAuthor = strings.Replace(repoAuthorInput, string('\n'), "", 1)
					fmt.Print("Project Name: ")
					projectNameInput, _ := reader.ReadString('\n')
					projectName = strings.Replace(projectNameInput, string('\n'), "", 1)
				}

				RunInitCommand(goPath, repoAuthor, projectName)
			},
		},
		{
			Name:    "migrate",
			Aliases: []string{"convert"},
			Usage:   "Migrate go package to use gophr.pm/<REPO_NAME>",
			Action: func(c *cli.Context) {
				var fileName string

				// TODO consider tabbing for arg if not present
				if c.NArg() == 0 {
					reader := bufio.NewReader(os.Stdin)
					fmt.Print("File Name: ")
					fileNameInput, _ := reader.ReadString('\n')
					fileName = strings.Replace(fileNameInput, string('\n'), "", 1)
				} else {
					fileName = c.Args().First()
				}

				RunMigrateCommand(fileName)
			},
		},
	}
	app.Run(os.Args)
}

func runUninstallCommand(depName string, fileName string) {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()
	depsArray := ParseDeps(fileName)

	// If a dep does not exist in the import statemtn, if it does not exist then throw an error
	if DepExistsInList(depName, depsArray) == false {
		red := color.New(color.FgRed).SprintFunc()
		magenta := color.New(color.FgMagenta).SprintFunc()
		s.Stop()
		fmt.Printf("%s gophr %s %s package %s not present in %s\n", red("✗"), red("ERROR"), magenta("uninstall"), magenta("'"+depName+"'"), magenta(fileName))
		os.Exit(3)
	}

	// If a dep exist begin process of removing it from the import statement
	file, err := os.Open("./" + fileName)
	newFileBuffer := make([]byte, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fileLine := scanner.Text() + "\n"
		if fileLine != "\t\""+depName+"\"\n" {
			byteBuffer := scanner.Bytes()
			byteBuffer = append(byteBuffer, byte('\n'))
			for _, token := range byteBuffer {
				newFileBuffer = append(newFileBuffer, token)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("./"+fileName, newFileBuffer, 0644)
	Check(err)

	depsArray = ParseDeps(fileName)
	if DepExistsInList(depName, depsArray) == false {
		magenta := color.New(color.FgMagenta).SprintFunc()
		s.Stop()
		// TODO turn this check mark green
		fmt.Printf("✓ %s was successfully uninstalled from %s\n", magenta("'"+depName+"'"), magenta(fileName))
		os.Exit(3)
	}
}

// Returns an array of built dependency structs from an array of dep names.
func buildDependencyStructs(depNames []string) {

}

// Return a map of dependencies that have the attributes installed or missing
func validateDepIsInstalled(depName string) {

}
