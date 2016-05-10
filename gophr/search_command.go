package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/codegangsta/cli"
	"github.com/skeswa/gophr/common"
)

func RunSearchCommand(c *cli.Context) error {
	spinner := InitSpinner()
	spinner.Start()

	searchQueryArg := c.Args().First()
	err := validateSearchQueryArg(searchQueryArg)
	if err != nil {
		fmt.Println(err)
		return err
	}
	searchResultsData, err := FetchSearchResultsData(searchQueryArg)
	if err != nil {
		fmt.Println(err)
		return err
	}
	searchResultsPackages, err := BuildPackageModelsFromRequestData(searchResultsData)
	if err != nil {
		fmt.Println(err)
		return err
	}

	spinner.Stop()
	PrintSearchResultPackageModels(searchResultsPackages)

	return nil
}

func validateSearchQueryArg(searchQuery string) error {
	if len(searchQuery) == 0 {
		InvalidArgumentError := NewInvalidArgumentError("Search Query", searchQuery, 1)
		return errors.New(InvalidArgumentError.Error())
	}

	return nil
}

func FetchSearchResultsData(searchQuery string) ([]byte, error) {
	requestURL := GetGophrBaseURL() + "/api/search?q=" + searchQuery
	request, err := http.Get(requestURL)
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
