package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

// TODO Spanish translation

// Strings holds each and every string on the program, allowing for easy translation
type Strings struct {
	Version string

	CheckingLocal string
	NoMatch       string
	AcceptEnter   string

	SidebarTitle string
	SidebarBy    string
	SidebarTags  string
	SidebarDesc  string

	ErrOpeningConfig string
	ErrParsingConfig string
	ErrCheckingLocal string
	ErrPickingFunc   string
	ErrSavingDB      string
	ErrSavingConfig  string
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

	if translate.Version != languageVer {
		return errors.New(errStr + "Language file is outdated")
	}

	return nil
}

func checkLangs() (msg string) {

	files, err := ioutil.ReadDir(languageDir + "/")
	if err != nil {
		msg = " Couldn't open '" + languageDir + "' : " + err.Error()
		return
	}

	msg += "Your avaliable languages:\n\n"

	for _, f := range files {
		tmpLang, err2 := ioutil.ReadFile(languageDir + "/" + f.Name())
		if err2 != nil {
			msg += errStr + f.Name() + " [Could not read]\n"

		} else if err2 = json.Unmarshal(tmpLang, &translate); err2 == nil {
			if translate.Version == languageVer {
				msg += okStr + f.Name() + " [OK]\n"
			} else {
				msg += errStr + f.Name() + " [Outdated]\n"
			}
		} else {
			msg += errStr + f.Name() + " [Corrupted]\n"
		}
	}

	return
}
