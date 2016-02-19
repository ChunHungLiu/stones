package core

// mazenode is a single node (in the graph sense) in an abstractmaze
type mazenode struct {
	Pos   Offset
	Edges map[Offset]*mazenode
}

// abstractmaze represents the graph structure of a maze
type abstractmaze map[Offset]*mazenode

// Data needed by abstractPerfect to iterate through directional Offsets.
var (
	orthogonal = [4]Offset{
		{0, 1},
		{0, -1},
		{1, 0},
		{-1, 0},
	}
	cardinal = [8]Offset{
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

// abstractPerfect generates an abstractmaze with the given number of nodes.
// This maze will be perfect, meaning that it has no loops.
func abstractPerfect(n int, runfactor float64) abstractmaze {
	// set up bookkeeping for growing tree algorithm
	start := &mazenode{Offset{}, make(map[Offset]*mazenode)}
	maze := abstractmaze{Offset{}: start}
	frontier := []*mazenode{start}

	// our frontier size is unbounded, so stop when we've added enough nodes.
	for len(maze) < n {
		// select a node at random, meaning we emulate Prim's algorithm
		var index int
		if RandChance(runfactor) {
			index = len(frontier) - 1
		} else {
			index = RandIntn(len(frontier))
		}
		curr := frontier[index]

		// candidate steps are on edges which lead to unseen node
		candidates := make([]Offset, 0, 4)
		for _, step := range orthogonal {
			if _, seen := maze[curr.Pos.Add(step)]; !seen {
				candidates = append(candidates, step)
			}
		}

		if len(candidates) > 0 {
			// create the adjacent node in the step direction
			step := candidates[RandIntn(len(candidates))]
			adjpos := curr.Pos.Add(step)
			adjnode := &mazenode{adjpos, make(map[Offset]*mazenode)}

			// link up the edges to and from the adjacent node
			maze[adjpos] = adjnode
			curr.Edges[step] = adjnode
			adjnode.Edges[step.Neg()] = curr

			// add the newly created ajacent node to the frontier
			frontier = append(frontier, adjnode)
		} else {
			// if we have no candidate edges, we'll never expand it
			frontier = append(frontier[:index], frontier[index+1:]...)
		}
	}

	return maze
}

// tilemap allows for lazy instantiation of Tile
type tilemap map[Offset]*Tile

// Get retrieves the Tile at the given offset, creating it if needed.
// The Tile is returned, along with a bool indicating if the Tile was newly
// created (as opposed to existing due to a previous call to Get).
func (m tilemap) Get(o Offset) (t *Tile, justcreated bool) {
	if tile, ok := m[o]; ok {
		return tile, false
	}
	m[o] = NewTile(o)
	return m[o], true
}

// PerfectMaze creates a set of Tile which form a perfect maze (meaning the
// maze has no loops). Note that n specifies the size of the underlying graph
// describing the maze, which is related to but equal to the number of Tile in
// the resulting maze. The runfactor specifies how often the algorithm will try
// to continue extending a corridor, as opposed to starting a new branch.
func PerfectMaze(n int, runfactor float64) map[*Tile]struct{} {
	// create an abstractmaze, which we will convert into a grid of Tile
	maze := abstractPerfect(n, runfactor)
	grid := make(tilemap)

	// create a Tile for every node and edge in the maze
	for off, node := range maze {
		// create the Tile corresponding to the graph node
		// scale node Offset by 2 so we can fit an edge Tile between node Tiles
		nodeOff := off.Scale(2)
		tile, _ := grid.Get(nodeOff)

		// for each edge, create the adajcent edge and node Tiles
		for step := range node.Edges {
			negStep := step.Neg()

			// add a tile corresponding to the graph edge
			edgeOff := nodeOff.Add(step)
			edge, _ := grid.Get(edgeOff)
			tile.Adjacent[step] = edge
			edge.Adjacent[negStep] = tile

			// add a tile corresponding to the adjacent node
			adjOff := edgeOff.Add(step)
			adj, _ := grid.Get(adjOff)
			edge.Adjacent[step] = adj
			adj.Adjacent[negStep] = edge
		}
	}

	// surround each passable tile with a wall tile
	for off, tile := range grid {
		if !tile.Pass {
			continue
		}
		for _, step := range cardinal {
			if _, ok := tile.Adjacent[step]; !ok {
				wall, newtile := grid.Get(off.Add(step))
				if newtile {
					wall.Face = Glyph{'#', ColorWhite}
					wall.Pass = false
				}
			}
		}
	}

	// add in the missing tile connections
	for off, tile := range grid {
		for _, step := range cardinal {
			_, link := tile.Adjacent[step]
			neighbor := grid[off.Add(step)]
			if !link && neighbor != nil {
				tile.Adjacent[step] = neighbor
				neighbor.Adjacent[step.Neg()] = tile
			}
		}
	}

	// convert and return tilemap as a set of tiles
	tiles := make(map[*Tile]struct{})
	for _, tile := range grid {
		tiles[tile] = struct{}{}
	}
	return tiles
}

// TODO Add optional z-levels to mazes
// TODO Add braid and half-braid mazes
// TODO Add dungeon
// TODO Add caveify
