package core

import (
	"strconv"
	"testing"
)

func FieldCase(g StrGrid) (origin *Tile, weights map[*Tile]int) {
	cols, rows := len(g[0]), len(g)
	tiles := make([][]Tile, cols)
	for x := 0; x < cols; x++ {
		tiles[x] = make([]Tile, rows)
		for y := 0; y < rows; y++ {
			tiles[x][y].Face = Glyph{'.', ColorWhite}
			tiles[x][y].Pass = true
			tiles[x][y].Adjacent = make(map[Offset]*Tile)
			tiles[x][y].Offset = Offset{x, y}
		}
	}

	origin, weights = nil, make(map[*Tile]int)

	link := func(x, y, dx, dy int) {
		nx, ny := x+dx, y+dy
		if 0 <= nx && nx < cols && 0 <= ny && ny < rows {
			tiles[x][y].Adjacent[Offset{dx, dy}] = &tiles[nx][ny]
		}
	}

	for x := 0; x < cols; x++ {
		for y := 0; y < rows; y++ {
			link(x, y, 1, 1)
			link(x, y, 1, 0)
			link(x, y, 1, -1)
			link(x, y, 0, 1)
			link(x, y, 0, -1)
			link(x, y, -1, 1)
			link(x, y, -1, 0)
			link(x, y, -1, -1)

			c := g[y][x] // care - it really is y,x here
			tiles[x][y].Face = Glyph{rune(c), ColorWhite}
			switch c {
			case '#':
				tiles[x][y].Pass = false
			case '@':
				origin = &tiles[x][y]
				weights[&tiles[x][y]] = 0
			default:
				if weight, err := strconv.Atoi(string(c)); err == nil {
					weights[&tiles[x][y]] = weight
				}
			}
		}
	}

	return origin, weights
}

func TestAttractiveField(t *testing.T) {
	cases := []struct {
		g StrGrid
		r int
	}{
		{
			StrGrid{
				"#######",
				"#@1234#",
				"#11234#",
				"#22234#",
				"#33334#",
				"#44444#",
				"#######",
			}, 10,
		}, {
			StrGrid{
				"########",
				"#987666#",
				"#####56#",
				"#544456#",
				"#543####",
				"#54321@#",
				"########",
			}, 10,
		}, {
			StrGrid{
				"########",
				"#...333#",
				"#####23#",
				"#21@123#",
				"#211####",
				"#22223.#",
				"########",
			}, 3,
		},
	}
	for i, c := range cases {
		origin, weights := FieldCase(c.g)
		actual := AttractiveField(c.r, origin)
		for tile, weight := range weights {
			off := actual.Follow(tile)
			adj := tile.Adjacent[off]
			if adjweight := weights[adj]; adjweight > weight || (adjweight == weight && weight != 0) {
				t.Errorf("AttractiveField failed case %d", i)
			}
		}
	}
}
