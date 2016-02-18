package core

type mazenode struct {
	Pos   Offset
	Edges map[Offset]*mazenode
}

type abstractmaze map[Offset]*mazenode

var (
	orthogonal = []Offset{
		{0, 1},
		{0, -1},
		{1, 0},
		{-1, 0},
	}
	cardinal = []Offset{
		{0, 1},
		{0, -1},
		{1, 0},
		{-1, 0},
		{1, 1},
		{-1, 1},
		{1, -1},
		{-1, -1},
	}
)

func abstractPerfect(n int) abstractmaze {
	start := &mazenode{Offset{}, make(map[Offset]*mazenode)}
	maze := abstractmaze{Offset{}: start}
	frontier := []*mazenode{start}

	for len(maze) < n {
		index := RandIntn(len(frontier))
		curr := frontier[index]

		candidates := make([]Offset, 0, 4)
		for _, step := range orthogonal {
			if _, seen := maze[curr.Pos.Add(step)]; !seen {
				candidates = append(candidates, step)
			}
		}

		if len(candidates) == 0 {
			frontier = append(frontier[:index], frontier[index+1:]...)
		} else {
			step := candidates[RandIntn(len(candidates))]
			adjpos := curr.Pos.Add(step)
			adjnode := &mazenode{adjpos, make(map[Offset]*mazenode)}

			maze[adjpos] = adjnode
			curr.Edges[step] = adjnode
			adjnode.Edges[step.Neg()] = curr

			frontier = append(frontier, adjnode)
		}
	}

	return maze
}

type tilemap map[Offset]*Tile

func (m tilemap) Get(o Offset) (t *Tile, newed bool) {
	if tile, ok := m[o]; ok {
		return tile, false
	}
	m[o] = NewTile(o)
	return m[o], true
}

func PerfectMaze(n int) map[Offset]*Tile {
	maze := abstractPerfect(n)
	tiles := make(tilemap)

	for off, node := range maze {
		nodeOff := off.Scale(2) // scale by 2, so we can fit edge tiles
		tile, _ := tiles.Get(nodeOff)

		for step := range node.Edges {
			negStep := step.Neg()

			// add a tile corresponding to the graph edge
			edgeOff := nodeOff.Add(step)
			edge, _ := tiles.Get(edgeOff)
			tile.Adjacent[step] = edge
			edge.Adjacent[negStep] = tile

			// add a tile corresponding to the adjacent node
			adjOff := edgeOff.Add(step)
			adj, _ := tiles.Get(adjOff)
			edge.Adjacent[step] = adj
			adj.Adjacent[negStep] = edge
		}
	}

	for off, tile := range tiles {
		if !tile.Pass {
			continue
		}
		for _, step := range cardinal {
			if _, ok := tile.Adjacent[step]; !ok {
				wall, newtile := tiles.Get(off.Add(step))
				if newtile {
					wall.Face = Glyph{'#', ColorWhite}
					wall.Pass = false
				}
				tile.Adjacent[step] = wall
				wall.Adjacent[step.Neg()] = tile
			}
		}
	}

	return tiles
}

// TODO Add braid and half-braid mazes
// TODO Add dungeon
// TODO Add caveify
