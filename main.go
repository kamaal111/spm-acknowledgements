package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	startTimer := time.Now()

	spmPath := getSPMPath()

	outputPath := initializeFlag("Output path", "", "output", "o")

	spmDirectoryContent, err := ioutil.ReadDir(spmPath)
	checkError(err)

	var acknowledgements []Acknowledgement

	for _, spmPackage := range spmDirectoryContent {
		if spmPackage.IsDir() {
			packagePath := appendFileToPath(spmPath, spmPackage.Name())
			packageDirectoryContent, err := ioutil.ReadDir(packagePath)
			checkError(err)

			acknowledgement := Acknowledgement{
				PackageName: spmPackage.Name(),
			}

			for _, packageFile := range packageDirectoryContent {
				if packageFile.Name() == "LICENSE" {
					licenseData, err := ioutil.ReadFile(appendFileToPath(packagePath, packageFile.Name()))
					checkError(err)

					acknowledgement.Content = string(licenseData)
					break
				}
			}

			acknowledgements = append(acknowledgements, acknowledgement)
		}
	}

	err = createJSONFile(acknowledgements, appendFileToPath(outputPath, "acknowledgements.json"))
	checkError(err)

	timeElapsed := time.Since(startTimer)
	fmt.Printf("Created acknowledgements file in %s ✨\n", timeElapsed)
}

// Acknowledgement - structure of an acknowledgement object
type Acknowledgement struct {
	PackageName string `json:"package_name,omitempty"`
	Content     string `json:"content,omitempty"`
	Version     string `json:"version,omitempty"`
	URL         string `json:"url,omitempty"`
	Author      string `json:"author,omitempty"`
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func appendFileToPath(path string, file string) string {
	if len(strings.TrimSpace(path)) < 1 {
		return file
	}

	pathRune := []rune(path)
	lastCharacter := string(pathRune[len(pathRune)-1:])
	if lastCharacter == "/" {
		return path + file
	}

	return fmt.Sprintf("%s/%s", path, file)
}

func createJSONFile(data []Acknowledgement, path string) error {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	createdFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer createdFile.Close()

	_, err = createdFile.Write(jsonBytes)
	if err != nil {
		return err
	}

	return nil
}

func getSPMPath() string {
	spmPath := os.Getenv("BUILD_DIR")
	if len(spmPath) > 0 {
		return spmPath + "/../../SourcePackages/checkouts"
	}
	spmPathFlag := initializeFlag("SPM path", "", "spm", "s")
	if len(spmPathFlag) < 1 {
		log.Fatalln(errors.New("please provide the SPM path with -s or -spm"))
	}
	return spmPathFlag

}

func initializeFlag(usage string, flagDefault string, longVariable string, shortVariable string) string {
	var value string
	flag.StringVar(&value, longVariable, flagDefault, usage)
	flag.StringVar(&value, shortVariable, flagDefault, usage)
	flag.Parse()
	return value
}