package main

import (
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
	"path/filepath"
)

func RunInitCommand(goPath string, repoAuthor string, projectName string) {
	initPath := filepath.Join(goPath, "src", "github.com", repoAuthor, projectName)
	os.MkdirAll(initPath, 0777)

	// Now we need to glob to make sure a file name like that doesn't already exists
	fls, err := filepath.Glob(initPath + "/*.go")
	check(err)

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
	}

	newFileBuffer := []byte(basicSkeleton)
	err = ioutil.WriteFile(filepath.Join(initPath, projectName)+".go", newFileBuffer, 0644)
	check(err)
}
