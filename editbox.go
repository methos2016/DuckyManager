package main

// Adaptation of the "editBox" example from termbox

import (
	"unicode/utf8"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

func fill(x, y, w, h int, cell termbox.Cell) {
	for ly := 0; ly < h; ly++ {
		for lx := 0; lx < w; lx++ {
			termbox.SetCell(x+lx, y+ly, cell.Ch, cell.Fg, cell.Bg)
		}
	}
}

func runeAdvanceLen(r rune, pos int) int {
	if r == '\t' {
		return tabstopLength - pos%tabstopLength
	}
	return runewidth.RuneWidth(r)
}

func vOffsetCOffset(text []byte, boffset int) (voffset, coffset int) {
	text = text[:boffset]
	for len(text) > 0 {
		r, size := utf8.DecodeRune(text)
		text = text[size:]
		coffset++
		voffset += runeAdvanceLen(r, voffset)
	}
	return
}

func byteSliceGrow(s []byte, desiredCap int) []byte {
	if cap(s) < desiredCap {
		ns := make([]byte, len(s), desiredCap)
		copy(ns, s)
		return ns
	}
	return s
}

func byteSliceRemove(text []byte, from, to int) []byte {
	size := to - from
	copy(text[from:], text[to:])
	text = text[:len(text)-size]
	return text
}

func byteSliceInsert(text []byte, offset int, what []byte) []byte {
	n := len(text) + len(what)
	text = byteSliceGrow(text, n)
	text = text[:n]
	copy(text[offset+len(what):], text[offset:])
	copy(text[offset:], what)
	return text
}

const preferedHorizontalThreshold = 5
const tabstopLength = 8

type editBox struct {
	text          []byte
	lineVOffset   int
	cursorBOffset int // cursor offset in bytes
	cursorVOffset int // visual cursor offset in termbox cells
	cursorCOffset int // cursor offset in unicode code points
}

func (eb *editBox) Draw(x, y, w, h int) {
	eb.AdjustVOffset(w)

	const coldef = termbox.ColorDefault
	fill(x, y, w, h, termbox.Cell{Ch: ' '})

	t := eb.text
	lx := 0
	tabstop := 0
	for {
		rx := lx - eb.lineVOffset
		if len(t) == 0 {
			break
		}

		if lx == tabstop {
			tabstop += tabstopLength
		}

		if rx >= w {
			termbox.SetCell(x+w-1, y, '→',
				coldef, coldef)
			break
		}

		r, size := utf8.DecodeRune(t)
		if r == '\t' {
			for ; lx < tabstop; lx++ {
				rx = lx - eb.lineVOffset
				if rx >= w {
					goto next
				}

				if rx >= 0 {
					termbox.SetCell(x+rx, y, ' ', coldef, coldef)
				}
			}
		} else {
			if rx >= 0 {
				termbox.SetCell(x+rx, y, r, coldef, coldef)
			}
			lx += runewidth.RuneWidth(r)
		}
	next:
		t = t[size:]
	}

	if eb.lineVOffset != 0 {
		termbox.SetCell(x, y, '←', coldef, coldef)
	}
}

// Adjusts line visual offset to a proper value depending on width
func (eb *editBox) AdjustVOffset(width int) {
	ht := preferedHorizontalThreshold
	maxHThreshold := (width - 1) / 2
	if ht > maxHThreshold {
		ht = maxHThreshold
	}

	threshold := width - 1
	if eb.lineVOffset != 0 {
		threshold = width - ht
	}
	if eb.cursorVOffset-eb.lineVOffset >= threshold {
		eb.lineVOffset = eb.cursorVOffset + (ht - width + 1)
	}

	if eb.lineVOffset != 0 && eb.cursorVOffset-eb.lineVOffset < ht {
		eb.lineVOffset = eb.cursorVOffset - ht
		if eb.lineVOffset < 0 {
			eb.lineVOffset = 0
		}
	}
}

func (eb *editBox) MoveCursorTo(boffset int) {
	eb.cursorBOffset = boffset
	eb.cursorVOffset, eb.cursorCOffset = vOffsetCOffset(eb.text, boffset)
}

func (eb *editBox) RuneUnderCursor() (rune, int) {
	return utf8.DecodeRune(eb.text[eb.cursorBOffset:])
}

func (eb *editBox) RuneBeforeCursor() (rune, int) {
	return utf8.DecodeLastRune(eb.text[:eb.cursorBOffset])
}

func (eb *editBox) MoveCursorOneRuneBackward() {
	if eb.cursorBOffset == 0 {
		return
	}
	if eb.cursorCOffset > len(eb.text) {
		eb.MoveCursorToBeginningOfTheLine()
		return
	}

	_, size := eb.RuneBeforeCursor()
	eb.MoveCursorTo(eb.cursorBOffset - size)
}

func (eb *editBox) MoveCursorOneRuneForward() {
	if eb.cursorBOffset == len(eb.text) {
		return
	}
	if eb.cursorCOffset > len(eb.text) {
		eb.MoveCursorTo(len(eb.text))
	}
	_, size := eb.RuneUnderCursor()
	eb.MoveCursorTo(eb.cursorBOffset + size)
}

func (eb *editBox) MoveCursorToBeginningOfTheLine() {
	eb.MoveCursorTo(0)
}

func (eb *editBox) MoveCursorToEndOfTheLine() {
	eb.MoveCursorTo(len(eb.text))
}

func (eb *editBox) DeleteRuneBackward() {
	if eb.cursorBOffset == 0 {
		return
	}

	eb.MoveCursorOneRuneBackward()
	_, size := eb.RuneUnderCursor()
	eb.text = byteSliceRemove(eb.text, eb.cursorBOffset, eb.cursorBOffset+size)
}

func (eb *editBox) DeleteRuneForward() {
	if eb.cursorBOffset == len(eb.text) {
		return
	}
	_, size := eb.RuneUnderCursor()
	eb.text = byteSliceRemove(eb.text, eb.cursorBOffset, eb.cursorBOffset+size)
}

func (eb *editBox) DeleteTheRestOfTheLine() {
	eb.text = eb.text[:eb.cursorBOffset]
}

func (eb *editBox) InsertRune(r rune) {
	var buf [utf8.UTFMax]byte

	if eb.cursorCOffset > len(eb.text) {
		eb.MoveCursorToBeginningOfTheLine()
	}

	n := utf8.EncodeRune(buf[:], r)
	eb.text = byteSliceInsert(eb.text, eb.cursorBOffset, buf[:n])
	eb.MoveCursorOneRuneForward()
}

func (eb *editBox) CursorX() int {
	return eb.cursorVOffset - eb.lineVOffset
}
