package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/shikkic/gophr-cli/gophr/common"
	"github.com/codegangsta/cli"
	"github.com/fatih/color"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/skeswa/gophr/common"
)

func main() {
	app := cli.NewApp()
	app.Name = "gophr"
	app.Usage = "An end-to-end package management solution for Go"
	app.Commands = []cli.Command{
		{
			Name:    "search",
			Aliases: []string{"d"},
			Usage:   "Search gophr dependency",
			Action: func(c *cli.Context) {
				searchQueryArg := c.Args().First()

				// Check for searchQueryArg
				if len(searchQueryArg) == 0 {
					fmt.Println("ERROR no query argument given")
					os.Exit(3)
				}

				// abstract this into gophr request lib
				res, err := http.Get("http://gophr.dev/api/search?q=" + searchQueryArg)
				data, err := ioutil.ReadAll(res.Body)
				magenta := color.New(color.FgMagenta).SprintFunc()

				if err != nil {
					fmt.Println(err)
					os.Exit(3)
				}

				var packageModels []common.PackageDTO
				err = ffjson.Unmarshal(data, &packageModels)
				Check(err)

				// abstract this into print search packagePrint
				for _, packageModel := range packageModels {
					fmt.Printf("%s \n", magenta(packageModel.Author+"/"+packageModel.Repo))
					fmt.Println("3123 Downloads")
					fmt.Println(packageModel.Description + "\n")
				}
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
				// TODO SHOULD INCLUDE BASH COMPLETION FOR CURRENT DEPENDENCIES IN FILE NAME
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
