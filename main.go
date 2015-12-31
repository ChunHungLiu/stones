package main

import (
	"github.com/jlund3/stones/core"
)

func main() {
	core.MustTermInit()
	defer core.TermDone()

	for x, ch := range "Hello World!" {
		core.TermDraw(x, 0, core.Glyph{ch, core.ColorWhite})
	}
	core.TermRefresh()

	core.GetKey()
}
