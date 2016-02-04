package core

import (
	"unicode"
)

// Label is a Visual which displays fixed text on screen.
type Label struct {
	Text string
	X, Y int
	Fg   Color
}

// NewLabel creates a new label with the given text.
func NewLabel(text string, x, y int) *Label {
	return &Label{text, x, y, ColorWhite}
}

// Update draws the Label text at the given location.
func (l *Label) Update() {
	for i, ch := range l.Text {
		TermDraw(l.X+i, l.Y, Glyph{ch, l.Fg})
	}
}

// Border is a Visual which displays a border
type Border struct {
	Widget
	UpperLeft, UpperRight, LowerLeft, LowerRight Glyph
	Vertical, Horizontal                         Glyph
}

// NewBorder creates a new Border with the given parameters.
func NewBorder(vert, horiz Glyph, x, y, w, h int) *Border {
	return &Border{Widget{x, y, w, h}, horiz, horiz, horiz, horiz, vert, horiz}
}

// Update draws the Border on screen.
func (w *Border) Update() {
	w.DrawRel(0, 0, w.UpperLeft)
	w.DrawRel(w.w-1, 0, w.UpperRight)
	w.DrawRel(0, w.h-1, w.LowerLeft)
	w.DrawRel(w.w-1, w.h-1, w.LowerRight)
	for y := 1; y < w.h-1; y++ {
		w.DrawRel(0, y, w.Vertical)
		w.DrawRel(w.w-1, y, w.Vertical)
	}
	for x := 1; x < w.w-1; x++ {
		w.DrawRel(x, 0, w.Horizontal)
		w.DrawRel(x, w.h-1, w.Horizontal)
	}
}

// TextBox is an Element which allows a user to enter custom text.
type TextBox struct {
	Text string
	Len  int
	X, Y int

	colorSelect
	ExtraCh rune
}

// NewTextBox returns a new TextBox with the given text.
func NewTextBox(text string, length, x, y int) *TextBox {
	return &TextBox{text, length, x, y, colorSelect{ColorWhite, ColorLightWhite}, '_'}
}

// Update draws the current text.
func (t *TextBox) Update(selected bool) {
	color := t.getColor(selected)

	for x := 0; x < t.Len; x++ {
		if x < len(t.Text) {
			TermDraw(t.X+x, t.Y, Glyph{rune(t.Text[x]), color})
		} else {
			TermDraw(t.X+x, t.Y, Glyph{t.ExtraCh, color})
		}
	}
}

// Activate lets the user enter text into the TextBox.
func (t *TextBox) Activate() FormResult {
	old := t.Text
	t.Text = ""
	t.Update(true)
	TermRefresh()

	var key Key
	for key != KeyEnter && key != KeyEsc {
		key = GetKey()
		if unicode.IsPrint(rune(key)) {
			t.Text += string(key)
		}
		t.Update(true)
		TermRefresh()
	}

	if key == KeyEsc {
		t.Text = old
	}
	return nil
}

// Submit is a button which submits a particular FormResult.
type Submit struct {
	Text   string
	X, Y   int
	Result FormResult

	colorSelect
}

// NewSubmit creates a new Submit with the given result.
func NewSubmit(text string, x, y int, result FormResult) *Submit {
	return &Submit{text, x, y, result, colorSelect{ColorWhite, ColorLightWhite}}
}

// Update displays the Submit on screen.
func (s *Submit) Update(selected bool) {
	color := s.getColor(selected)
	for i, ch := range s.Text {
		TermDraw(s.X+i, s.Y, Glyph{ch, color})
	}
}

// Activate returns the Submit FormResult.
func (s *Submit) Activate() FormResult {
	return s.Result
}

// colorSelect is used to let an Element have customizable Color selection.
type colorSelect struct {
	NormalFg, SelectedFg Color
}

func (s colorSelect) getColor(selected bool) Color {
	if selected {
		return s.SelectedFg
	} else {
		return s.NormalFg
	}
}
