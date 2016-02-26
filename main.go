package main

import (
	"github.com/rauko1753/stones/core"
	"github.com/rauko1753/stones/habilis"
)

func genMaze() *core.Tile {
	numNodes := 50
	runProb := .5
	weaveProb := .5
	loopProb := .5
	gen := core.BoolTileFactory(func(o core.Offset, pass bool) *core.Tile {
		t := core.NewTile(o)
		t.Pass = pass
		t.Lite = pass
		if pass {
			t.Face = core.Glyph{'.', core.ColorLightRed}
		} else {
			t.Face = core.Glyph{'#', core.ColorRed}
		}
		return t
	})
	return gen.HalfBraidMaze(numNodes, runProb, weaveProb, loopProb)
}

func genOverworld() *core.Tile {
	gen := core.FloatTileFactory(func(o core.Offset, height float64) *core.Tile {
		t := core.NewTile(o)
		switch {
		case height < .5:
			t.Face = core.Glyph{'~', core.ColorBlue}
			t.Pass = false
		default:
			t.Face = core.Glyph{'.', core.ColorGreen}
		}
		return t
	})
	h := core.NewHeightmap(40, 40)
	h.Generate()
	tiles := gen.Apply(h)
	for x := 0; x < h.Cols(); x++ {
		tiles[x][0].Pass = false
		tiles[x][h.Rows()-1].Pass = false
	}
	for y := 0; y < h.Rows(); y++ {
		tiles[0][y].Pass = false
		tiles[h.Cols()-1][y].Pass = false
	}
	return tiles[20][20]
}

func main() {
	core.MustTermInit()
	defer core.TermDone()

	origin := genOverworld()

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
