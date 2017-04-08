package model

import (
"encoding/json"
"os"
)

const dataFile = "data/data.json"


// RetrieveSites reads and unmarshals the sites data file.
func RetrieveSites(fileName string) ([]Site, error) {
	configFile := fileName
	if configFile == "" {
		configFile = dataFile
	}
	
	// Open the file.
	file, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}

	// Schedule the file to be closed once
	// the function returns.
	defer file.Close()

	// Decode the file into a slice of pointers
	// to Feed values.
	var servers []Site
	err = json.NewDecoder(file).Decode(&servers)

	// We don't need to check for errors, the caller can do this.
	return servers, err
}

