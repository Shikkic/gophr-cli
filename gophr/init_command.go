package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/fatih/color"
)

const basicSkeleton string = "package main\n\nimport (\n\t\"fmt\"\n)\n\nfunc main () {\n\tfmt.Println(\"hello world!\")\n}"

// TODO TOTAL REFACTOR
func RunInitCommand(c *cli.Context) {
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

	initPath := filepath.Join(goPath, "src", "github.com", repoAuthor, projectName)
	os.MkdirAll(initPath, 0777)

	// Now we need to glob to make sure a file name like that doesn't already exists
	fls, err := filepath.Glob(initPath + "/*.go")
	Check(err)

	if len(fls) > 0 {
		// check if the .go file names match your project name
		for _, fileName := range fls {
			if fileName == initPath+"/"+projectName+".go" {
				red := color.New(color.FgRed).SprintFunc()
				magenta := color.New(color.FgMagenta).SprintFunc()
				fmt.Printf("%s gophr %s %s file with that name already exists\n", red("✗"), red("ERROR"), magenta("init"))
				os.Exit(3)
			}
		}
	} else {
		// TODO throw error or gracefully exit
	}

	newFileBuffer := []byte(basicSkeleton)
	err = ioutil.WriteFile(filepath.Join(initPath, projectName)+".go", newFileBuffer, 0644)
	Check(err)
}
