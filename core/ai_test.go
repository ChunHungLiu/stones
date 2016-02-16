package core

type StrGrid []string

func (g StrGrid) Convert(callback func(*Tile, byte)) [][]Tile {
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

			tile := &tiles[x][y]
			ch := g[y][x] // care - it really is y, x (not x, y)
			callback(tile, ch)
		}
	}

	return tiles
}
