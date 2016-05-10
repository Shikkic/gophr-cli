package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/codegangsta/cli"
	"github.com/skeswa/gophr/common"
)

func RunSearchCommand(c *cli.Context) {
	spinner := InitSpinner()
	spinner.Start()
	searchQueryArg := c.Args().First()
	validateSearchQueryArg(searchQueryArg)
	searchResultsData, err := FetchSearchResultsData(searchQueryArg)
	Check(err)
	searchResultsPackages, err := BuildPackageModelsFromRequestData(searchResultsData)
	Check(err)
	spinner.Stop()
	PrintSearchResultPackageModels(searchResultsPackages)
}

func validateSearchQueryArg(searchQuery string) {
	if len(searchQuery) == 0 {
		newError := NewInvalidArgumentError("Search Query", searchQuery, 1)
		newError.PrintErrorAndExit()
		os.Exit(1)
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
