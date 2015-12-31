package core

// GenStub is a temporary map gen for testing.
func GenStub(cols, rows int) [][]Tile {
	backing := make([]Tile, cols*rows)
	tiles := make([][]Tile, cols)
	for x := 0; x < cols; x++ {
		tiles[x] = backing[x*rows : (x+1)*rows]
		for y := 0; y < rows; y++ {
			if x == 0 || x == cols-1 || y == 0 || y == rows-1 {
				tiles[x][y].Face = Glyph{'#', ColorWhite}
			} else if RandChance(.1) {
				tiles[x][y].Face = Glyph{'%', ColorGreen}
			} else {
				tiles[x][y].Face = Glyph{'.', ColorWhite}
			}
		}
	}
	return tiles
}
