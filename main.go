package main

import (
	"github.com/rauko1753/stones/core"
	"github.com/rauko1753/stones/habilis"
)

func main() {
	core.MustTermInit()
	defer core.TermDone()

	cols, rows := 20, 10
	tiles := core.GenStub(cols, rows)

	hero := habilis.Skin{
		Name: "you",
		Face: core.Glyph{'@', core.ColorWhite},
		Pos:  tiles[10][5],
	}
	tiles[10][5].Occupant = &hero

	goblin := habilis.Skin{
		Name: "goblin",
		Face: core.Glyph{'g', core.ColorYellow},
		Pos:  tiles[5][5],
	}
	tiles[5][5].Occupant = &goblin

	log := core.NewLogWidget(0, 11, 80, 10)
	view := core.NewCameraWidget(&hero, 0, 0, 11, 11)
	screen := core.Screen{log, view}
	hero.View = view

	hero.Logger = log

	for !hero.Expired {
		screen.Update()
		hero.Handle(&habilis.Action{})
	}
}
