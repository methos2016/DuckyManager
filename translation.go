package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
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
	MainLoop         string
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
func parseLang() error {
	files, err := ioutil.ReadDir(languageDir + "/")
	if err != nil {
		return errors.New(" Couldn't open '" + languageDir + "' : " + err.Error())
	}

	// Needs a language file
	if len(os.Args) != 2 || os.Args[1] == "" {
		fmt.Println("Usage: DuckyManager <lang>")
		fmt.Println("Your avaliable languages:")

		for _, f := range files {
			tmpLang, err := ioutil.ReadFile(languageDir + "/" + f.Name())
			if err != nil {
				fmt.Println(errStr + f.Name() + " [Could not read]")

			} else if err = json.Unmarshal(tmpLang, &translate); err == nil {
				if translate.Version == languageVer {
					fmt.Println(okStr + f.Name() + " [OK]")
				} else {
					fmt.Println(errStr + f.Name() + " [Outdated]")
				}
			} else {
				fmt.Println(errStr + f.Name() + " [Corrupted]")
			}
		}

		// So it returns an error and doesn't keep going on init
		return errors.New("")
	}

	lang, err := ioutil.ReadFile(languageDir + "/" + os.Args[1])
	if err != nil {
		return errors.New(errStr + "Error opening language file: " + err.Error())
	}

	err = json.Unmarshal(lang, &translate)
	if err != nil {
		return errors.New(errStr + "Error parsing language file: " + err.Error())
	}

	return nil
}
