package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

// Strings holds each and every string on the program, allowing for easy translation
type Strings struct {
	Version          string
	CheckingLocal    string
	ErrClearScreen   string
	ErrOpeningConfig string
	ErrParsingConfig string
	ErrCheckingLocal string
	Info             string
	Valid            string
	Deleted          string
	Modified         string
	TermInputMode    string
	TermOutputMode   string
	NewScripts       string
	SidebarTitle     string
	SidebarBy        string
	SidebarTags      string
	SidebarDesc      string
	ErrTermboxInit   string
	ErrDrawing       string
	ErrEvent         string
	TooLong          string
	Any              string
	NoMatch          string
	AcceptEnter      string
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

func checkLangs(args []string, debug bool) (err error) {
	files, err := ioutil.ReadDir(languageDir + "/")
	if err != nil {
		return errors.New(" Couldn't open '" + languageDir + "' : " + err.Error())
	}

	if len(args) != 2 || args[1] == "" {
		err = errors.New("Incorrect args")

		if !debug {
			fmt.Println("Usage: DuckyManager <lang>")
			fmt.Println("Your avaliable languages:")
		}

		for _, f := range files {
			tmpLang, err := ioutil.ReadFile(languageDir + "/" + f.Name())
			if err != nil && !debug {
				fmt.Println(errStr + f.Name() + " [Could not read]")

			} else if err = json.Unmarshal(tmpLang, &translate); err == nil && !debug {
				if translate.Version == languageVer {
					fmt.Println(okStr + f.Name() + " [OK]")
				} else {
					fmt.Println(errStr + f.Name() + " [Outdated]")
				}
			} else if !debug {
				fmt.Println(errStr + f.Name() + " [Corrupted]")
			}
		}
	}

	return
}
