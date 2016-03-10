package main

import (
	"strings"

	"github.com/nsf/termbox-go"
)

func pickFunctionality(ev termbox.Event, currentState State) {
	if ev.Ch != 0 {
		switch ev.Ch {
		case 's', 'S':

			res, err := search(currentState.Scripts)
			if err != nil {
				//TODO Handle err
			}

			if len(res) != 0 {
				mainLoop(State{
					Scripts:       res,
					Position:      0,
					PositionUpper: 0,
				})

			} else {
				err := showErrorMsg(translate.NoMatch)
				if err != nil {
					//TODO Handle err
				}
			}

		case 'e', 'E':
			edit(currentState)
		}
	}
}

func search(scripts []Script) (res []Script, err error) {

	var eB editBox

	var functions = []func([]Script, string) []Script{ListByName, ListByUser, ListByTags, ListByDesc}
	var titles = []string{translate.SidebarTitle, translate.SidebarBy, translate.SidebarTags, translate.SidebarDesc}
	var values = make([]string, len(titles))
	var currentValue = 0
	var done, tab = false, false

	for !done {
		eB.text = []byte(values[currentValue])
		if err = printEditBox(eB, 30, titles[currentValue]); err != nil {
			return
		}

		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			// Iterate
			case termbox.KeyTab:
				values[currentValue] = string(eB.text)
				if currentValue+1 == len(values) {
					currentValue = 0
				} else {
					currentValue++
				}
				tab = true

			// Save
			case termbox.KeyEnter:
				done = true
			// Close without saving
			case termbox.KeyCtrlC, termbox.KeyEsc:
				return

			// Editing stuff
			case termbox.KeyArrowLeft:
				eB.MoveCursorOneRuneBackward()
			case termbox.KeyArrowRight:
				eB.MoveCursorOneRuneForward()
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				eB.DeleteRuneBackward()
			case termbox.KeyDelete:
				eB.DeleteRuneForward()
			case termbox.KeySpace:
				eB.InsertRune(' ')
			case termbox.KeyHome:
				eB.MoveCursorToBeginningOfTheLine()
			case termbox.KeyEnd:
				eB.MoveCursorToEndOfTheLine()
			default:
				if ev.Ch != 0 {
					eB.InsertRune(ev.Ch)
				}
			}

		case termbox.EventError:
			l.Println(translate.ErrEvent + ": " + ev.Err.Error())
		}

		if !tab {
			values[currentValue] = string(eB.text)
		} else {
			tab = false
		}
	}

	// Search!
	for i, value := range values {
		if value != "" {
			// Search multiple, for tags
			if strings.Contains(",", value) {
				for _, v := range strings.Split(value, ",") {
					res = append(res, functions[i](scripts, v)...)
				}
			} else {
				res = append(res, functions[i](scripts, value)...)
			}

		}
	}
	SortScripts(res)
	res = TrimRepeated(res)
	return
}

func edit(currentState State) {
	script := currentState.Scripts[currentState.Position]

	// Iterate through the fields
	var eB editBox
	var values = []*string{&script.Name, &script.User, &script.Tags, &script.Desc}
	var titles = []string{translate.SidebarTitle, translate.SidebarBy, translate.SidebarTags, translate.SidebarDesc}
	var currentValue = 0
	var done, tab = false, false

	for !done {
		eB.text = []byte(*values[currentValue])
		if err := printEditBox(eB, 30, titles[currentValue]); err != nil {
			//TODO handle
		}

		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			// Iterate
			case termbox.KeyTab:
				*values[currentValue] = string(eB.text)
				if currentValue+1 == len(values) {
					currentValue = 0
				} else {
					currentValue++
				}
				tab = true

			// Save
			case termbox.KeyCtrlS, termbox.KeyEnter:
				done = true
			// Close without saving
			case termbox.KeyCtrlC, termbox.KeyEsc:
				return

			// Editing stuff
			case termbox.KeyArrowLeft:
				eB.MoveCursorOneRuneBackward()
			case termbox.KeyArrowRight:
				eB.MoveCursorOneRuneForward()
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				eB.DeleteRuneBackward()
			case termbox.KeyDelete:
				eB.DeleteRuneForward()
			case termbox.KeySpace:
				eB.InsertRune(' ')
			case termbox.KeyHome:
				eB.MoveCursorToBeginningOfTheLine()
			case termbox.KeyEnd:
				eB.MoveCursorToEndOfTheLine()
			default:
				if ev.Ch != 0 {
					eB.InsertRune(ev.Ch)
				}
			}

		case termbox.EventError:
			l.Println(translate.ErrEvent + ": " + ev.Err.Error())
		}

		if !tab {
			*values[currentValue] = string(eB.text)
		} else {
			tab = false
		}
	}

	// Update values of script
	currentState.Scripts[currentState.Position] = script

	SortScripts(currentState.Scripts)
}

func waitForEnter() {
	ev := termbox.PollEvent()
	for ev.Key != termbox.KeyEnter {
		ev = termbox.PollEvent()
	}
}
