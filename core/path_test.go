package core

import (
	"math"
	"testing"
)

type StrGrid []string

func SearchCase(g StrGrid) (origin, goal *Tile, path map[*Tile]struct{}) {
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

	origin, goal = nil, nil
	path = make(map[*Tile]struct{})
	haspath := false

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
			case '$':
				goal = &tiles[x][y]
				path[&tiles[x][y]] = struct{}{}
			case '@':
				origin = &tiles[x][y]
			case 'x':
				path[&tiles[x][y]] = struct{}{}
				haspath = true
			}
		}
	}

	if haspath {
		return
	}
	return origin, goal, nil
}

func PathValid(path []*Tile) bool {
	if len(path) == 0 {
		return true
	}
	for i := 0; i < len(path)-1; i++ {
		curr, next := path[i], path[i+1]
		step := next.Offset.Sub(curr.Offset)
		if actual, ok := curr.Adjacent[step]; !ok || next != actual {
			return false
		}
	}
	return true
}

func PathsEqual(actual []*Tile, expected map[*Tile]struct{}) bool {
	if len(actual) != len(expected) {
		return false
	}
	for _, t := range actual {
		if _, ok := expected[t]; !ok {
			return false
		}
	}
	return true
}

type SearchAlgo func(*Tile, *Tile) []*Tile

func RunCase(t *testing.T, name string, testnum int, algo SearchAlgo, g StrGrid) {
	origin, goal, expected := SearchCase(g)
	actual := algo(origin, goal)
	if !PathValid(actual) || !PathsEqual(actual, expected) {
		t.Errorf("%s failed case %d", name, testnum)
	}
}

func TestAStarPath(t *testing.T) {
	cases := []StrGrid{
		{
			"#######",
			"#$....#",
			"#.x...#",
			"#..x..#",
			"#...x.#",
			"#....@#",
			"#######",
		}, {
			"#######",
			"#$....#",
			"#######",
			"#.....#",
			"#.....#",
			"#....@#",
			"#######",
		}, {
			"########",
			"#$xxx..#",
			"#####x.#",
			"#...x..#",
			"#..x####",
			"#...xx@#",
			"########",
		}, {
			"###########",
			"#.xxx$....#",
			"#x#######.#",
			"#x###...#.#",
			"#x###.#.#.#",
			"#x###.#.#.#",
			"#.xxx@#...#",
			"###########",
		},
	}
	for i, c := range cases {
		RunCase(t, "AStarPath", i, AStarPath, c)
	}
}

func TestGreedyPath(t *testing.T) {
	cases := []StrGrid{
		{
			"#######",
			"#$....#",
			"#.x...#",
			"#..x..#",
			"#...x.#",
			"#....@#",
			"#######",
		}, {
			"#######",
			"#$....#",
			"#######",
			"#.....#",
			"#.....#",
			"#....@#",
			"#######",
		}, {
			"########",
			"#$xxx..#",
			"#####x.#",
			"#...x..#",
			"#..x####",
			"#...xx@#",
			"########",
		}, {
			"###########",
			"#....$xxx.#",
			"#.#######x#",
			"#.###.x.#x#",
			"#.###x#x#x#",
			"#.###x#x#x#",
			"#.###@#x#x#",
			"#.###.#x#x#",
			"#.....#.x.#",
			"###########",
		},
	}
	for i, c := range cases {
		RunCase(t, "GreedyPath", i, GreedyPath, c)
	}
}

func TestCustomSearch(t *testing.T) {
	cost := func(a, b *Tile) float64 {
		delta := b.Offset.Sub(a.Offset)
		if delta.X != 0 && delta.Y != 0 {
			return math.Inf(1)
		}
		return delta.Euclidean()
	}
	search := CreateGraphSearch(cost, euclidean)
	cases := []StrGrid{
		{
			"#######",
			"#$....#",
			"#######",
			"#.....#",
			"#.....#",
			"#....@#",
			"#######",
		}, {
			"########",
			"#$xxxx.#",
			"#####x.#",
			"#..xxx.#",
			"#..x####",
			"#..xxx@#",
			"########",
		}, {
			"###########",
			"#xxxx$....#",
			"#x#######.#",
			"#x###...#.#",
			"#x###.#.#.#",
			"#x###.#.#.#",
			"#xxxx@#...#",
			"###########",
		},
	}
	for i, c := range cases {
		RunCase(t, "GraphSearch", i, search, c)
	}
}
