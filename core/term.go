package core

import (
	"github.com/nsf/termbox-go"
)

// Color represents the color of a Glyph
type Color uint16

// Color constants for use with ColorChar.
const (
	ColorRed     = Color(termbox.ColorRed)
	ColorBlue    = Color(termbox.ColorBlue)
	ColorCyan    = Color(termbox.ColorCyan)
	ColorBlack   = Color(termbox.ColorBlack)
	ColorGreen   = Color(termbox.ColorGreen)
	ColorWhite   = Color(termbox.ColorWhite)
	ColorYellow  = Color(termbox.ColorYellow)
	ColorMagenta = Color(termbox.ColorMagenta)

	ColorLightRed     = Color(termbox.ColorRed | termbox.AttrBold)
	ColorLightBlue    = Color(termbox.ColorBlue | termbox.AttrBold)
	ColorLightCyan    = Color(termbox.ColorCyan | termbox.AttrBold)
	ColorLightBlack   = Color(termbox.ColorBlack | termbox.AttrBold)
	ColorLightGreen   = Color(termbox.ColorGreen | termbox.AttrBold)
	ColorLightWhite   = Color(termbox.ColorWhite | termbox.AttrBold)
	ColorLightYellow  = Color(termbox.ColorYellow | termbox.AttrBold)
	ColorLightMagenta = Color(termbox.ColorMagenta | termbox.AttrBold)
)

// Glyph pairs a rune with a color.
type Glyph struct {
	Ch rune
	Fg Color
}

// Key represents a single keypress.
type Key rune

// Offset translates a keypress into a direction for both vi-keys and numpad.
// The ok bool indicates whether the keypress corresponded to a direction.
func (k Key) Offset() (dx, dy int, ok bool) {
	switch k {
	case 'h', '4':
		return -1, 0, true
	case 'l', '6':
		return 1, 0, true
	case 'k', '8':
		return 0, -1, true
	case 'j', '2':
		return 0, 1, true
	case 'u', '9':
		return 1, -1, true
	case 'y', '7':
		return -1, -1, true
	case 'n', '3':
		return 1, 1, true
	case 'b', '1':
		return -1, 1, true
	case '.', '5':
		return 0, 0, true
	default:
		return 0, 0, false
	}
}

// Key constants which normally require escapes.
const (
	KeyEsc   Key = 27
	KeyEnter Key = '\r'
	KeyCtrlC Key = Key(termbox.KeyCtrlC)
)

// TermInit readies the terminal for use by the term functions in the core
// package. TermInit should be called before any other term functions are used.
// After a succesful call to TermInit, a call to TermDone should be deferred.
func TermInit() error {
	return termbox.Init()
}

// MustTermInit is like TermInit, except that any errors result in a panic.
func MustTermInit() {
	if err := TermInit(); err != nil {
		panic(err)
	}
}

// TermDone cleans up any setup from TermInit, and reverts the terminal to its
// original state. TermDone should be called after TermInit when the term
// functions in the core package are no longer needed.
func TermDone() {
	termbox.Close()
}

// TermDraw places a Glyph into the internal buffer at the given location.
// No changes are made on screen until TermRefresh is called.
func TermDraw(x, y int, g Glyph) {
	termbox.SetCell(x, y, g.Ch, termbox.Attribute(g.Fg), termbox.ColorBlack)
}

// TermClear erases everything in the internal buffer.
// No changes are made on screen until TermRefresh is called.
func TermClear() {
	termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
}

// TermRefresh ensures that the screen reflects the internal buffer state.
func TermRefresh() {
	termbox.Flush()
}

// State stores the nessesary information to restore a terminal buffer to a
// particular state.
type State [][]termbox.Cell

// TermSave captures the current state of the internal buffer so it can be
// restored later on.
func TermSave() State {
	cols, rows := termbox.Size()
	cells := termbox.CellBuffer()

	state := make(State, rows)
	for y := 0; y < rows; y++ {
		state[y] = make([]termbox.Cell, cols)
		for x := 0; x < cols; x++ {
			state[y][x] = cells[y*cols+x]
		}
	}

	return state
}

// Restore reverts the state of the buffer to the previously saved state.
func (s State) Restore() {
	for y, row := range s {
		for x, cell := range row {
			termbox.SetCell(x, y, cell.Ch, cell.Fg, cell.Bg)
		}
	}
}

// GetKey returns the next keypress. It blocks until there is one.
func GetKey() Key {
	for {
		event := termbox.PollEvent()
		if event.Type == termbox.EventKey {
			return Key(event.Ch) | Key(event.Key)
		}
	}
}
