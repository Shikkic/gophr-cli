package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/skeswa/gophr/common"
)

// Define Dependency Struct
type Dependency struct {
	name, version string
	installed     bool
}

// Parse Dependencies from a .go file
// TODO consider returning if the file name does not exist
func ParseDeps(fileName string) []string {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, fileName, nil, parser.ImportsOnly)
	Check(err)

	depsArray := make([]string, len(f.Imports))
	for index, s := range f.Imports {
		depName := strings.Replace(s.Path.Value, string('"'), " ", 2)
		depName = strings.Replace(depName, " ", "", 10)
		depsArray[index] = depName
	}

	return depsArray
}

func DepExistsInList(depName string, depArray []string) bool {
	for _, currDepName := range depArray {
		if currDepName == depName {
			return true
		}
	}

	return false
}

// TODO need to look up error codes and print message
func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func FetchSearchResultsData(searchQuery string) ([]byte, error) {
	request, err := http.Get("http://gophr.dev/api/search?q=" + searchQuery)
	if err != nil {
		// TODO
		// return an error code
		return nil, err
	}
	requestData, err := ioutil.ReadAll(request.Body)
	if err != nil {
		// TODO
		// return an error code
		return nil, err
	}

	return requestData, nil
}

func PrintSearchResultPackageModels(packageModels []common.PackageDTO) {
	if len(packageModels) == 0 {
		PrintEmptySearchResults()
		return
	}

	for _, packageModel := range packageModels {
		fmt.Printf("%s \n", Magenta(packageModel.Author+"/"+packageModel.Repo))
		// TODO fetch the real download numbers
		fmt.Println("3123 Downloads")
		fmt.Println(packageModel.Description + "\n")
	}
}

func PrintEmptySearchResults() {
	fmt.Println("No results found with that query")
}

func BuildPackageModelsFromRequestData(packageModelData []byte) ([]common.PackageDTO, error) {
	var packageModels []common.PackageDTO
	err := ffjson.Unmarshal(packageModelData, &packageModels)
	if err != nil {
		// return error with error code
		return nil, err
	}

	return packageModels, nil
}

func InitSpinner() *spinner.Spinner {
	return spinner.New(spinner.CharSets[14], 100*time.Millisecond)
}

/*
// Text Color Functions
# Use these color functions with printf() to make stdout colored
*/

func Magenta(text string) string {
	magenta := color.New(color.FgMagenta).SprintFunc()
	return magenta(text)
}

func Red(text string) string {
	red := color.New(color.FgRed).SprintFunc()
	return red(text)
}

func Green(text string) string {
	green := color.New(color.FgGreen).SprintFunc()
	return green(text)
}

func Blue(text string) string {
	blue := color.New(color.FgBlue).SprintFunc()
	return blue(text)
}
