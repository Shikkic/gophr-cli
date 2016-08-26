package main

import (
	"os"

	"github.com/codegangsta/cli"
)

var DEV_MODE bool

func main() {
	app := cli.NewApp()
	app.Name = "gophr"
	app.Usage = "An end-to-end package management solution for Go"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "dev",
			Usage:       "enable developer mode on commands",
			Destination: &DEV_MODE,
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "search",
			Aliases: []string{"s"},
			Usage:   "Search go packages on gophr.pm",
			Action:  RunSearchCommand,
		},
		{
			Name:    "deps",
			Aliases: []string{"d"},
			Usage:   "List go packages of a specified go file or folder",
			Action:  RunDepsCommand,
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
			Action: RunInstallCommand,
		},
		{
			Name:    "uninstall",
			Aliases: []string{"uninstall dep"},
			Usage:   "Uninstall dependency",
			Action:  RunUninstallCommand,
		},
		{
			Name:    "init",
			Aliases: []string{"new"},
			Usage:   "Initialize new project",
			Flags: []cli.Flag{
				// TODO create lib flag to generate library
				cli.StringFlag{
					Name:  "lib",
					Value: "library",
					Usage: "create a basic library",
				},
			},
			Action: RunInitCommand,
		},
		{
			Name:    "lock",
			Aliases: []string{"l"},
			Usage:   "Lock a file(s) github go packages to use gophr.pm/<REPO_NAME>",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "latest",
					Usage: "Lock all github go package dependencies to latest version or master SHA",
				},
			},
			Action: RunLockCommand,
		},
		{
			Name:    "sub",
			Aliases: []string{"t"},
			Usage:   "TEST: This is a test command for the subversion process that will runing during a lock command",
			Action:  RunSubVersioningCommand,
		},
	}
	app.Run(os.Args)
}
