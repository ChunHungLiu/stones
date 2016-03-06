package core

import (
	"fmt"
)

// Widget serves as a base to various Visual which need relative drawing.
type Widget struct {
	x, y, w, h int
}

// NewWidget creates a Widget with the given location and size.
func NewWidget(x, y, w, h int) Widget {
	return Widget{x, y, w, h}
}

// DrawRel performs a TermDraw relative to the location of the Widget.
// Nothing outside the bounds of the Widget will be drawn.
func (w *Widget) DrawRel(x, y int, g Glyph) {
	if InBounds(x, y, w.w, w.h) {
		TermDraw(x+w.x, y+w.y, g)
	}
}

// TextWidget displays dynamically bound text on the screen.
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

// logmsg is a cached message in LogWidget.
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

// CameraWidget is a Widget which displays an Entity field of view
type CameraWidget struct {
	Widget
	Camera Entity
}

// NewCameraWidget creates a new CameraWidget with the given camera Entity.
func NewCameraWidget(camera Entity, x, y, w, h int) *CameraWidget {
	return &CameraWidget{Widget{x, y, w, h}, camera}
}

// Update draws the camera field of view on screen.
func (w *CameraWidget) Update() {
	req := FoVRequest{}
	w.Camera.Handle(&req)
	cx, cy := w.center()

	for offset, tile := range req.FoV {
		req := RenderRequest{}
		tile.Handle(&req)
		w.DrawRel(cx+offset.X, cy+offset.Y, req.Render)
	}
}

// Mark draws a Glyph on screen relative to the Camera center.
func (w *CameraWidget) Mark(offset Offset, mark Glyph) {
	cx, cy := w.center()
	w.DrawRel(cx+offset.X, cy+offset.Y, mark)
}

// center computes the offset of the camera center relative to the Widget.
func (w *CameraWidget) center() (x, y int) {
	return w.w / 2, w.h / 2
}

// FoVRequest is an Event querying an Entity for a field of view.
type FoVRequest struct {
	FoV map[Offset]*Tile
}

// PercentBarWidget displays a percent bar based on a bound percent function.
type PercentBarWidget struct {
	Widget
	Binding     func() float64
	Vertical    bool
	Invert      bool
	RoundDigits int
	Fill, Empty Glyph
}

// NewPercentBarWidget creates a new PercentBarWidget with the given binding.
func NewPercentBarWidget(binding func() float64, x, y, w, h int) *PercentBarWidget {
	return &PercentBarWidget{Widget{x, y, w, h}, binding, false, false, 2, Glyph{'*', ColorWhite}, Glyph{'-', ColorWhite}}
}

// fillsize computes the size of filled part of the bar on the binding func.
func (b *PercentBarWidget) fillsize() int {
	var max int
	if b.Vertical {
		max = b.h
	} else {
		max = b.w
	}
	return Clamp(0, int(float64(max)*Round(b.Binding(), b.RoundDigits)), max)
}

// isfill returns true if the given x, y is a fill char under the fillsize.
func (b *PercentBarWidget) isfill(x, y, fillsize int) bool {
	if b.Vertical && b.Invert {
		return y < fillsize
	} else if b.Vertical && !b.Invert {
		return b.h-y < fillsize
	} else if !b.Vertical && b.Invert {
		return b.w-x < fillsize
	}
	// else !b.Vertical && !b.Invert
	return x < fillsize
}

// Update displays the PercentBar on screen.
func (b *PercentBarWidget) Update() {
	fillsize := b.fillsize()
	for x := 0; x < b.w; x++ {
		for y := 0; y < b.h; y++ {
			var ch Glyph
			if b.isfill(x, y, fillsize) {
				ch = b.Fill
			} else {
				ch = b.Empty
			}
			b.DrawRel(x, y, ch)
		}
	}
}

// TODO Add non-centering version of CameraWidget
