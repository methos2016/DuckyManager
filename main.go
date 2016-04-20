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
	var noOnline = flag.Bool("offline", false, "Deactivates online repositories. They'll be listed but not updated, nor interacted with")
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

	// TODO counter to show progress to user

	if !*noOnline {
		// Check all repositories at the same time
		c := make(chan Scripts)
		for _, rep := range config.Repositories {
			go func(ch chan Scripts, rep Repository) {
				s, err2 := rep.GetUpdates()

				if err2 != nil {
					fmt.Println(errStr + translate.ErrUpdatingOnline + " [" + rep.Repo + "]: " + err2.Error())
					l.Println(errStr + translate.ErrUpdatingOnline + " [" + rep.Repo + "]: " + err2.Error())
					os.Exit(errExitCode)
				}

				ch <- s

			}(c, rep)
		}

		// Get all new scripts and add them
		for i := 1; i == len(config.Repositories); i++ {
			newScripts := <-c
			scripts = append(scripts, newScripts...)
		}

		if err != nil {
			fmt.Println(errStr + translate.ErrCheckingLocal + " : " + err.Error())
			l.Println(errStr + translate.ErrCheckingLocal + " : " + err.Error())
			os.Exit(errExitCode)
		}
	}

	scripts = TrimRepeated(scripts)

	// TODO list online ones here too
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

	// make sure we save any changes to the DB and the config file

	if err = Save(config.LocalDBFile, currentState.Scripts); err != nil {
		l.Println(translate.ErrSavingDB + ": " + err.Error())
		fmt.Println(translate.ErrSavingDB + ": " + err.Error())
	}

	if err = Save(configFile, config); err != nil {
		l.Println(translate.ErrSavingConfig + ": " + err.Error())
		fmt.Println(translate.ErrSavingConfig + ": " + err.Error())
	}

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
