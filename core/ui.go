package core

import (
	"fmt"
	"strings"

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

type Marker interface {
	Mark(o Offset, g Glyph)
}

type Targetter struct {
	Camera Entity
	Marker
	Reticle Glyph
	Trace   *Glyph
	Accept  string
}

// Aim allows the user to select a target from an on-screen Camera view.
func (t Targetter) Aim() (target *Tile, ok bool) {
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
				t.Mark(o, *t.Trace)
			}
		}
		t.Mark(offset, Glyph{'*', ColorRed})
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
func Aim(camera Entity, marker Marker, accept string) (target *Tile, ok bool) {
	return Targetter{camera, marker, Glyph{'*', ColorRed}, nil, accept}.Aim()
}
