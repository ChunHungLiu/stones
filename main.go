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

	hero := habilis.Skin{core.Glyph{'@', core.ColorWhite}, &tiles[10][5], false}
	tiles[10][5].Occupant = &hero

	for !hero.Expired {
		core.TermClear()
		for off, tile := range core.FoV(hero.Pos, radius) {
			req := core.RenderRequest{}
			tile.Handle(&req)
			core.TermDraw(radius+off.X, radius+off.Y, req.Render)
		}
		core.TermRefresh()

		hero.Handle(&core.Action{})
	}
}
