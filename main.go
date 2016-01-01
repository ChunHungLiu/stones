package main

import (
	"github.com/jlund3/stones/core"
	"github.com/jlund3/stones/habilis"
)

func main() {
	core.MustTermInit()
	defer core.TermDone()

	cols, rows := 20, 10
	tiles := core.GenStub(cols, rows)
	radius := 5

	hero := habilis.Skin{core.Glyph{'@', core.ColorWhite}, &tiles[10][5]}
	key := core.Key(0)

	for key != core.KeyEsc {
		core.TermClear()
		for off, tile := range core.FoV(hero.Pos, radius) {
			core.TermDraw(off.X+radius, off.Y+radius, tile.Face)
		}
		core.TermDraw(radius, radius, hero.Face)
		core.TermRefresh()

		key = core.GetKey()
		if dx, dy, ok := key.Offset(); ok {
			if adj := hero.Pos.Adjacent[core.Offset{dx, dy}]; adj.Pass {
				hero.Pos = adj
			}
		}
	}
}
