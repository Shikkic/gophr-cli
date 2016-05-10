package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/codegangsta/cli"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/skeswa/gophr/common"
)

// TODO COMPLETE REFACTOR OF RUN LOCK COMMAND
func RunLockCommand(c *cli.Context) {
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

	// Identify the file actually exists
	file, err := os.Open("./" + fileName)
	if err != nil {
		// TODO fix this up
		fmt.Printf("%s %s %s %s \n", Red("✗"), Magenta("GOPHR LOCK"), Red("could not open"), Magenta(fileName))
		os.Exit(3)
	}
	defer file.Close()

	// Retreive a list of dependencies in file
	fmt.Printf("Looing for unversioned packages in %s\n\n", Blue(fileName))
	packageNames := ParseDeps(fileName)
	githubPackageURLs := filterPackageURLsForGithubURLs(packageNames)

	// Verify there are packages eligible to be versioned
	if len(githubPackageURLs) == 0 {
		fmt.Println("No unversioned github urls found in this go file")
	}

	// Build versioned package URLs
	var versionedPackageURLs [][]byte
	if c.Bool("latest") == true {
		versionedPackageURLs = versionPackageURLsLatest(githubPackageURLs)
	} else {
		// Version list of github packages
		versionedPackageURLs = versionPackageURLs(githubPackageURLs)
	}

	// Insert new packages into import statements
	replaceVersionedPackages(file, fileName, versionedPackageURLs)

	// Run go fmt on file
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Start()
	cmd := exec.Command("go", "fmt", fileName)
	err = cmd.Run()
	s.Stop()
	Check(err)

	// Run go get and install the packages
	s.Start()
	fmt.Println("preparing to run go get")
	cmd = exec.Command("go", "get", "--insecure")
	err = cmd.Run()
	s.Stop()
	Check(err)

	// Print out new packages for go file(s)
	PrintDepsFromFileName(fileName)
}

func replaceVersionedPackages(file io.Reader, fileName string, versionedPackages [][]byte) {
	// If a dep exist begin process of removing it from the import statement
	foundImport := false
	newFileBuffer := make([]byte, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fileLine := scanner.Text() + "\n"
		byteBuffer := scanner.Bytes()

		if foundImport == true {
			if fileLine == ")\n" {
				foundImport = false
			} else {
				// TODO ask for a version number
				if strings.Contains(fileLine, "github.com") {
					depName := strings.Replace(scanner.Text(), string('"'), "", 2)
					depName = strings.Replace(depName, string('\t'), "", 1)
					for _, packageURL := range versionedPackages {
						depNameTokens := strings.Split(depName, "/")
						author := depNameTokens[1]
						repo := depNameTokens[2]
						if strings.Contains(string(packageURL), author+"/"+repo) {
							byteBuffer = packageURL
						}
					}
				}
			}
		} else if fileLine == "import (\n" {
			foundImport = true
		}

		byteBuffer = append(byteBuffer, byte('\n'))
		for _, token := range byteBuffer {
			newFileBuffer = append(newFileBuffer, token)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	err := ioutil.WriteFile("./"+fileName, newFileBuffer, 0644)
	Check(err)
}

func versionPackageURLsLatest(packageURLs []string) [][]byte {
	versionedPackageURLs := [][]byte{}
	for _, packageURL := range packageURLs {
		versionedPackageURL := buildVersionedPackageURLLatest(packageURL)
		versionedPackageURLs = append(versionedPackageURLs, versionedPackageURL)
	}

	return versionedPackageURLs
}

func versionPackageURLs(packageURLs []string) [][]byte {
	versionedPackageURLs := [][]byte{}
	for _, packageURL := range packageURLs {
		versionedPackageURL := buildVersionedPackageURL(packageURL)
		versionedPackageURLs = append(versionedPackageURLs, versionedPackageURL)
	}

	return versionedPackageURLs
}

func buildVersionedPackageURL(packageURL string) []byte {
	reader := bufio.NewReader(os.Stdin)
	versionList := retrieveVersionList(packageURL)
	fmt.Printf("%s \n", Magenta(packageURL))
	fmt.Println(versionList)
	fmt.Print("version number: ")
	versionNumberInput, _ := reader.ReadString('\n')
	versionNumber := strings.Replace(versionNumberInput, string('\n'), "", 1)
	// TODO make this a function
	byteBuffer := []byte(strings.Replace(packageURL, "github.com/", "gophr.dev/", 1) + "@" + versionNumber + "\"")

	return byteBuffer
}

func buildVersionedPackageURLLatest(packageURL string) []byte {
	s := spinner.New(spinner.CharSets[14], 50*time.Millisecond)
	s.Start()
	versionLatest := retrieveVersionLatest(packageURL)
	s.Stop()

	fmt.Printf("Found %s \n", Magenta(packageURL))
	fmt.Printf("latest version %s\n", Blue(versionLatest))
	url := strings.Replace(packageURL, "github.com/", "\"gophr.dev/", 1) + "@" + versionLatest + "\""
	fmt.Printf("%s", Green("✓ package was successfully versioned at "+url+"\n"))
	fmt.Println("")
	// TODO make this a function
	byteBuffer := []byte(url)

	return byteBuffer
}

func filterPackageURLsForGithubURLs(packageURLs []string) []string {
	githubPackageURLs := []string{}
	for _, packageURL := range packageURLs {
		if strings.HasPrefix(packageURL, "github.com") {
			githubPackageURLs = append(githubPackageURLs, packageURL)
		}
	}

	return githubPackageURLs
}

func retrieveVersionList(packageName string) []string {
	packageNamesArray := strings.Split(packageName, "/")
	author := packageNamesArray[1]
	repo := packageNamesArray[2]

	url := "http://gophr.dev/api/" + author + "/" + repo + "/versions"
	res, err := http.Get(url)
	Check(err)
	data, err := ioutil.ReadAll(res.Body)
	Check(err)

	// TODO This could be abstracted
	var packageModels []common.VersionDTO

	buildVersionDTO(data, &packageModels)
	versionList := []string{}
	for _, versions := range packageModels {
		versionList = append(versionList, versions.Value)
	}

	return versionList
}

func retrieveVersionLatest(packageName string) string {
	packageNamesArray := strings.Split(packageName, "/")
	author := packageNamesArray[1]
	repo := packageNamesArray[2]

	url := "http://gophr.dev/api/" + author + "/" + repo + "/versions/latest"
	res, err := http.Get(url)
	Check(err)
	data, err := ioutil.ReadAll(res.Body)
	Check(err)

	// TODO This could be abstracted
	var packageModels common.VersionDTO
	err = ffjson.Unmarshal(data, &packageModels)

	return packageModels.Value
}

func buildVersionDTO(data []byte, versionStruct *[]common.VersionDTO) {
	err := ffjson.Unmarshal(data, &versionStruct)
	Check(err)
}
