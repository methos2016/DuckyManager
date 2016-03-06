package main

import (
	"strings"

	"github.com/nsf/termbox-go"
)

func search(scripts []Script) (res []Script) {

	var eB editBox

	var functions = []func([]Script, string) []Script{ListByName, ListByUser, ListByTags, ListByDesc}
	var titles = []string{translate.SidebarTitle, translate.SidebarBy, translate.SidebarTags, translate.SidebarDesc}
	var values = make([]string, len(titles))
	var currentValue = 0
	var done, tab = false, false

	for !done {
		eB.text = []byte(values[currentValue])
		printEditBox(eB, 30, titles[currentValue])

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

	res = TrimRepeated(SortScripts(res))
	return
}

func edit(position int, scripts []Script) {
	script := scripts[position]

	// Iterate through the fields
	var eB editBox
	var values = []*string{&script.Name, &script.User, &script.Tags, &script.Desc}
	var titles = []string{translate.SidebarTitle, translate.SidebarBy, translate.SidebarTags, translate.SidebarDesc}
	var currentValue = 0
	var done, tab = false, false

	for !done {
		eB.text = []byte(*values[currentValue])
		printEditBox(eB, 30, titles[currentValue])

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
	scripts[position] = script

	scripts = SortScripts(scripts)
}

func waitForEnter() {
	ev := termbox.PollEvent()
	for ev.Key != termbox.KeyEnter {
		ev = termbox.PollEvent()
	}
}
