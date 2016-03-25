package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

func loadLang() (err error) {
	if err = checkLangs(os.Args); err != nil {
		return
	}
	err = parseLang(os.Args[1])
	return
}

func loadConfig() (config Config, err error) {
	var cf []byte
	if cf, err = ioutil.ReadFile(configFile); err != nil {
		return Config{}, errors.New(translate.ErrOpeningConfig + ": " + err.Error())
	}

	if err = json.Unmarshal(cf, &config); err != nil {
		return Config{}, errors.New(errStr + translate.ErrParsingConfig + ": " + err.Error())
	}

	return
}
