package main

import (
	"github.com/jlund3/stones/core"
)

func main() {
	core.MustTermInit()
	defer core.TermDone()

	cols, rows := 20, 10
	tiles := core.GenStub(cols, rows)
	radius := 5

	x, y := cols/2, rows/2
	key := core.Key(0)

	for key != core.KeyEsc {
		core.TermClear()
		for off, tile := range core.FoV(&tiles[x][y], radius) {
			core.TermDraw(off.X+radius, off.Y+radius, tile.Face)
		}
		core.TermDraw(radius, radius, core.Glyph{'@', core.ColorWhite})
		core.TermRefresh()

		key = core.GetKey()
		if dx, dy, ok := key.Offset(); ok {
			if tiles[x+dx][y+dy].Pass {
				x += dx
				y += dy
			}
		}
	}
}
