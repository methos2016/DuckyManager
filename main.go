package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/nsf/termbox-go"
)

// TODO Fill this with logs, use the correct err/ok str and set size limit for logs

func main() {
	// Load lang
	if err := loadLang(); err != nil {
		fmt.Println(err)
		os.Exit(errExitCode)
	}

	// Load config
	config, err := loadConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(errExitCode)
	}

	// Init log
	f, err := os.OpenFile(config.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(errStr + translate.ErrParsingConfig + ": " + err.Error())
		os.Exit(errExitCode)
	}

	l = log.New(f, "", log.Ltime)

	l.SetOutput(f)

	// Load scripts
	l.Println("+------------------------------+")
	l.Println(translate.CheckingLocal)

	scripts, valid, deleted, modified, newOnes, err := CheckLocal(config.LocalDBFile, config.ScriptsPath)

	if err != nil {
		fmt.Println(errStr + translate.ErrCheckingLocal + " : " + err.Error())
		l.Println(errStr + translate.ErrCheckingLocal + " : " + err.Error())
		os.Exit(errExitCode)
	}

	// make sure we save any changes to the DB
	defer func() {
		if err = Save(config.LocalDBFile, scripts); err != nil {
			l.Println(translate.ErrSavingDB + ": " + err.Error())
			fmt.Println(translate.ErrSavingDB + ": " + err.Error())
		}
	}()

	l.Println("[" + strconv.Itoa(int(valid)) + "] " + translate.Valid + " , " +
		"[" + strconv.Itoa(int(deleted)) + "] " + translate.Deleted + " , " +
		"[" + strconv.Itoa(int(modified)) + "] " + translate.Modified + " , " +
		"[" + strconv.Itoa(int(newOnes)) + "] " + translate.NewScripts)

	// GUI
	err = termbox.Init()
	if err != nil {
		fmt.Println(errStr + translate.ErrTermboxInit + ": " + err.Error())
		l.Println(errStr + translate.ErrTermboxInit + ": " + err.Error())
		os.Exit(errExitCode)
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc)
	termbox.SetOutputMode(termbox.Output256)
	l.Println(okStr + translate.TermInputMode + ": InputESC || " + translate.TermOutputMode + ": Output256")

	currentState := State{
		Scripts:       scripts,
		Position:      0,
		PositionUpper: 0,
	}

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
