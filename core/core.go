// Package core contains (somewhat) generic roguelike functionality.
package core

import (
	"math"

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
	KeyEsc   Key = Key(termbox.KeyEsc)
	KeyEnter Key = Key(termbox.KeyEnter)
	KeyCtrlC Key = Key(termbox.KeyCtrlC)
	KeyPgup  Key = Key(termbox.KeyPgup)
	KeyPgdn  Key = Key(termbox.KeyPgdn)
)

// Offset stores a 2-dimensional int vector.
type Offset struct {
	X, Y int
}

// KeyMap stores default directional Key values. This dictionary can be edited
// to affect any core functions which require knowledge of directional keys.
var KeyMap = map[Key]Offset{
	'h': {-1, 0}, '4': {-1, 0},
	'l': {1, 0}, '6': {1, 0},
	'k': {0, -1}, '8': {0, -1},
	'j': {0, 1}, '2': {0, 1},
	'u': {1, -1}, '9': {1, -1},
	'y': {-1, -1}, '7': {-1, -1},
	'n': {1, 1}, '3': {1, 1},
	'b': {-1, 1}, '1': {-1, 1},
}

// Sub returns the result of subtracting another Offset from this one.
func (o Offset) Sub(a Offset) Offset {
	return Offset{o.X - a.X, o.Y - a.Y}
}

// Add returns the result of adding another Offset to this one.
func (o Offset) Add(a Offset) Offset {
	return Offset{o.X + a.X, o.Y + a.Y}
}

// Manhattan returns the L_1 distance off the Offset.
func (o Offset) Manhattan() int {
	return Abs(o.X) + Abs(o.Y)
}

// Euclidean returns the L_2 distance off the Offset.
func (o Offset) Euclidean() float64 {
	return math.Hypot(float64(o.X), float64(o.Y))
}

// Chebyshev returns the L_inf distance off the Offset.
func (o Offset) Chebyshev() int {
	return Max(Abs(o.X), Abs(o.Y))
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

// InBounds returns true if x in [0, w) and y in [0, h).
func InBounds(x, y, w, h int) bool {
	return 0 <= x && x < w && 0 <= y && y < h
}
