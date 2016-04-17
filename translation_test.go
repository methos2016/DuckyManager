package main

import (
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"testing"
)

// Unused code is tested at compilation time using gometalinter,
// no need to test it here

func clearTranslate(t *testing.T) {
	s := reflect.ValueOf(&translate).Elem()

	for i := 0; i < s.NumField(); i++ {
		s.Field(i).SetString("")
	}
}

func TestTranslationStructUsedOnLang(t *testing.T) {

	// Load languages
	files, err := ioutil.ReadDir(languageDir + "/")
	if err != nil {
		t.Fatal(" Couldn't open '" + languageDir + "' : " + err.Error())
	}

	for _, f := range files {

		// Clean translation from last iter
		clearTranslate(t)

		// Fill it
		checkLangs()

		if err := parseLang(f.Name()); err != nil {
			t.Fatal(err.Error())
			os.Exit(errExitCode)
		}

		// Only for current languages (the mantained ones)
		if translate.Version != languageVer {
			continue
		}

		// Now test if everything is filled
		v := reflect.ValueOf(translate)
		val := reflect.Indirect(reflect.ValueOf(&translate))

		for i := 0; i < v.NumField(); i++ {
			if v.Field(i).String() == "" {
				t.Error("Unused value for language " + f.Name() + ": (n " + strconv.Itoa(i+1) + ") '" + val.Type().Field(i).Name + "'")
			}
		}
	}
}
