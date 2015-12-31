package core

// GenStub is a temporary map gen for testing.
func GenStub(cols, rows int) [][]Tile {
	tiles := make([][]Tile, cols)
	for x := 0; x < cols; x++ {
		tiles[x] = make([]Tile, rows)
		for y := 0; y < rows; y++ {
			tiles[x][y].Face = Glyph{'.', ColorWhite}
			tiles[x][y].Pass = true
			tiles[x][y].Adjacent = make(map[Offset]*Tile)
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
		}
	}

	for x := 0; x < cols; x++ {
		for y := 0; y < rows; y++ {
			if x == 0 || x == cols-1 || y == 0 || y == rows-1 {
				tiles[x][y].Face = Glyph{'#', ColorWhite}
				tiles[x][y].Pass = false
			} else if RandChance(.1) {
				tiles[x][y].Face = Glyph{'%', ColorGreen}
				tiles[x][y].Pass = false
			}
		}
	}

	return tiles
}
