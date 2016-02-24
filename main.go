package main

import (
	"github.com/rauko1753/stones/core"
	"github.com/rauko1753/stones/habilis"
)

func main() {
	origin := core.HalfBraidMaze(50, .5, 1, .5)

	core.MustTermInit()
	defer core.TermDone()

	hero := habilis.Skin{
		Name: "you",
		Face: core.Glyph{Ch: '@', Fg: core.ColorWhite},
		Pos:  origin,
	}
	origin.Occupant = &hero

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
