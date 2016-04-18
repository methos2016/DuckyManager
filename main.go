package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/nsf/termbox-go"
)

// TODO Create debug lines (logs)
// TODO Script editor with syntax

func main() {
	// Debug flag
	flag.BoolVar(&debug, "debug", false, "Activates debug output to logs. Meant for bug fixing/reports")
	flag.StringVar(&lang, "lang", "en", "Sets the language for the program")
	var checkLang = flag.Bool("sanity", false, "Sanity checks installed languages")
	flag.Parse()

	if *checkLang {
		fmt.Println(checkLangs())
		return
	}

	// Load lang
	if err := parseLang(lang); err != nil {
		fmt.Println(err)
		os.Exit(errExitCode)
	}

	// Load config
	config, err := LoadConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(errExitCode)
	}

	// Init log if debug is active
	if debug {
		var f *os.File
		f, err = os.OpenFile(config.LogFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			fmt.Println(errStr + translate.ErrParsingConfig + ": " + err.Error())
			os.Exit(errExitCode)
		}

		l = log.New(f, "", log.Ltime)
		l.SetOutput(f)
	} else {
		l = log.New(ioutil.Discard, "", log.Ltime)
	}

	// Load scripts
	l.Println("+------------------------------+")
	l.Println(translate.CheckingLocal)
	scripts, valid, deleted, modified, newOnes, err := CheckLocal(config.LocalDBFile, config.ScriptsPath)

	if err != nil {
		fmt.Println(errStr + translate.ErrCheckingLocal + " : " + err.Error())
		l.Println(errStr + translate.ErrCheckingLocal + " : " + err.Error())
		os.Exit(errExitCode)
	}

	// make sure we save any changes to the DB and the config file
	defer func() {
		if err = Save(config.LocalDBFile, scripts); err != nil {
			l.Println(translate.ErrSavingDB + ": " + err.Error())
			fmt.Println(translate.ErrSavingDB + ": " + err.Error())
		}

		if err = Save(configFile, config); err != nil {
			l.Println(translate.ErrSavingConfig + ": " + err.Error())
			fmt.Println(translate.ErrSavingConfig + ": " + err.Error())
		}

	}()

	l.Println("[" + strconv.Itoa(int(deleted)) + "] Deleted , " +
		"[" + strconv.Itoa(int(modified)) + "] Modified , " +
		"[" + strconv.Itoa(int(newOnes)) + "] New , " +
		"[" + strconv.Itoa(int(valid)) + "] Valid")

	// GUI
	currentState, err := loadGUI(scripts)
	if err != nil {
		fmt.Println(errStr + translate.ErrTermboxInit + ": " + err.Error())
		l.Println(errStr + translate.ErrTermboxInit + ": " + err.Error())
		os.Exit(errExitCode)
	}
	defer termbox.Close()

	mainLoop(currentState)
}

func mainLoop(currentState State) {
	exit := false

	for !exit {
		if err := redrawMain(currentState); err != nil {
			l.Println(errStr + translate.ErrDrawing + ": " + err.Error())
			os.Exit(errExitCode)
		}

		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc, termbox.KeyCtrlC:
				exit = true
			default:
				currentState.SwitchKey(ev)
			}
		case termbox.EventError:
			l.Println(errStr + translate.ErrEvent + ": " + ev.Err.Error())
			os.Exit(errExitCode)
		}

	}
}
