package core

type MapGen func(o Offset) *Tile

func createTileGrid(cols, rows int, o Offset, f MapGen) []*Tile {
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

// isDiag returns true if the Offset is a single diagonal step.
func isDiag(o Offset) bool {
	return Abs(o.X) == 1 && Abs(o.Y) == 1
}

// connectDiagonals takes an orthogonally connected maze, and connects each Tile
// diagonally through its neighbors.
func connectDiagonals(origin *Tile) {
	// FIXME only connect diagonals on the same z-level
	// setup breadth-first graph traversal bookkeeping
	frontier := []*Tile{origin}
	visited := map[*Tile]struct{}{origin: {}}

	for len(frontier) != 0 {
		// pop the queue in breadth first fashion
		curr := frontier[0]
		frontier = frontier[1:]

		for _, step1 := range orthogonal {
			// take two orthogonal steps to for a single diagonal step from curr
			// if there is something there, connect curr and the resulting Tile
			if adj, ok := curr.Adjacent[step1]; ok {
				for step2, tile := range adj.Adjacent {
					if diag := step1.Add(step2); isDiag(diag) {
						curr.Adjacent[diag] = tile
						tile.Adjacent[diag.Neg()] = curr
					}
				}

				// enqueue any unvisted Tile
				if _, seen := visited[adj]; !seen {
					frontier = append(frontier, adj)
					visited[adj] = struct{}{}
				}
			}
		}
	}
}

// addWalls connects each passable tile to a wall where needed. Any added wall
// Tile will be in the returned Tile slice. The provided MapGen will be used to
// create any new wall Tile, so any wrapped MapGens should be created
// accordingly. The wall Tiles will be fully connected, but those connection may
// not be bi-directional. Consequently, the wall Tile adjancency maps should not
// be used for field of view purposes. The returned slice of Tile will contain
// any added wall Tile.
func addWalls(origin *Tile, f MapGen) []*Tile {
	var walls []*Tile

	// setup breadth-first graph traversal bookkeeping
	frontier := []*Tile{origin}
	visited := map[*Tile]struct{}{origin: {}}

	for len(frontier) != 0 {
		// pop queue in breadth first fashion
		curr := frontier[0]
		frontier = frontier[1:]

		// for each step, ensure we have an edge.
		for _, step := range cardinal {
			if adj, ok := curr.Adjacent[step]; ok {
				// if the edge already exists, just add it if passable.
				if _, seen := visited[adj]; !seen && adj.Pass {
					frontier = append(frontier, adj)
					visited[adj] = struct{}{}
				}
			} else {
				// try and find a pre-made wall through my neighbors
				off := curr.Offset.Add(step)
				wall, ok := findWall(curr, off)
				if !ok {
					wall = f(off)
					walls = append(walls, wall)
				}
				// connect to the wall - note that this connection means that
				// nearby nodes will reuse this one when they are popped off
				// the queue in breadth first fashion.
				curr.Adjacent[step] = wall
				wall.Adjacent[step.Neg()] = curr
			}
		}
	}

	return walls
}
