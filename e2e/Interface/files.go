package Interface

import (
	"encoding/json"
	"os"
	"strings"
)

type resultParameters struct {
	Code   int      `json:"code"`
	Header []params `json:"header"`
	Body   string   `json:"body"`
}

type params struct {
	Parameter string `json:"parameter"`
	Value     string `json:"value"`
}

type ConfigFile struct {
	Name             string             `json:"name"`
	Path             string             `json:"path"`
	Verb             string             `json:"verb"`
	RequestHeader    []params           `json:"requestHeader"`
	RequestPath      []params           `json:"requestPath"`
	RequestBody      string             `json:"requestBody"`
	AvailableResults []resultParameters `json:"results"`
}

func ReadTestFile(filePath string) ([]ConfigFile, error) {
	var configFileData []ConfigFile
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return configFileData, err
	}

	err = json.Unmarshal(fileData, &configFileData)
	if err != nil {
		return configFileData, err
	}

	return configFileData, nil
}

func SearchForTestFiles(dirPath string) ([]string, []string) {
	var acceptedList []string
	var rejectedList []string

	files, err := os.ReadDir(dirPath)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		fileWithPath := dirPath + "/" + file.Name()
		if file.IsDir() {
			accepted, rejected := SearchForTestFiles(fileWithPath)
			acceptedList = addFileToList(acceptedList, accepted)
			rejectedList = addFileToList(rejectedList, rejected)
		}

		if strings.Contains(file.Name(), ".json") {
			acceptedList = append(acceptedList, fileWithPath)
		} else {
			rejectedList = append(rejectedList, fileWithPath)
		}
	}

	return acceptedList, rejectedList
}

func addFileToList(list []string, fileWithPath []string) []string {
	for _, f := range fileWithPath {
		list = append(list, f)
	}

	return list
}
