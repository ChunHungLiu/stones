package core

func createTileGrid(cols, rows int, o Offset, f func(o Offset) *Tile) []*Tile {
	backing := make([]*Tile, cols*rows)

	tiles := make([][]*Tile, cols)
	for x := 0; x < cols; x++ {
		tiles[x] = backing[x*rows : (x+1)*rows]
		for y := 0; y < rows; y++ {
			tiles[x][y] = f(o.Add(Offset{x, y}))
		}
	}

	link := func(x, y, dx, dy int) {
		nx, ny := x+dx, y+dy
		if 0 <= nx && nx < cols && 0 <= ny && ny < rows {
			tiles[x][y].Adjacent[Offset{dx, dy}] = tiles[nx][ny]
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

	return backing
}
