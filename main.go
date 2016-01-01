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
	tiles[10][5].Occupant = &hero
	key := core.Key(0)

	for key != core.KeyEsc {
		core.TermClear()
		for off, tile := range core.FoV(hero.Pos, radius) {
			req := core.RenderRequest{}
			tile.Handle(&req)
			core.TermDraw(radius+off.X, radius+off.Y, req.Render)
		}
		core.TermRefresh()

		key = core.GetKey()
		if dx, dy, ok := key.Offset(); ok {
			if adj := hero.Pos.Adjacent[core.Offset{dx, dy}]; adj.Pass {
				hero.Pos.Occupant, adj.Occupant = nil, &hero
				hero.Pos = adj
			}
		}
	}
}
