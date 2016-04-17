package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/nsf/termbox-go"
)

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

func loadGUI(scripts []Script) (currentState State, err error) {
	if err = termbox.Init(); err != nil {
		return State{}, err
	}

	termbox.SetInputMode(termbox.InputEsc)
	termbox.SetOutputMode(termbox.Output256)

	currentState = DefaultState(scripts)

	return
}
