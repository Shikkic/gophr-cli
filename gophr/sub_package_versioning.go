package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/briandowns/spinner"
	"github.com/codegangsta/cli"
)

// TODO consider taking a struct of
// depName
// version
// function to return filepath

// TODO right now it takes the full depName+version as a path
func RunSubVersioningCommand(c *cli.Context) error {
	searchDir := "/Users/shikkic/dev/go/src/gophr.dev/codegangsta/cli@2.0.0"
	// Step 1 retrieve all the go paths from a dir
	goFilePaths, err := BuildGoFilePathsFromDir(searchDir)
	// Step 2 For each file in goFilePaths run a lock command on each file
	for _, goFilePath := range goFilePaths {
		// Identify the file actually exists
		fileName := goFilePath
		file, err := os.Open(fileName)
		if err != nil {
			// TODO fix this up
			fmt.Printf("%s %s %s %s \n", Red("✗"), Magenta("GOPHR LOCK"), Red("could not open"), Magenta(fileName))
			os.Exit(3)
		}
		defer file.Close()

		// Retreive a list of dependencies in file
		fmt.Printf("Looking for unversioned packages in %s\n\n", Blue(fileName))
		packageNames := ParseDeps(fileName)
		fmt.Println("filtering")
		githubPackageURLs := filterPackageURLsForGithubURLs(packageNames)

		fmt.Println("STUCK")

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

		fmt.Println(versionedPackageURLs)

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

	return err
}

func BuildGoFilePathsFromDir(searchDir string) ([]string, error) {
	fileList := []string{}
	err := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return nil // but continue walking elsewhere
		}
		if f.IsDir() {
			return nil // not a file.  ignore.
		}
		matched, err := filepath.Match("*.go", f.Name())
		if err != nil {
			return err // this is fatal.
		}
		if matched {
			fileList = append(fileList, path)
		}
		return nil
	})

	return fileList, err
}
