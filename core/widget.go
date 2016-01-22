package core

import (
	"fmt"
)

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

// Widget serves as a base to various Visual which need relative drawing.
type Widget struct {
	x, y, w, h int
}

// NewWidget creates a Widget with the given location and size.
func NewWidget(x, y, w, h int) Widget {
	return Widget{x, y, w, h}
}

// DrawRel performs a TermDraw relative to the location and size of the Widget.
func (w *Widget) DrawRel(x, y int, g Glyph) {
	if InBounds(x, y, w.w, w.h) {
		TermDraw(x+w.x, y+w.y, g)
	}
}

// TextWidget displays dynamic text on the screen.
type TextWidget struct {
	Widget
	Binding func() string
}

// NewTextWidget creates a new TextWidget with the given binding.
func NewTextWidget(binding func() string, x, y, w, h int) *TextWidget {
	return &TextWidget{Widget{x, y, w, h}, binding}
}

// Update draws the bound text on screen.
func (w *TextWidget) Update() {
	x, y := 0, 0
	for _, ch := range w.Binding() {
		if ch == '\n' {
			x, y = 0, y+1
		} else {
			w.DrawRel(x, y, Glyph{ch, ColorWhite})
			x++
		}
	}
}

type logmsg struct {
	Text  string
	Count int
	Seen  bool
}

// String implements fmt.Stringer for logmsg.
func (m *logmsg) String() string {
	if m.Count == 1 {
		return m.Text
	}
	return fmt.Sprintf("%s (x%d)", m.Text, m.Count)
}

// LogWidget is a Widget which stores and display log messages.
type LogWidget struct {
	Widget
	cache []*logmsg
}

// NewLogWidget creates a new empty LogWidget.
func NewLogWidget(x, y, w, h int) *LogWidget {
	return &LogWidget{Widget{x, y, w, h}, make([]*logmsg, 0)}
}

// Log places a new message in the LogWidget cache.
func (w *LogWidget) Log(msg string) {
	last := len(w.cache) - 1
	// if cache is empty, or last message text was different than this one
	if last < 0 || w.cache[last].Text != msg {
		w.cache = append(w.cache, &logmsg{msg, 1, false})
		// truncate cache if too long to show on the widget
		if len(w.cache) > w.h {
			w.cache = w.cache[len(w.cache)-w.h:]
		}
	} else { // duplicate text, so just reuse last message
		w.cache[last].Count++
		w.cache[last].Seen = false
	}
}

func (w *LogWidget) Clear() {
	w.cache = w.cache[:0]
}

// Update draws the cached log messages on screen.
func (w *LogWidget) Update() {
	for y, msg := range w.cache {
		// determine color based on seen
		var fg Color
		if msg.Seen {
			fg = ColorLightBlack
		} else {
			fg = ColorWhite
		}

		// note we assume no newlines, unlike TextWidget.
		for x, ch := range msg.String() {
			w.DrawRel(x, y, Glyph{ch, fg})
		}

		// we just displayed the message, so next time should be seen
		msg.Seen = true
	}
}

type FoVRequest struct {
	FoV map[Offset]*Tile
}

type CameraWidget struct {
	Widget
	Camera Entity
}

func NewCameraWidget(camera Entity, x, y, w, h int) *CameraWidget {
	return &CameraWidget{Widget{x, y, w, h}, camera}
}

func (w *CameraWidget) Update() {
	req := FoVRequest{}
	w.Camera.Handle(&req)
	centerx, centery := w.w/2, w.h/2

	for offset, tile := range req.FoV {
		req := RenderRequest{}
		tile.Handle(&req)
		w.DrawRel(centerx+offset.X, centery+offset.Y, req.Render)
	}
}
