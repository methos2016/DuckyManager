package main

import (
	"strings"

	"github.com/nsf/termbox-go"
)

func pickFunctionality(ev termbox.Event, currentState State) error {
	if ev.Ch != 0 {
		switch ev.Ch {
		case 's', 'S':

			res, err := search(currentState.Scripts)
			if err != nil {
				return err
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
					return err
				}
			}

		case 'e', 'E':
			if err := edit(currentState); err != nil {
				return err
			}

		}
	}

	return nil
}

func search(scripts []Script) (res []Script, err error) {

	var functions = []func([]Script, string) []Script{ListByName, ListByUser, ListByTags, ListByDesc}
	var titles = []string{translate.SidebarTitle, translate.SidebarBy, translate.SidebarTags, translate.SidebarDesc}
	var values = make([]*string, len(titles))

	if err = editableMenu(titles, values); err != nil {
		return
	}

	res = searchFromValues(scripts, values, functions)
	SortScripts(res)
	res = TrimRepeated(res)
	return
}

func edit(currentState State) (err error) {
	script := currentState.Scripts[currentState.Position]

	var values = []*string{&script.Name, &script.User, &script.Tags, &script.Desc}
	var titles = []string{translate.SidebarTitle, translate.SidebarBy, translate.SidebarTags, translate.SidebarDesc}

	if err = editableMenu(titles, values); err != nil {
		return
	}

	// Update values of script
	currentState.Scripts[currentState.Position] = script

	SortScripts(currentState.Scripts)

	return
}

func searchFromValues(scripts []Script, values []*string, functions []func([]Script, string) []Script) (res []Script) {
	for i, value := range values {
		if *value != "" {
			// Search multiple, for tags
			if strings.Contains(",", *value) {
				for _, v := range strings.Split(*value, ",") {
					res = append(res, functions[i](scripts, v)...)
				}
			} else {
				res = append(res, functions[i](scripts, *value)...)
			}

		}
	}

	return
}

func waitForEnter() {
	ev := termbox.PollEvent()
	for ev.Key != termbox.KeyEnter {
		ev = termbox.PollEvent()
	}
}
