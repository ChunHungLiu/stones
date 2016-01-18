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

	log := core.NewLogWidget(0, 11, 80, 10)
	hero := habilis.Skin{"you", core.Glyph{'@', core.ColorWhite}, &tiles[10][5], log, false}
	tiles[10][5].Occupant = &hero

	goblin := habilis.Skin{"goblin", core.Glyph{'g', core.ColorYellow}, &tiles[5][5], nil, false}
	tiles[5][5].Occupant = &goblin

	for !hero.Expired {
		core.TermClear()
		for off, tile := range core.FoV(hero.Pos, radius) {
			req := core.RenderRequest{}
			tile.Handle(&req)
			core.TermDraw(radius+off.X, radius+off.Y, req.Render)
		}
		log.Update()
		core.TermRefresh()

		hero.Handle(&core.Action{})
	}
}
