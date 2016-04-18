package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

// Config holds the configuration for the program
type Config struct {
	// LogFile is the filename for the log file
	LogFile string
	// LocalDBFile is the filename for the local database
	LocalDBFile string
	// ScriptsPath is the local path for your scripts
	ScriptsPath string
	// Repositories is the saved repositories from which the aplication feeds
	Repositories []Repository
}

// LoadConfig will create a config object from the config file
func LoadConfig() (config Config, err error) {
	var cf []byte
	if cf, err = ioutil.ReadFile(configFile); err != nil {
		return Config{}, errors.New(translate.ErrOpeningConfig + ": " + err.Error())
	}

	if err = json.Unmarshal(cf, &config); err != nil {
		return Config{}, errors.New(errStr + translate.ErrParsingConfig + ": " + err.Error())
	}

	return
}
