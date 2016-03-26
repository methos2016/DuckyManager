package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/nsf/termbox-go"
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

func loadGUI(scripts []Script) (currentState State, err error) {
	if err = termbox.Init(); err != nil {
		return State{}, err
	}

	termbox.SetInputMode(termbox.InputEsc)
	termbox.SetOutputMode(termbox.Output256)
	l.Println(okStr + translate.TermInputMode + ": InputESC || " + translate.TermOutputMode + ": Output256")

	currentState = DefaultState(scripts)

	return
}
