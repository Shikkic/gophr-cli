package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/skeswa/gophr/common"
)

func RunInstallCommand(depName string, fileName string) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Start()
	// Step 1 Determine if we need to download and install dependencies of the folder or for specified dependency
	// TODO instead use a flag to determine if it should install to all packages
	if len(depName) == 0 {
		fls, err := filepath.Glob("*.go")
		Check(err)

		if len(fls) == 0 {
			red := color.New(color.FgRed).SprintFunc()
			magenta := color.New(color.FgMagenta).SprintFunc()
			s.Stop()
			fmt.Printf("gophr %s %s not run in go a package\n", red("ERROR"), magenta("install"))
			os.Exit(3)
		}

		depName = "./"
	}

	// Do you want
	magenta := color.New(color.FgMagenta).SprintFunc()
	s.Stop()
	fmt.Printf("gophr %s version lock %s\n", magenta("INSTALL"), magenta(depName))
	fmt.Println("y/n?")
	reader := bufio.NewReader(os.Stdin)

	input, _ := reader.ReadString('\n')
	versionLock := strings.Replace(input, string('\n'), "", 2)

	s.Start()
	var depVersion string
	if versionLock == "y" || versionLock == "yes" {
		// Print known versions
		res, err := http.Get("http://gophr.dev/api/" + depName + "/versions")
		Check(err)
		data, err := ioutil.ReadAll(res.Body)
		Check(err)
		s.Stop()
		var things []common.VersionDTO
		err = ffjson.Unmarshal(data, &things)
		Check(err)
		if len(things) == 0 {
			fmt.Println("No version branches detected, fetching master branch SHA")
			s.Start()
			res, err := http.Get("http://gophr.dev/api/" + depName + "/versions/latest")
			data, err := ioutil.ReadAll(res.Body)
			var latestVersion common.VersionDTO
			err = ffjson.Unmarshal(data, &latestVersion)
			Check(err)
			s.Stop()
			fmt.Println("Master branch SHA: " + latestVersion.Value)
			depVersion = "@" + latestVersion.Value
		} else {
			// TODO ask if you want latest hash?
			fmt.Println("Known versions: ")
			for _, lol := range things {
				fmt.Println(lol.Value)
			}
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			depVersion = "@" + strings.Replace(input, string('\n'), "", 2)
			// TODO verify version is one of the ones that exists
		}
	}

	s.Start()
	depName = "gophr.dev/" + depName + depVersion
	cmd := exec.Command("go", "get", "--insecure", depName)
	err := cmd.Run()
	Check(err)

	// Step 3 if command was successful, append to file
	if len(fileName) > 0 {
		file, errz := ioutil.ReadFile(fileName)
		Check(errz)
		augmentImportStatement(file, fileName, depName)
	}

	// Step 4 after adding it to import statement run go fmt on file
	cmd = exec.Command("go", "fmt", fileName)
	err = cmd.Run()
	Check(err)

	// Check if exits in file
	depsArray := ParseDeps(fileName)
	if DepExistsInList(depName, depsArray) == true {
		magenta := color.New(color.FgMagenta).SprintFunc()
		s.Stop()
		fmt.Printf("✓ %s was successfully installed into %s\n", magenta("'"+depName+"'"), magenta(fileName))
		os.Exit(3)
	} else {
		red := color.New(color.FgRed).SprintFunc()
		magenta := color.New(color.FgMagenta).SprintFunc()
		s.Stop()
		fmt.Printf("x %s failed to install %s\n", red("ERROR"), magenta("'"+depName+"'"))
		os.Exit(3)
	}

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
	Check(err)
}
