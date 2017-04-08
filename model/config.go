package model

import (
"encoding/json"
"os"
	"log"
	"errors"
)

const defaultConfigFile = "data/config.json"

// RetrieveConfig reads and unmarshals the main config data file.
func RetrieveConfig(fileName string) ([]Configuration, error) {
	// Open the file.
	configFile := fileName
	if configFile == "" {
		configFile = defaultConfigFile
	}
	log.Println(configFile)
	file, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}

	// Schedule the file to be closed once
	// the function returns.
	defer file.Close()

	// Decode the file into a slice of pointers
	// to Feed values.
	var config []Configuration
	err = json.NewDecoder(file).Decode(&config)
	
	if err == nil {
		if len(config) == 0 {
			return nil, errors.New("Please provide at least one configuration ...")
		}
	}

	// We don't need to check for errors, the caller can do this.
	return config, err
}

