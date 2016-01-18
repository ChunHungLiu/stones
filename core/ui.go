package core

import (
	"fmt"
)

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

func (f Form) Run() FormResult {
	curr := 0
	for {
		TermClear()
		for _, v := range f.Visuals {
			v.Update()
		}
		for i, e := range f.Elements {
			e.Update(i == curr)
		}
		TermRefresh()

		switch key := GetKey(); key {
		case KeyEnter:
			if result := f.Elements[curr].Activate(); result != nil {
				return result
			}
		case KeyEsc:
			return ResultEsc
		default:
			if dx, dy, ok := key.Offset(); ok && dx == 0 {
				curr = Mod(curr+dy, len(f.Elements))
			}
		}
	}
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

// Widget serves as a base to various Visual which need relative drawing.
type Widget struct {
	x, y, w, h int
}

func NewWidget(x, y, w, h int) Widget {
	return Widget{x, y, w, h}
}

func (w *Widget) DrawRel(x, y int, g Glyph) {
	if InBounds(x, y, w.w, w.h) {
		TermDraw(x+w.x, y+w.h, g)
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

func (m *logmsg) String() string {
	if m.Count == 1 {
		return m.Text
	}
	return fmt.Sprintf("%s (x%d)", m.Text, m.Count)
}

type LogWidget struct {
	Widget
	cache []*logmsg
}

func NewLogWidget(x, y, w, h int) *LogWidget {
	return &LogWidget{Widget{x, y, w, h}, make([]*logmsg, 0)}
}

func (w *LogWidget) AddMessage(msg string) {
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
