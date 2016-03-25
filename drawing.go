package main

import (
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

func guiPrint(x, y, w int,
	fg, bg termbox.Attribute,
	msg string,
) {
	for _, c := range msg {
		if x == w {
			termbox.SetCell(x-1, y, '→', fg, bg)
			break
		}
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}

// TODO show title on main screen as "main, search of "x"...
func redrawMain(currentState State) error {
	if err := termbox.Clear(coldef, coldef); err != nil {
		return err
	}

	w, h := termbox.Size()

	sidebarDraw(w, h, currentState)
	listScripts(w, h, currentState)
	return termbox.Sync()
}

func listScripts(totalW, totalH int, currentState State) {

	w := totalW * 2 / 3
	h := totalH

	x := 0
	y := totalH - h

	for i, c := currentState.PositionUpper, 0; c < h && c < len(currentState.Scripts); c++ {

		name := currentState.Scripts[i].GetName()

		if i == currentState.Position {
			guiPrint(x+1, y+c, w, termbox.AttrBold, coldef, name)
		} else {
			guiPrint(x+1, y+c, w, coldef, coldef, name)
		}
		i++
	}

	SortScripts(currentState.Scripts)
}

func sidebarDraw(totalW, totalH int, currentState State) {

	w := totalW / 3
	h := totalH

	x := totalW - w
	y := totalH - h

	// Draw left side line
	for i := y; i < h; i++ {
		termbox.SetCell(x, i, '│', coldef, coldef)
	}

	// Title
	lines := printSideInfo(x+2, y+1, w, h, translate.SidebarTitle, currentState.GetCurrentScript().Name)

	// By
	lines = printSideInfo(x+2, y+2+lines, w, h, translate.SidebarBy, currentState.GetCurrentScript().User)

	// Targets
	lines = printSideInfo(x+2, y+2+lines, w, h, translate.SidebarTags, currentState.GetCurrentScript().Tags)

	// Desc
	printSideInfo(x+2, y+2+lines, w, h, translate.SidebarDesc, currentState.GetCurrentScript().Desc)
}

func printSideInfo(x, y, w, h int,
	title, msg string,
) (line int) {

	dotLen := len(": ")
	titleLen := len(title)

	line = y

	guiPrint(x, line, w, termbox.AttrUnderline, coldef, title)
	guiPrint(x+titleLen, line, w, coldef, coldef, ": ")

	letterCount := 0

	for _, c := range msg {
		currentLen := 0

		// Starting line
		if line == y {
			currentLen = x + titleLen + dotLen + letterCount
		} else {
			currentLen = x + letterCount
		}

		guiPrint(currentLen, line, w, termbox.AttrBold, coldef, string(c))

		l.Println(currentLen-x, w, string(c))

		if currentLen-x+4 > w {
			line++
			letterCount = 0
		}

		letterCount++
	}

	/* Used when the tags where an slice.. may come up again later
	// Add comma if there is more
	if msgLen > 1 && i+1 < msgLen {
		if line == y {
			guiPrint(x+titleLen+dotLen+letterCount, line, w, coldef, coldef, ", ")
		} else {
			guiPrint(x+letterCount, line, w, coldef, coldef, ", ")
		}
		letterCount += 2
	}*/

	return
}

func printEditBox(eB editBox, editBoxWidth int, title string) (err error) {

	w, h := termbox.Size()

	midy := h / 2
	midx := (w - editBoxWidth) / 2

	// unicode box drawing chars around the edit box
	termbox.SetCell(midx-1, midy, '│', coldef, coldef)
	termbox.SetCell(midx+editBoxWidth, midy, '│', coldef, coldef)
	termbox.SetCell(midx-1, midy-1, '┌', coldef, coldef)
	termbox.SetCell(midx-1, midy+1, '└', coldef, coldef)
	termbox.SetCell(midx+editBoxWidth, midy-1, '┐', coldef, coldef)
	termbox.SetCell(midx+editBoxWidth, midy+1, '┘', coldef, coldef)
	fill(midx, midy-1, editBoxWidth, 1, termbox.Cell{Ch: '─'})
	fill(midx, midy+1, editBoxWidth, 1, termbox.Cell{Ch: '─'})

	// Title
	guiPrint(midx, midy-1, editBoxWidth, termbox.AttrBold, coldef, title)

	eB.Draw(midx, midy, editBoxWidth, 1)
	termbox.SetCursor(midx+eB.CursorX(), midy)

	err = termbox.Flush()
	return
}

func showErrorMsg(msg string) (err error) {
	w, h := termbox.Size()

	midy := h / 2
	midx := (w - len(msg)) / 2

	if err = termbox.Clear(coldef, coldef); err != nil {
		return
	}

	guiPrint(midx, midy, w, termbox.ColorRed, coldef, msg)
	guiPrint(midx, midy+2, w, termbox.AttrUnderline, coldef, translate.AcceptEnter)
	if err = termbox.Sync(); err != nil {
		return
	}

	waitForEnter()

	return
}

func editableMenu(titles []string, values []*string) (err error) {
	var eB editBox
	var currentValue = 0
	var action = -1

	for action != actionEnter {
		eB.text = []byte(*values[currentValue])
		if err = printEditBox(eB, 30, titles[currentValue]); err != nil {
			return
		}

		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			action = editableMenuSwitchKey(&eB, values, currentValue, ev)
			if action == actionEsc {
				return
			}

		case termbox.EventError:
			l.Println(translate.ErrEvent + ": " + ev.Err.Error())
		}

		*values[currentValue] = string(eB.text)

		if action == actionTab {
			if currentValue+1 == len(values) {
				currentValue = 0
			} else {
				currentValue++
			}
		}
	}

	return
}

// TODO 3 booleans seems dirty. Maybe fix this sometime.

const (
	actionTab   = 0
	actionEnter = 1
	actionEsc   = 2
)

// editableMenuSwitchKey switches between the posibilities on the editable menu.
// Returns 0 if the user tabbed,
// 		   1 if the user pressed Enter,
// 		   2 if the user pressed Ctrl+C or Esc
func editableMenuSwitchKey(eB *editBox, values []*string, currentValue int, ev termbox.Event) (action int) {
	action = -1
	switch ev.Key {
	// Iterate
	case termbox.KeyTab:
		action = actionTab

	// Save
	case termbox.KeyEnter:
		action = actionEnter
	// Close without saving
	case termbox.KeyCtrlC, termbox.KeyEsc:
		action = actionEsc

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
	return
}

/* To be used...

//TODO freaking use this as standard on all menus..
func drawBox(x, y, w, h int,
	title string,
	titleEffect, titleBGEffect termbox.Attribute,
) (err error) {

	for i := 0; i < h; i++ {
		termbox.SetCell(x, y+i, '│', coldef, coldef)
		termbox.SetCell(w, y+i, '│', coldef, coldef)
	}

	fill(x, y, w, 1, termbox.Cell{Ch: '─'})
	fill(x, h, w, 1, termbox.Cell{Ch: '─'})

	termbox.SetCell(x, y, '┌', coldef, coldef)
	termbox.SetCell(x, h, '└', coldef, coldef)

	termbox.SetCell(w, y, '┐', coldef, coldef)
	termbox.SetCell(w, h, '┘', coldef, coldef)

	guiPrint(x+1, y, w, titleEffect, titleBGEffect, title)

	err = termbox.Flush()
	return
}


func printOptionsBox(maxOptionsLine, selected int,
	options []string,
	title string,
) (err error) {

	w, h := termbox.Size()

	midy := h / 2
	midx := w / 2

	// Calculate lines
	var nLines int
	if len(options)%maxOptionsLine != 0 {
		nLines = (len(options) / maxOptionsLine) + 1
	} else {
		nLines = len(options) / maxOptionsLine
	}
	nL2 := nLines / 2

	// Calculate max width
	maxW := 0
	for i := 0; i < len(options); i = i + maxOptionsLine {
		tmpW := 0

		for y := 0; y < maxOptionsLine && y+i < len(options); y++ {
			tmpW += len(options[i+y]) + len(" ")
		}
		if tmpW > maxW {
			maxW = tmpW
		}
	}

	// Print options
	le := 0
	y := 0
	for i, option := range options {
		if i == selected {
			guiPrint(midx+le-(maxW/2)+1, midy-nL2+y, w, termbox.AttrBold, termbox.AttrReverse, option+" ")
		} else {
			guiPrint(midx+le-(maxW/2)+1, midy-nL2+y, w, coldef, coldef, option+" ")
		}

		le += len(option + " ")

		if i%maxOptionsLine-1 == 0 && i != 0 {
			y++
			le = 0
		}
	}

	err = drawBox(midx-maxW/2, midy-nL2, maxW, midy+nL2, title, termbox.AttrBold, coldef)
	return
} */
