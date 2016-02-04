package core

import (
	"github.com/nsf/termbox-go"
)

// TermInit readies the terminal for use by the term functions in the core
// package. TermInit should be called before any other term functions are used.
// After a successful call to TermInit, a call to TermDone should be deferred.
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

// Visual represents something which can be drawn in the terminal.
type Visual interface {
	Update()
}

// Screen is a collection of Visual.
type Screen []Visual

// Update clears the screen, and draws each Visual in the Screen.
func (s Screen) Update() {
	TermClear()
	for _, v := range s {
		v.Update()
	}
	TermRefresh()
}

// FormResult describes the result from running a Form.
type FormResult interface {
	Result() string
}

// resultstr is the default implementation of FormResult.
type resultstr string

// NewFormResult wraps a string as a FormResult.
func NewFormResult(s string) FormResult {
	return resultstr(s)
}

// ResultEsc is the result from Form.Run when escape is pressed.
var ResultEsc = NewFormResult("ESCAPE")

// Result unwraps the resultstr into a string.
func (r resultstr) Result() string {
	return string(r)
}

// Element represents an activatable element on a Form.
type Element interface {
	Update(selected bool)
	Activate() FormResult
}

// Form is a collection for Visual and Element for building a TUI screen.
type Form struct {
	Visuals  []Visual
	Elements []Element
}

// update runs Update on each Visual and Element. Used in both Update and Run.
func (f Form) update(selected int) {
	TermClear()
	for _, v := range f.Visuals {
		v.Update()
	}
	for i, e := range f.Elements {
		e.Update(i == selected)
	}
	TermRefresh()
}

// Update clears the screen, and draws each Visual and Element in the Form.
// No Element is considered selected during the update.
func (f Form) Update() {
	f.update(-1)
}

// Run allows the user to select and activate Form Elements. Run returns any
// non-nil FormResult from an activated Element. Additionally, ResultEsc is
// returned if the user hits escape.
func (f Form) Run() FormResult {
	curr := 0
	for {
		f.update(curr)

		switch key := GetKey(); key {
		case KeyEnter:
			if result := f.Elements[curr].Activate(); result != nil {
				return result
			}
		case KeyEsc:
			return ResultEsc
		default:
			if delta, ok := KeyMap[key]; ok && delta.X == 0 {
				curr = Mod(curr+delta.Y, len(f.Elements))
			}
		}
	}
}

// TODO Add Term interface with termbox and gopherlibterminal implementations
