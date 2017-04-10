package model

import (
	"encoding/json"
	"os"
)

const indexDataFile = "data/indexservice.json"

// RetrieveIndex reads and unmarshals the index data file.
func RetrieveIndex(fileName string) (IndexConfig, error) {
	configFile := fileName
	if configFile == "" {
		configFile = indexDataFile
	}
	println("Index: " + configFile)
	indexServer := IndexConfig{}
	// Open the file.
	file, err := os.Open(configFile)
	if err != nil {
		return indexServer, err
	}

	// Schedule the file to be closed once
	// the function returns.
	defer file.Close()

	// Decode the file into a slice of pointers
	// to Feed values.
	err = json.NewDecoder(file).Decode(&indexServer)

	// We don't need to check for errors, the caller can do this.
	return indexServer, err

}
