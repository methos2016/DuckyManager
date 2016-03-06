package main

import (
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

const coldef = termbox.ColorDefault

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

func redrawMain(positionUpper, position int, scripts []Script) error {
	termbox.Clear(coldef, coldef)

	w, h := termbox.Size()
	sidebarDraw(w, h, position, scripts)
	listScripts(w, h, positionUpper, position, scripts)

	return termbox.Sync()
}

func listScripts(totalW, totalH, positionUpper, position int,
	scripts []Script,
) {

	w := totalW * 2 / 3
	h := totalH

	x := 0
	y := totalH - h

	for i, c := positionUpper, 0; c < h && c < len(scripts); c++ {

		name := scripts[i].GetName()

		if i == position {
			guiPrint(x+1, y+c, w, termbox.AttrBold, coldef, name)
		} else {
			guiPrint(x+1, y+c, w, coldef, coldef, name)
		}
		i++
	}

	scripts = SortScripts(scripts)
}

func sidebarDraw(totalW, totalH, position int, scripts []Script) {

	w := totalW / 3
	h := totalH

	x := totalW - w
	y := totalH - h

	// Draw left side line
	for i := y; i < h; i++ {
		termbox.SetCell(x, i, '│', coldef, coldef)
	}

	// Title
	lines := printSideInfo(x+2, y+1, w, h, translate.SidebarTitle, scripts[position].Name)

	// By
	lines = printSideInfo(x+2, y+2+lines, w, h, translate.SidebarBy, scripts[position].User)

	// Targets
	lines = printSideInfo(x+2, y+2+lines, w, h, translate.SidebarTags, scripts[position].Tags)

	// Desc
	printSideInfo(x+2, y+2+lines, w, h, translate.SidebarDesc, scripts[position].Desc)
}

func printSideInfo(x, y, w, h int,
	title string,
	msg string,
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

func drawBox(x, y, w, h int,
	title string,
	titleEffect, titleBGEffect termbox.Attribute,
) {

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

	termbox.Flush()
}

func printEditBox(eB editBox, editBoxWidth int, title string) {

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

	termbox.Flush()
}

func printOptionsBox(maxOptionsLine, selected int,
	options []string,
	title string,
) {

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

	drawBox(midx-maxW/2, midy-nL2, maxW, midy+nL2, title, termbox.AttrBold, coldef)
}

func showErrorMsg(msg string) {
	w, h := termbox.Size()

	midy := h / 2
	midx := (w - len(msg)) / 2

	termbox.Clear(coldef, coldef)

	guiPrint(midx, midy, w, termbox.ColorRed, coldef, msg)
	guiPrint(midx, midy+2, w, termbox.AttrUnderline, coldef, translate.AcceptEnter)
	termbox.Sync()

	waitForEnter()
}
