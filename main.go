package main

import (
	"github.com/rauko1753/stones/core"
	"github.com/rauko1753/stones/habilis"
)

var boolgen = core.MapGenBool(func(o core.Offset, pass bool) *core.Tile {
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

var biomes = core.BiomeList{
	core.Biome{
		Boundary:   .3,
		PassChance: 0,
		ImpassTiles: []core.Glyph{
			{Ch: '~', Fg: core.ColorBlue},
			{Ch: '~', Fg: core.ColorLightBlue},
		},
		ImpassLite: true,
	},
	core.Biome{
		Boundary:   .4,
		PassChance: 0,
		ImpassTiles: []core.Glyph{
			{Ch: '~', Fg: core.ColorCyan},
			{Ch: '~', Fg: core.ColorLightCyan},
		},
		ImpassLite: true,
	},
	core.Biome{
		Boundary:   1,
		PassChance: .95,
		PassTiles: []core.Glyph{
			{Ch: '.', Fg: core.ColorGreen},
		},
		ImpassTiles: []core.Glyph{
			{Ch: '%', Fg: core.ColorGreen},
		},
	},
}

func genMaze() []*core.Tile {
	numNodes := 10
	runProb := .5
	weaveProb := 0.
	loopProb := .5
	return boolgen.HalfBraidMaze(numNodes, runProb, weaveProb, loopProb)
}

func genOverworld() *core.Tile {
	h := core.NewHeightmap(200, 400)
	h.Generate()
	overworld := biomes.NewMapGen().Overworld(h)

	for _, tile := range overworld {
		if len(tile.Adjacent) < 8 {
			tile.Pass = false
			tile.Lite = false
			tile.Face = core.Glyph{'#', core.ColorWhite}
		}
	}

	return core.RandPassTile(overworld)
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
