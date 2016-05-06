package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/codegangsta/cli"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/skeswa/gophr/common"
)

func RunLockCommand(fileName string, c *cli.Context) {
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

	// Print out new packages for go file(s)
	ReadFile(fileName)
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
						packageURLTokens := strings.Split(string(packageURL), "/")
						author := packageURLTokens[1]
						repo := packageURLTokens[2]
						fmt.Println(string(packageURL))
						if strings.Contains(depName, author+"/"+repo) {
							byteBuffer = packageURL
							fmt.Println(string(byteBuffer))
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

////////////////
// DEPRICATED //
////////////////
/*
func RunLock(fileName string) {
	// If a dep exist begin process of removing it from the import statement
	file, err := os.Open("./" + fileName)
	foundImport := false
	newFileBuffer := make([]byte, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

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
					Magenta := color.New(color.FgMagenta).SprintFunc()
					reader := bufio.NewReader(os.Stdin)
					depName := strings.Replace(scanner.Text(), string('"'), "", 2)
					depName = strings.Replace(depName, string('\t'), "", 1)
					versionList := retrieveVersionList(depName)
					fmt.Printf("%s \n", Magenta(depName))
					fmt.Println(versionList)
					fmt.Print("version number:")
					versionNumberInput, _ := reader.ReadString('\n')
					versionNumber := strings.Replace(versionNumberInput, string('\n'), "", 1)
					byteBuffer = []byte("\t\"" + strings.Replace(depName, "github.com/", "gophr.dev/", 1) + "@" + versionNumber + "\"")
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

	err = ioutil.WriteFile("./"+fileName, newFileBuffer, 0644)
	Check(err)
}
*/
