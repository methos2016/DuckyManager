package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

// TODO check for unused strings

// Strings holds each and every string on the program, allowing for easy translation
type Strings struct {
	Version string

	CheckingLocal string
	Info          string
	TooLong       string
	NoMatch       string
	AcceptEnter   string
	LoadingLocal  string

	NewScripts string
	Valid      string
	Deleted    string
	Modified   string
	Any        string

	TermInputMode  string
	TermOutputMode string

	SidebarTitle string
	SidebarBy    string
	SidebarTags  string
	SidebarDesc  string

	ErrClearScreen   string
	ErrOpeningConfig string
	ErrParsingConfig string
	ErrCheckingLocal string
	ErrPickingFunc   string
	ErrSavingDB      string
	ErrTermboxInit   string
	ErrDrawing       string
	ErrEvent         string
}

// Meant to be called on init
func parseLang(langFile string) error {

	lang, err := ioutil.ReadFile(languageDir + "/" + langFile)
	if err != nil {
		return errors.New(errStr + "Error opening language file: " + err.Error())
	}

	err = json.Unmarshal(lang, &translate)
	if err != nil {
		return errors.New(errStr + "Error parsing language file: " + err.Error())
	}

	return nil
}

// Everything OK if BOTH msgs and err are empty ("" && nil)
func checkLangs(args []string) (msgs string, err error) {
	files, err := ioutil.ReadDir(languageDir + "/")
	if err != nil {
		return "", errors.New(" Couldn't open '" + languageDir + "' : " + err.Error())
	}

	// If incorrect args, is not considered an error, but a msg will be returned.
	if len(args) != 2 || args[1] == "" {
		msgs += "Usage: DuckyManager <lang>\n"
		msgs += "Your avaliable languages:\n\n"

		for _, f := range files {
			tmpLang, err := ioutil.ReadFile(languageDir + "/" + f.Name())
			if err != nil {
				msgs += errStr + f.Name() + " [Could not read]\n"

			} else if err = json.Unmarshal(tmpLang, &translate); err == nil {
				if translate.Version == languageVer {
					msgs += okStr + f.Name() + " [OK]\n"
				} else {
					msgs += errStr + f.Name() + " [Outdated]\n"
				}
			} else {
				msgs += errStr + f.Name() + " [Corrupted]\n"
			}
		}
	}

	return
}
