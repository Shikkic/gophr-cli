package main

import (
	"bufio"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func runUninstallCommand(depName string, fileName string) {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()
	depsArray := ParseDeps(fileName)

	// If a dep does not exist in the import statemtn, if it does not exist then throw an error
	if DepExistsInList(depName, depsArray) == false {
		red := color.New(color.FgRed).SprintFunc()
		magenta := color.New(color.FgMagenta).SprintFunc()
		s.Stop()
		fmt.Printf("%s gophr %s %s package %s not present in %s\n", red("✗"), red("ERROR"), magenta("uninstall"), magenta("'"+depName+"'"), magenta(fileName))
		os.Exit(3)
	}

	// If a dep exist begin process of removing it from the import statement
	file, err := os.Open("./" + fileName)
	newFileBuffer := make([]byte, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fileLine := scanner.Text() + "\n"
		if fileLine != "\t\""+depName+"\"\n" {
			byteBuffer := scanner.Bytes()
			byteBuffer = append(byteBuffer, byte('\n'))
			for _, token := range byteBuffer {
				newFileBuffer = append(newFileBuffer, token)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("./"+fileName, newFileBuffer, 0644)
	Check(err)

	depsArray = ParseDeps(fileName)
	if DepExistsInList(depName, depsArray) == false {
		magenta := color.New(color.FgMagenta).SprintFunc()
		s.Stop()
		// TODO turn this check mark green
		fmt.Printf("✓ %s was successfully uninstalled from %s\n", magenta("'"+depName+"'"), magenta(fileName))
		os.Exit(3)
	}
}
