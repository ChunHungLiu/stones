package core

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
