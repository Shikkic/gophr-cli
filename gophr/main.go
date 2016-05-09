package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "gophr"
	app.Usage = "An end-to-end package management solution for Go"
	app.Commands = []cli.Command{
		{
			Name:    "search",
			Aliases: []string{"s"},
			Usage:   "Search go packages on gophr.pm",
			Action: func(c *cli.Context) {
				spinner := InitSpinner()
				spinner.Start()
				searchQueryArg := c.Args().First()

				// TODO create validation and error handling helper
				if len(searchQueryArg) == 0 {
					fmt.Println("ERROR no query argument given")
					os.Exit(1)
				}

				searchResultsData, err := FetchSearchResultsData(searchQueryArg)

				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				searchResultsPackages, err := BuildPackageModelsFromRequestData(searchResultsData)

				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				spinner.Stop()
				PrintSearchResultPackageModels(searchResultsPackages)
			},
		},
		{
			Name:    "deps",
			Aliases: []string{"d"},
			Usage:   "List dependencies of a go file or folder",
			Action: func(c *cli.Context) {
				fileNameArg := c.Args().First()
				switch {
				case len(fileNameArg) != 0:
					// TODO Rename this
					ReadFile(fileNameArg)
				default:
					fls, err := filepath.Glob("*.go")
					Check(err)
					// TODO Rename this
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
					fmt.Printf("%s gophr %s %s not run with a package name\n", Red("✗"), Red("ERROR"), Magenta("install"))
					fmt.Printf("run %s for more help\n", Magenta("gophr install -h"))
					os.Exit(3)
				}

				// TODO check if type string with reflect!
				if c.NArg() < 2 {
					// TODO move these into functions
					fmt.Printf("%s gophr %s %s not run with a file name\n", Red("✗"), Red("ERROR"), Magenta("install"))
					fmt.Printf("run %s for more help\n", Magenta("gophr install -h"))
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
				// TODO SHOULD INCLUDE BASH COMPLETION FOR CURRENT DEPENDENCIES IN FILE NAME
				var depName string
				var fileName string

				// TODO Consider using -a or --all flag to re-install all dependencies
				if c.NArg() == 0 {
					// TODO move these into functions
					fmt.Printf("%s gophr %s %s not run with a package name\n", Red("✗"), Red("ERROR"), Magenta("uninstall"))
					fmt.Printf("run %s for more help\n", Magenta("gophr uninstall -h"))
					os.Exit(3)
				}

				// TODO check if type string with reflect
				if c.NArg() < 2 {
					// TODO move these into functions
					fmt.Printf("%s gophr %s %s not run with a file name\n", Red("✗"), Red("ERROR"), Magenta("uninstall"))
					fmt.Printf("run %s for more help\n", Magenta("gophr uninstall -h"))
					os.Exit(3)
				}

				if c.NArg() > 0 {
					depName = c.Args()[0]
				}

				// TODO consider tabbing for arg if not present
				if c.NArg() > 1 {
					fileName = c.Args()[1]
				}

				RunUninstallCommand(depName, fileName)
			},
		},
		{
			Name:    "init",
			Aliases: []string{"new"},
			Usage:   "initialize new project",
			Flags: []cli.Flag{
				// TODO create lib flag to generate library
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

					fmt.Printf("%s gophr %s %s $GOPATH not set\n", Red("✗"), Red("ERROR"), Magenta("init"))
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
			Name:    "lock",
			Aliases: []string{"convert"},
			Usage:   "lock a file(s) github go packages to use gophr.pm/<REPO_NAME>",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "latest",
					Usage: "lock all github go package dependencies to latest master SHA",
				},
			},
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

				RunLockCommand(fileName, c)
			},
		},
	}
	app.Run(os.Args)
}
