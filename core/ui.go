package core

import (
	"fmt"
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

func Target(camera Entity, m Marker) (target *Tile, ok bool) {
	state := TermSave()
	defer state.Restore()

	req := FoVRequest{}
	camera.Handle(&req)
	offset := Offset{}

	var key Key
	for key != KeyEnter {
		state.Restore()
		for _, o := range Trace(offset) {
			m.Mark(o, Glyph{'*', ColorBlue})
		}
		m.Mark(offset, Glyph{'*', ColorRed})
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
