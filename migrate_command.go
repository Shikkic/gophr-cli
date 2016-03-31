package main

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func RunMigrateCommand(fileName string) {
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
					magenta := color.New(color.FgMagenta).SprintFunc()
					reader := bufio.NewReader(os.Stdin)
					depName := strings.Replace(scanner.Text(), string('"'), "", 2)
					depName = strings.Replace(depName, string('\t'), "", 1)
					fmt.Printf("%s version: ", magenta(depName))
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
