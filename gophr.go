package main

import (
	"bufio"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/codegangsta/cli"
	"github.com/fatih/color"
	//"io"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	//"reflect"
	"strings"
	"time"
)

// Define Constants
const readBufferSize int = 7
const basicSkeleton string = "package main\n\nimport (\n\t\"fmt\"\n)\n\nfunc main () {\n\tfmt.Println(\"hello world!\")\n}"

// Define Dependency Struct
type Dependency struct {
	name, version string
	installed     bool
}

// TODO Consider breaking up each command into seperate go file
func main() {
	app := cli.NewApp()
	app.Name = "gophr"
	app.Usage = "A good go package manager"
	app.Commands = []cli.Command{
		{
			Name:    "deps",
			Aliases: []string{"d"},
			Usage:   "List dependencies of a go file or folder",
			Action: func(c *cli.Context) {
				fileName := c.Args().First()
				switch {
				case len(fileName) != 0:
					readFile(fileName)
				default:
					fls, err := filepath.Glob("*.go")
					check(err)
					readFiles(fls)
				}
			},
		},
		{
			Name:    "install",
			Aliases: []string{"install deps"},
			Usage:   "Install dependency",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "a",
					Value: "all",
					Usage: "install dependency to all files",
				},
			},
			Action: func(c *cli.Context) {
				var depName string
				var fileName string

				// TODO Consider using -a or --all flag to re-install all dependencies
				if c.NArg() == 0 {
					// TODO move these into functions
					red := color.New(color.FgRed).SprintFunc()
					magenta := color.New(color.FgMagenta).SprintFunc()
					fmt.Printf("%s gophr %s %s not run with a package name\n", red("✗"), red("ERROR"), magenta("install"))
					fmt.Printf("run %s for more help\n", magenta("gophr install -h"))
					os.Exit(3)
				}

				// TODO check if type string with reflect!
				if c.NArg() < 2 {
					// TODO move these into functions
					red := color.New(color.FgRed).SprintFunc()
					magenta := color.New(color.FgMagenta).SprintFunc()
					fmt.Printf("%s gophr %s %s not run with a file name\n", red("✗"), red("ERROR"), magenta("install"))
					fmt.Printf("run %s for more help\n", magenta("gophr install -h"))
					os.Exit(3)
				}

				if c.NArg() > 0 {
					depName = c.Args()[0]
				}

				// TODO consider tabbing for arg if not present
				if c.NArg() > 1 {
					fileName = c.Args()[1]
				}

				runGoGetCommand(depName, fileName)
			},
		},
		{
			Name:    "uninstall",
			Aliases: []string{"uninstall dep"},
			Usage:   "Uninstall dependency",
			Action: func(c *cli.Context) {
				var depName string
				var fileName string

				// TODO Consider using -a or --all flag to re-install all dependencies
				if c.NArg() == 0 {
					// TODO move these into functions
					red := color.New(color.FgRed).SprintFunc()
					magenta := color.New(color.FgMagenta).SprintFunc()
					fmt.Printf("%s gophr %s %s not run with a package name\n", red("✗"), red("ERROR"), magenta("uninstall"))
					fmt.Printf("run %s for more help\n", magenta("gophr uninstall -h"))
					os.Exit(3)
				}

				// TODO check if type string with reflect
				if c.NArg() < 2 {
					// TODO move these into functions
					red := color.New(color.FgRed).SprintFunc()
					magenta := color.New(color.FgMagenta).SprintFunc()
					fmt.Printf("%s gophr %s %s not run with a file name\n", red("✗"), red("ERROR"), magenta("uninstall"))
					fmt.Printf("run %s for more help\n", magenta("gophr uninstall -h"))
					os.Exit(3)
				}

				if c.NArg() > 0 {
					depName = c.Args()[0]
				}

				// TODO consider tabbing for arg if not present
				if c.NArg() > 1 {
					fileName = c.Args()[1]
				}

				runUninstallCommand(depName, fileName)
			},
		},
		{
			Name:    "init",
			Aliases: []string{"new"},
			Usage:   "initialize new project",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "lib",
					Value: "library",
					Usage: "create a basic library",
				},
			},
			Action: func(c *cli.Context) {
				var repoAuthor string
				var projectName string

				// First check if GOPATH is set, err if not
				goPath := os.Getenv("GOPATH")
				if len(goPath) < 0 {
					red := color.New(color.FgRed).SprintFunc()
					magenta := color.New(color.FgMagenta).SprintFunc()
					fmt.Printf("%s gophr %s %s $GOPATH not set\n", red("✗"), red("ERROR"), magenta("init"))
					os.Exit(3)
				}

				// TODO consider tabbing for arg if not present
				if c.NArg() == 0 {
					reader := bufio.NewReader(os.Stdin)
					fmt.Print("Repo Author: ")
					repoAuthorInput, _ := reader.ReadString('\n')
					repoAuthor = strings.Replace(repoAuthorInput, string('\n'), "", 1)
					fmt.Print("Project Name: ")
					projectNameInput, _ := reader.ReadString('\n')
					projectName = strings.Replace(projectNameInput, string('\n'), "", 1)
				}

				initPath := filepath.Join(goPath, "src", "github.com", repoAuthor, projectName)
				os.MkdirAll(initPath, 0777)

				// Now we need to glob to make sure a file name like that doesn't already exists
				fls, err := filepath.Glob(initPath + "/*.go")
				check(err)
				fmt.Println(fls)

				if len(fls) > 0 {
					// check if the .go file names match your project name
					for _, fileName := range fls {
						if fileName == initPath+"/"+projectName+".go" {
							red := color.New(color.FgRed).SprintFunc()
							magenta := color.New(color.FgMagenta).SprintFunc()
							fmt.Printf("%s gophr %s %s file with that name already exists\n", red("✗"), red("ERROR"), magenta("init"))
							os.Exit(3)
						}
					}
				}

				newFileBuffer := []byte(basicSkeleton)
				err = ioutil.WriteFile(filepath.Join(initPath, projectName)+".go", newFileBuffer, 0644)
				check(err)
			},
		},
		{
			Name:    "migrate",
			Aliases: []string{"convert"},
			Usage:   "Migrate go package to use gophr.pm/<REPO_NAME>",
			Action: func(c *cli.Context) {
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

				RunMigrateCommand(fileName)
			},
		},
	}
	app.Run(os.Args)
}

func runUninstallCommand(depName string, fileName string) {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()
	depsArray := parseDeps(fileName)

	// If a dep does not exist in the import statemtn, if it does not exist then throw an error
	if depExistsInList(depName, depsArray) == false {
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
	check(err)

	depsArray = parseDeps(fileName)
	if depExistsInList(depName, depsArray) == false {
		magenta := color.New(color.FgMagenta).SprintFunc()
		s.Stop()
		// TODO turn this check mark green
		fmt.Printf("✓ %s was successfully uninstalled from %s\n", magenta("'"+depName+"'"), magenta(fileName))
		os.Exit(3)
	}
}

func depExistsInList(depName string, depArray []string) bool {
	for _, currDepName := range depArray {
		if currDepName == depName {
			return true
		}
	}

	return false
}

/*
 Deps Command Functions
*/

// TODO consider renaming to more specific
func readFiles(goFiles []string) {
	if len(goFiles) == 0 {
		path, err := os.Getwd()
		check(err)
		fmt.Println(path)
		fmt.Println("└── (empty)\n")
		os.Exit(3)
	}

	for _, goFile := range goFiles {
		readFile(goFile)
	}
}

func readFile(goFilePath string) {
	depsArray := parseDeps(goFilePath)
	// TODO Check to determine all github which packages are installed for
	// use map to distinguish
	printDeps(depsArray, goFilePath)
}

func printDeps(depsArray []string, goFileName string) {
	fmt.Print("Go Dependecies for ")
	color.Blue(goFileName)

	for index, depName := range depsArray {
		if index == (len(depsArray) - 1) {
			if strings.Contains(depName, "github") || strings.Contains(depName, "gophr.dev") {
				color.Green("└── " + depName + "\n")
			} else {
				fmt.Println("└── " + depName + "\n")
			}
		} else {
			if strings.Contains(depName, "github") || strings.Contains(depName, "gophr.dev") {
				color.Green("├─┬ " + depName)
			} else {
				fmt.Println("├─┬ " + depName)
			}
		}
	}
}

/*
	Install Command Functions
*/

func runGoGetCommand(depName string, fileName string) {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()

	// Step 1 Determine if we need to download and install dependencies of the folder or for specified dependency
	// TODO instead use a flag to determine if it should install to all packages
	if len(depName) == 0 {
		fls, err := filepath.Glob("*.go")
		check(err)

		if len(fls) == 0 {
			red := color.New(color.FgRed).SprintFunc()
			magenta := color.New(color.FgMagenta).SprintFunc()
			s.Stop()
			fmt.Printf("gophr %s %s not run in go a package\n", red("ERROR"), magenta("install"))
			os.Exit(3)
		}

		depName = "./"
	}

	// Step 2 run get command if
	if strings.Contains(depName, "github.com") {
		cmd := exec.Command("go", "get", depName)
		//var out bytes.Buffer
		//cmd.Stdout = &out
		err := cmd.Run()
		check(err)
	}

	// Step 3 if command was successful, append to file
	if len(fileName) > 0 {
		// add to file
		file, err := ioutil.ReadFile(fileName)
		check(err)
		augmentImportStatement(file, fileName, depName)
	}

	// Step 4 after adding it to import statement run go fmt on file
	cmd := exec.Command("go", "fmt", fileName)
	err := cmd.Run()
	check(err)

	// Check if exits in file
	depsArray := parseDeps(fileName)
	if depExistsInList(depName, depsArray) == true {
		magenta := color.New(color.FgMagenta).SprintFunc()
		s.Stop()
		fmt.Printf("✓ %s was successfully installed into %s\n", magenta("'"+depName+"'"), magenta(fileName))
		os.Exit(3)
	} else {
		//PANIC ITS NOT THERE
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
	check(err)
}

func appendDepsToBuffer(buffer []byte, depName []byte) []byte {
	for _, token := range depName {
		buffer = append(buffer, token)
	}

	return buffer
}

/*
Helper Functions
*/

// Parse Dependencies from a .go file
func parseDeps(fileName string) []string {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, fileName, nil, parser.ImportsOnly)
	check(err)

	depsArray := make([]string, len(f.Imports))
	for index, s := range f.Imports {
		depName := strings.Replace(s.Path.Value, string('"'), " ", 2)
		depName = strings.Replace(depName, " ", "", 10)
		depsArray[index] = depName
	}

	return depsArray
}

// Returns an array of built dependency structs from an array of dep names.
func buildDependencyStructs(depNames []string) {

}

// Return a map of dependencies that have the attributes installed or missing
func validateDepIsInstalled(depName string) {

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
