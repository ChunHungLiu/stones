package core

import (
	"strconv"
	"testing"
	"unicode"
)

func AttractiveFieldCase(g StrGrid) (goals []*Tile, weights map[*Tile]int) {
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

	goals, weights = make([]*Tile, 0), make(map[*Tile]int)

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
				goals = append(goals, &tiles[x][y])
				weights[&tiles[x][y]] = 0
			default:
				if weight, err := strconv.Atoi(string(c)); err == nil {
					weights[&tiles[x][y]] = weight
				}
			}
		}
	}

	return goals, weights
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
				"#000333#",
				"#####23#",
				"#21@123#",
				"#211####",
				"#222230#",
				"########",
			}, 3,
		}, {
			StrGrid{
				"#############",
				"#@1234555555#",
				"#11234444444#",
				"#22234333333#",
				"#33334322222#",
				"#44444321112#",
				"#55554321@12#",
				"#66654321112#",
				"#############",
			}, 10,
		}, {
			StrGrid{
				"########",
				"#@12300#",
				"#####00#",
				"#000000#",
				"#003####",
				"#00321@#",
				"########",
			}, 3,
		},
	}
	for i, c := range cases {
		goals, weights := AttractiveFieldCase(c.g)
		actual := AttractiveField(c.r, goals...)
		for tile, weight := range weights {
			off := actual.Follow(tile)
			adj := tile.Adjacent[off]
			if adjweight := weights[adj]; adjweight > weight || (adjweight == weight && weight != 0) {
				t.Errorf("AttractiveField failed case %d", i)
			}
		}
	}
}

func ReplusiveFieldCase(g StrGrid) (ungoals []*Tile, weights map[*Tile]int) {
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

	ungoals, weights = make([]*Tile, 0), make(map[*Tile]int)

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

			c := rune(g[y][x]) // care - it really is y,x here
			tiles[x][y].Face = Glyph{rune(c), ColorWhite}
			switch c {
			case '#':
				tiles[x][y].Pass = false
			default:
				if unicode.IsLetter(c) {
					weights[&tiles[x][y]] = int(unicode.ToLower(c) - 'a')
					if !unicode.IsLower(c) {
						ungoals = append(ungoals, &tiles[x][y])
					}
				} else if unicode.IsDigit(c) {
					weights[&tiles[x][y]] = int(c - '0')
				}
			}
		}
	}

	return ungoals, weights
}

func TestRepusliveField(t *testing.T) {
	cases := []struct {
		g StrGrid
		r int
	}{
		{
			StrGrid{
				"#######",
				"#Edcba#",
				"#ddcba#",
				"#cccba#",
				"#bbbba#",
				"#aaaaa#",
				"#######",
			}, 10,
		}, {
			StrGrid{
				"########",
				"#mlkjii#",
				"######H#",
				"######g#",
				"######f#",
				"#abcdef#",
				"########",
			}, 10,
		}, {
			StrGrid{
				"########",
				"#000aaa#",
				"#####ba#",
				"#bcDcba#",
				"#bcc####",
				"#bbbba0#",
				"########",
			}, 3,
		}, {
			StrGrid{
				"#############",
				"#Ggggggggghi#",
				"#ffffffffghi#",
				"#eeeeeeefghi#",
				"#ddddddefghi#",
				"#cccccdefghi#",
				"#bbbbcdefGhi#",
				"#aaabcdefghi#",
				"#############",
			}, 10,
		},
	}
	for i, c := range cases {
		goals, weights := ReplusiveFieldCase(c.g)
		actual := ReplusiveField(c.r, goals...)
		for tile, weight := range weights {
			off := actual.Follow(tile)
			adj := tile.Adjacent[off]
			if adjweight := weights[adj]; adjweight > weight || (adjweight == weight && weight != 0) {
				t.Errorf("ReplusiveField failed case %d", i)
			}
		}
	}
}

// TODO TONS of copy-paste here
// StrGrid Tile creation copied from Path tests
// TestCast funcs duplicated in this file
// Test runs are close to copied
