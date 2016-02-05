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

// Button is an Element which runs a callback upon activation.
type Button struct {
	Text    string
	X, Y    int
	Binding func() FormResult

	colorSelect
}

// NewButton creats a new Button with the given callback.
func NewButton(text string, x, y int, callback func() FormResult) *Button {
	return &Button{text, x, y, callback, colorSelect{ColorWhite, ColorLightWhite}}
}

// NewSubmit creates a new Button which simply returns a FormResult.
func NewSubmit(text string, x, y int, result FormResult) *Button {
	return NewButton(text, x, y, func() FormResult { return result })
}

// Update displays the Button on screen.
func (b *Button) Update(selected bool) {
	color := b.getColor(selected)
	for i, ch := range b.Text {
		TermDraw(b.X+i, b.Y, Glyph{ch, color})
	}
}

// Activate runs the Button callback and returns the FormResult.
func (b *Button) Activate() FormResult {
	return b.Binding()
}

// colorSelect is used to let an Element have customizable Color selection.
type colorSelect struct {
	NormalFg, SelectedFg Color
}

// getColor returns a color based on whether the Element is selected or not.
func (s colorSelect) getColor(selected bool) Color {
	if selected {
		return s.SelectedFg
	}
	return s.NormalFg
}
