// Package core contains (somewhat) generic roguelike functionality.
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

// Key constants which normally require escapes.
const (
	KeyEsc   Key = 27
	KeyEnter Key = '\r'
	KeyCtrlC Key = Key(termbox.KeyCtrlC)
)

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

// Offset stores a 2-dimensional int vector.
type Offset struct {
	X, Y int
}

// Diff returns the result of subtracting another Offset from this one.
func (o Offset) Diff(a Offset) Offset {
	return Offset{o.X - a.X, o.Y - a.Y}
}

// Max returns the maximum of x and y.
func Max(x, y int) int {
	if y > x {
		return y
	}
	return x
}

// Min returns the minimum of x and y.
func Min(x, y int) int {
	if y < x {
		return y
	}
	return x
}

// Mod returns x modulo y (not the same as x % y, which is remainder).
func Mod(x, y int) int {
	z := x % y
	if z < 0 {
		z += y
	}
	return z
}

// Abs returns the absolute value of x.
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Clamp limits a value to a specific range.
func Clamp(min, val, max int) int {
	if val < min {
		return min
	} else if val > max {
		return max
	}
	return val
}
