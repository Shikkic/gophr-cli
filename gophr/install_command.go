package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/skeswa/gophr/common"
)

func RunInstallCommand(c *cli.Context) error {
	depName, err := getFirstArgDepName(c)
	if err != nil {
		return err
	}

	fileName, err := getSecondArgFileName(c)
	if err != nil {
		return err
	}

	depVersion, err := getDepVersionFromUser(depName)
	depGophrURL := BuildVersionedGophrDepURL(depName, depVersion)
	err = RunGoGetDep(depGophrURL)
	if err != nil {
		// TODO special error here
		return err
	}

	file, err := ioutil.ReadFile(fileName)
	Check(err)
	augmentGoFileImportStatement(file, fileName, depGophrURL)
	err = RunGoFMTOnFileName(fileName)
	if err != nil {
		return err
	}

	err = ValidateDepWasInstalledIntoFileName(depName, fileName)
	if err != nil {
		return err
	}

	return nil
}

// TODO deps_helper
func BuildVersionedGophrDepURL(depName string, depVersion string) string {
	url := GetGophrBaseURL() + "depName" + "@" + "depVersion"
	return url
}

// TODO move this to HELPER
func RunGoGetDep(depURL string) error {
	cmd := exec.Command("go", "get", "--insecure", depURL)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// TODO move this to HELPER
func RunGoFMTOnFileName(fileName string) error {
	cmd := exec.Command("go", "fmt", fileName)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// TODO move this to deps_helper
func ValidateDepWasInstalledIntoFileName(depName string, fileName string) error {
	file, err := OpenASTFilePointerFromFileName(fileName)
	if err != nil {
		return nil
	}

	depURLs := ParseDepURLsFromFile(file)
	// TODO seperate this into seperate function
	if DepExistsInList(depName, depURLs) == true {
		fmt.Printf("✓ %s was successfully installed into %s\n", Magenta("'"+depName+"'"), Magenta(fileName))
		return nil
	} else {
		// TODO CREATE NEW UNIQUE ERROR HERE
		fmt.Printf("x %s failed to install %s\n", Red("ERROR"), Magenta("'"+depName+"'"))
		return nil
	}
}

func getDepVersionFromUser(depName string) (string, error) {
	depVersion, err := promptUserForDepVersion(depName)
	if err != nil {
		return "", err
	}

	return depVersion, nil
}

// TODO move this to generic helper file
func GetUserInput() string {
	bufferReader := bufio.NewReader(os.Stdin)
	userInput, _ := bufferReader.ReadString('\n')
	userInput = strings.Replace(userInput, string('\n'), "", 2)

	return userInput
}

// TODO move this to deps_helper
func FetchVersionsForDep(depName string) []common.VersionDTO {
	requestURL := GetGophrBaseURL() + "/api" + depName + "/versions"
	request, err := http.Get(requestURL)
	Check(err)
	data, err := ioutil.ReadAll(request.Body)
	Check(err)
	var things []common.VersionDTO
	err = ffjson.Unmarshal(data, &things)
	Check(err)

	return things
}

// TODO move this to deps_helper
func FetchLatestVersionForDep(depName string) common.VersionDTO {
	fmt.Println("No version branches detected, fetching master branch SHA")
	res, err := http.Get("http://gophr.dev/api/" + depName + "/versions/latest")
	data, err := ioutil.ReadAll(res.Body)
	var latestVersion common.VersionDTO
	err = ffjson.Unmarshal(data, &latestVersion)
	Check(err)

	return latestVersion
}

func promptUserForDepVersion(depName string) (string, error) {
	// TODO create function
	fmt.Printf("gophr %s version lock %s\n", Magenta("INSTALL"), Magenta(depName))
	fmt.Println("y/n?")
	versionLockDep := GetUserInput()

	var depVersion string
	if versionLockDep == "y" || versionLockDep == "yes" {
		depVersions := FetchVersionsForDep(depName)

		if len(depVersions) == 0 {
			latestVersion := FetchLatestVersionForDep(depName)
			fmt.Println("Master branch SHA: " + latestVersion.Value)
			depVersion = latestVersion.Value
		} else {
			// TODO ask if you want latest hash?
			fmt.Println("Known versions: ")
			for _, depVersion := range depVersions {
				fmt.Println(depVersion.Value)
			}
			depVersion = GetUserInput()
			// TODO verify version is one of the ones that exists
		}
	} else {
		return "", nil
	}

	return depVersion, nil
}

func getFirstArgDepName(c *cli.Context) (string, error) {
	err := validateFirstArgDepNameExists(c)
	if err != nil {
		return "", err
	}
	depName := c.Args()[0]

	return depName, nil
}

func validateFirstArgDepNameExists(c *cli.Context) error {
	if c.NArg() == 0 {
		// TODO create new error type for this
		fmt.Printf("%s gophr %s %s not run with a package name\n", Red("✗"), Red("ERROR"), Magenta("install"))
		fmt.Printf("run %s for more help\n", Magenta("gophr install -h"))
		return nil
	}

	return nil
}

func getSecondArgFileName(c *cli.Context) (string, error) {
	err := validateSecondArgFileNameExists(c)
	if err != nil {
		return "", err
	}
	fileName := c.Args()[1]

	// TODO validate that the files is indeed a go file and can be opened

	return fileName, nil
}

func validateSecondArgFileNameExists(c *cli.Context) error {
	if c.NArg() < 2 {
		// TODO create new error type for sthis
		fmt.Printf("%s gophr %s %s not run with a file name\n", Red("✗"), Red("ERROR"), Magenta("install"))
		fmt.Printf("run %s for more help\n", Magenta("gophr install -h"))
		return nil
	}

	return nil
}

func augmentGoFileImportStatement(file []byte, fileName string, depName string) {
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

func appendDepsToBuffer(buffer []byte, depName []byte) []byte {
	for _, token := range depName {
		buffer = append(buffer, token)
	}

	return buffer
}
