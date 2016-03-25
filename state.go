package main

import "github.com/nsf/termbox-go"

// State is just a holder for the current state of the program
type State struct {
	Scripts       []Script
	Position      int
	PositionUpper int
}

// Down handles a keystroke on the "down" arrow
func (s *State) Down() {
	if s.Position+1 < len(s.Scripts) {
		s.Position++

		_, h := termbox.Size()

		if s.Position-s.PositionUpper > h-1 {
			s.PositionUpper++
		}
	}
}

// Up handles a keystroke on the "up" arrow
func (s *State) Up() {
	if s.Position-1 >= 0 {
		s.Position--

		if s.Position < s.PositionUpper {
			s.PositionUpper--
		}
	}
}

// Home handles a keystroke on the "home" key
func (s *State) Home() {
	s.Position = 0
	s.PositionUpper = 0
}

// End handles a keystroke on the "end" key
func (s *State) End() {
	s.Position = len(s.Scripts) - 1
	_, h := termbox.Size()
	s.PositionUpper = len(s.Scripts) - h
	if s.PositionUpper < 0 {
		s.PositionUpper = 0
	}
}

// SwitchKey will pick the correct function for the pressed key
func (s *State) SwitchKey(ev termbox.Event) {
	switch ev.Key {
	case termbox.KeyArrowDown:
		s.Down()

	case termbox.KeyArrowUp:
		s.Up()

	case termbox.KeyHome:
		s.Home()

	case termbox.KeyEnd:
		s.End()

	default:
		// will call mainLoop by itself if needed, then come back
		if err := pickFunctionality(ev, *s); err != nil {
			l.Println(translate.ErrPickingFunc + ": " + err.Error())
		}
	}
}

// GetCurrentScript will return the script on the current position
func (s *State) GetCurrentScript() Script {
	return s.Scripts[s.Position]
}
