package core

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/nsf/termbox-go"
)

// ListSelect displays a list of items and allows the user to select one item.
func ListSelect(title string, items []interface{}) (index int, ok bool) {
	state := TermSave()
	defer state.Restore()

	rows := []string{title}
	cols := len(title)
	for i, item := range items {
		row := fmt.Sprintf("%c) %v", i+'a', item)
		rows = append(rows, row)
		cols = Max(cols, len(row))
	}

	for y, row := range rows {
		for x, ch := range row {
			TermDraw(x, y, Glyph{ch, ColorWhite})
		}
		for x := len(row); x < cols; x++ {
			TermDraw(x, y, Glyph{' ', ColorWhite})
		}
	}
	TermRefresh()

	index = int(GetKey() - 'a')
	if index < 0 || index >= len(items) {
		return 0, false
	}
	return index, true
}

// TermTint recolors every glyph in the buffer to have the given color.
// No changes are made on screen until RefreshScreen is called.
func TermTint(c Color) {
	fg := termbox.Attribute(c)
	cells := termbox.CellBuffer()
	for i := 0; i < len(cells); i++ {
		cells[i].Fg = fg
	}
}

// Targeter allows for customization of on-screen targeting.
type Targeter struct {
	Camera  Entity
	Canvas  Entity
	Reticle Glyph
	Trace   *Glyph
	Accept  string
}

// Aim allows the user to select a target from an on-screen Camera view.
func (t Targeter) Aim() (target *Tile, ok bool) {
	state := TermSave()
	defer state.Restore()

	req := FoVRequest{}
	t.Camera.Handle(&req)
	offset := Offset{}

	var key Key
	for !strings.Contains(t.Accept, string(key)) && key != KeyEsc {
		state.Restore()

		if t.Trace != nil {
			for _, o := range Trace(offset) {
				t.Canvas.Handle(&Mark{o, *t.Trace})
			}
		}
		t.Canvas.Handle(&Mark{offset, t.Reticle})
		TermRefresh()

		key = GetKey()
		delta, ok := KeyMap[key]
		_, visible := req.FoV[offset.Add(delta)]
		if ok && visible {
			offset = offset.Add(delta)
		}
	}

	return req.FoV[offset], key != KeyEsc
}

// Aim allows the user to select a target from an on-screen Camera view.
func Aim(camera, canvas Entity, accept string) (target *Tile, ok bool) {
	return Targeter{camera, canvas, Glyph{'*', ColorRed}, nil, accept}.Aim()
}

// Mark is an Event requesting that a Glyph be drawn on Screen.
type Mark struct {
	Offset Offset
	Mark   Glyph
}

// Label is a Visual which displays fixed text on screen.
type Label struct {
	Text string
	X, Y int
}

// Update draws the Label text at the given location.
func (l Label) Update() {
	for i, ch := range l.Text {
		TermDraw(l.X+i, l.Y, Glyph{ch, ColorWhite})
	}
}

// Border is a Visual which displays a border
type Border struct {
	Widget
	Vertical, Horizontal Glyph
}

// NewBorder creates a new Border with the given parameters.
func NewBorder(vert, horiz Glyph, x, y, w, h int) *Border {
	return &Border{Widget{x, y, w, h}, vert, horiz}
}

// Update draws the Border on screen.
func (w *Border) Update() {
	for y := 0; y < w.h; y++ {
		w.DrawRel(0, y, w.Vertical)
		w.DrawRel(w.w-1, y, w.Vertical)
	}
	for x := 0; x < w.w; x++ {
		w.DrawRel(x, 0, w.Horizontal)
		w.DrawRel(x, w.h-1, w.Horizontal)
	}
}

// TextBox is an Element which allows a user to enter custom text.
type TextBox struct {
	Text string
	Len  int
	X, Y int
}

// Update draws the current text.
func (t *TextBox) Update(selected bool) {
	var color Color
	if selected {
		color = ColorLightWhite
	} else {
		color = ColorWhite
	}

	for x := 0; x < t.Len; x++ {
		if x < len(t.Text) {
			TermDraw(t.X+x, t.Y, Glyph{rune(t.Text[x]), color})
		} else {
			TermDraw(t.X+x, t.Y, Glyph{'_', color})
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
