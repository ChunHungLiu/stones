package core

// MapGenBool generates Tiles from bool values to form various mazes.
type MapGenBool func(o Offset, pass bool) *Tile

// PerfectMaze creates a set of Tile which form a perfect maze (meaning the
// maze has no loops). The value of n specifies the size of the underlying graph
// describing the maze, which is related to but not equal to the number of Tile
// in the result maze. The runProb specifies how often the algorithm will try
// to continue extending a corridor, as opposed to starting a new branch.
func (f MapGenBool) PerfectMaze(n int, runProb, weaveProb float64) []*Tile {
	return applyMaze(abstractPerfect(n, runProb, weaveProb), f)
}

// BraidMaze creates a set of Tile which form a braid maze (meaning the maze
// has no dead ends). The value of n specifies the size of the underlying graph
// describing the maze, which is related to but not equal to the number of Tile
// in the result maze. The runProb specifies how often the algorithm will try
// to continue extending a corridor, as opposed to starting a new branch.
func (f MapGenBool) BraidMaze(n int, runProb, weaveProb float64) []*Tile {
	return applyMaze(abstractBraid(n, runProb, weaveProb, 1), f)
}

// HalfBraidMaze creates a set of Tile which form a half-braid maze (meaning the
// maze will have some dead ends and some loops). The value of n specifies the
// size of the underlying graph describing the maze, which is related to but not
// equal to the number of Tile in the result maze. The runProb specifies how
// often the algorithm will try to continue extending a corridor, as opposed to
// starting a new branch. The loopProb is the probability of removing a
// deadend by creating a loop.
func (f MapGenBool) HalfBraidMaze(n int, runProb, weaveProb, loopProb float64) []*Tile {
	return applyMaze(abstractBraid(n, runProb, weaveProb, loopProb), f)
}

// defaultMapGenBool is used in the generic versions of each MapGenBool method.
// It generates white '.' for passable Tile, and a white '#' for wall Tile.
var defaultMapGenBool = func(o Offset, pass bool) *Tile {
	t := NewTile(o)
	t.Pass = pass
	t.Lite = pass
	if !pass {
		t.Face = Glyph{'#', ColorWhite}
	}
	return t
}

// PerfectMaze creates a set of Tile which form a perfect maze (meaning the
// maze has no loops). The value of n specifies the size of the underlying graph
// describing the maze, which is related to but not equal to the number of Tile
// in the result maze. The runProb specifies how often the algorithm will try
// to continue extending a corridor, as opposed to starting a new branch.
// PerfectMazes uses a default MapGenBool which generates white '.' for passable
// Tile and white '#' for wall Tile.
func PerfectMaze(n int, runProb, weaveProb float64) []*Tile {
	return applyMaze(abstractPerfect(n, runProb, weaveProb), defaultMapGenBool)
}

// BraidMaze creates a set of Tile which form a braid maze (meaning the maze
// has no dead ends). The value of n specifies the size of the underlying graph
// describing the maze, which is related to but not equal to the number of Tile
// in the result maze. The runProb specifies how often the algorithm will try
// to continue extending a corridor, as opposed to starting a new branch.
// BraidMaze uses a default MapGenBool which generates white '.' for passable
// Tile and white '#' for wall Tile.
func BraidMaze(n int, runProb, weaveProb float64) []*Tile {
	return applyMaze(abstractBraid(n, runProb, weaveProb, 1), defaultMapGenBool)
}

// HalfBraidMaze creates a set of Tile which form a half-braid maze (meaning the
// maze will have some dead ends and some loops). The value of n specifies the
// size of the underlying graph describing the maze, which is related to but not
// equal to the number of Tile in the result maze. The runProb specifies how
// often the algorithm will try to continue extending a corridor, as opposed to
// starting a new branch. The loopProb is the probability of removing a
// deadend by creating a loop. HalfBraidMaze uses a default MapGenBool which
// generates white '.' for passable Tile and white '#' for wall Tile.
func HalfBraidMaze(n int, runProb, weaveProb, loopProb float64) []*Tile {
	return applyMaze(abstractBraid(n, runProb, weaveProb, loopProb), defaultMapGenBool)
}

// mazenode is a single node (in the graph sense) in an abstractmaze
type mazenode struct {
	Pos   Offset
	Edges map[Offset]*mazenode
}

// abstractmaze represents the graph structure of a maze
type abstractmaze struct {
	Nodes map[Offset][]*mazenode
}

// GetArbitraryNode returns a mazenode from the maze node list. Since the node
// is arbitrarily chosen, the caller cannot depend on the node being the same
// across multiple calls, nor can it be depended upon to be chosen randomly.
func (m *abstractmaze) GetArbitraryNode() *mazenode {
	for _, nodelist := range m.Nodes {
		for _, node := range nodelist {
			return node
		}
	}

	return nil
}

// Data needed by maze generation to iterate through directional Offsets.
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
func abstractPerfect(n int, runProb, weaveProb float64) *abstractmaze {
	// set up bookkeeping for growing tree algorithm
	origin := &mazenode{Offset{}, make(map[Offset]*mazenode)}
	maze := &abstractmaze{map[Offset][]*mazenode{origin.Pos: {origin}}}
	frontier := []*mazenode{origin}
	nodesAdded := 1

	// TODO Separate bookkeeping from node adding to create extendMaze

	// our frontier size is unbounded, so stop when we've added enough nodes.
	for nodesAdded < n {
		// select a node at random, meaning we emulate Prim's algorithm
		var index int
		if RandChance(runProb) {
			index = len(frontier) - 1
		} else {
			index = RandIntn(len(frontier))
		}
		curr := frontier[index]

		// candidate steps are on edges which lead to unseen node
		candidates := make([]Offset, 0, 4)
		for _, step := range orthogonal {
			_, used := curr.Edges[step]
			_, exists := maze.Nodes[curr.Pos.Add(step)]
			if !used && (!exists || RandChance(weaveProb)) {
				candidates = append(candidates, step)
			}
		}

		if len(candidates) > 0 {
			// create the adjacent node in the step direction
			step := candidates[RandIntn(len(candidates))]
			adjpos := curr.Pos.Add(step)
			adjnode := &mazenode{adjpos, make(map[Offset]*mazenode)}

			// link up the edges to and from the adjacent node
			maze.Nodes[adjpos] = append(maze.Nodes[adjpos], adjnode)
			curr.Edges[step] = adjnode
			adjnode.Edges[step.Neg()] = curr

			// add the newly created ajacent node to the frontier and inc size
			frontier = append(frontier, adjnode)
			nodesAdded++
		} else {
			// if we have no candidate edges, we'll never expand it
			frontier = append(frontier[:index], frontier[index+1:]...)
		}
	}

	return maze
}

// abstractPerfect generates an abstractmaze with the given number of nodes.
// This maze will be a braid, meaning that it has no deadends.
func abstractBraid(n int, runProb, weaveProb, loopProb float64) *abstractmaze {
	maze := abstractPerfect(n, runProb, weaveProb)
	removeDeadends(maze, loopProb)
	return maze
}

// removeDeadends removes a given percent of deadends from an abstractmaze.
// Deadends are removed by adding an edge to an unconnected but adjacent node.
// Deadends which have no unconnected adjacent node are simply removed.
func removeDeadends(m *abstractmaze, loopProb float64) {
	origin := m.GetArbitraryNode()

	// find all the dead ends - nodes which have only one edge
	deadends := []*mazenode{}
	frontier := []*mazenode{origin}
	visited := map[*mazenode]struct{}{origin: {}}
	for len(frontier) != 0 {
		curr := frontier[len(frontier)-1]
		frontier = frontier[:len(frontier)-1]

		if len(curr.Edges) == 1 {
			deadends = append(deadends, curr)
		}

		for _, adj := range curr.Edges {
			if _, seen := visited[adj]; !seen {
				frontier = append(frontier, adj)
				visited[adj] = struct{}{}
			}
		}
	}

	// remove all the dead ends by adding edges which connects deadends
	// some deadends must simply be removed, since they have no adjacent nodes
	for len(deadends) != 0 {
		deadend := deadends[len(deadends)-1]
		deadends = deadends[:len(deadends)-1]

		// check that deadend is still a deadend as it could have been used
		// as a connection for another deadend
		if len(deadend.Edges) > 1 {
			continue
		}

		if !RandChance(loopProb) {
			continue
		}

		// find all the adjancent nodes we connected to - deadend must have an
		// open edge towards the neighbor, and the neighbor needs an unused edge
		// back to the deadend.
		var candidates []*mazenode
		for _, step := range orthogonal {
			if _, used := deadend.Edges[step]; used {
				continue
			}

			negStep := step.Neg()
			for _, adj := range m.Nodes[deadend.Pos.Add(step)] {
				if _, used := adj.Edges[negStep]; !used {
					candidates = append(candidates, adj)
				}
			}
		}

		if len(candidates) == 0 {
			// since there was no valid edge to connect, we have a straggler
			// and we just delete the node
			m.Nodes[deadend.Pos] = remove(m.Nodes[deadend.Pos], deadend)

			// find the node adjacent to the deadend, and delete its edge
			// to the deadend. it will either be the next dead end to prune
			// or will be left alone if it has 2+ edges remaining.
			for step, adj := range deadend.Edges {
				delete(adj.Edges, step.Neg())
				if len(adj.Edges) == 1 {
					deadends = append(deadends, adj)
				}
			}
			// FIXME add a replacement node for each deleted node
		} else {
			// pick a neighboring node, and connect an edge to the neighbor
			neighbor := candidates[RandIntn(len(candidates))]
			step := neighbor.Pos.Sub(deadend.Pos)
			deadend.Edges[step] = neighbor
			neighbor.Edges[step.Neg()] = deadend
		}
	}
}

// remove returns a slice which is l with n removed.
func remove(l []*mazenode, n *mazenode) []*mazenode {
	for i, v := range l {
		if n == v {
			return append(l[:i], l[i+1:]...)
		}
	}
	return l
}

// applyMaze converts each node and edge an abstract maze to a single Tile.
// The origin Tile of the maze is returned.
func applyMaze(m *abstractmaze, f MapGenBool) []*Tile {
	origin := createPassTiles(m, f)
	connectDiagonals(origin)
	addWalls(origin, f)

	var maze []*Tile
	frontier := []*Tile{origin}
	visited := map[*Tile]struct{}{origin: {}}
	for len(frontier) != 0 {
		curr := frontier[len(frontier)-1]
		frontier = frontier[:len(frontier)-1]
		for _, adj := range curr.Adjacent {
			if _, seen := visited[adj]; !seen {
				frontier = append(frontier, adj)
				visited[adj] = struct{}{}
				maze = append(maze, adj)
			}
		}
	}
	return maze
}

// createPassTiles returns a set of passable Tiles (given by the origin Tile)
// corresponding to the nodes and edges of an abstractmaze.
func createPassTiles(m *abstractmaze, f MapGenBool) *Tile {
	origin := m.GetArbitraryNode()

	// graph traversal bookkeeping
	frontier := []*mazenode{origin}
	visited := map[*mazenode]*Tile{origin: f(Offset{}, true)}

	for len(frontier) != 0 {
		// pop the next node and grab its corresponding Tile
		node := frontier[len(frontier)-1]
		frontier = frontier[:len(frontier)-1]
		nodeTile := visited[node]

		// for each edge, create a Tile for the edge, and (if needed) a Tile for
		// the adjacent node
		for step, adj := range node.Edges {
			if _, edgeExists := nodeTile.Adjacent[step]; !edgeExists {
				negStep := step.Neg()

				// create a Tile corresponding to the edge between node and adj
				edgeOff := nodeTile.Offset.Add(step)
				edgeTile := f(edgeOff, true)
				nodeTile.Adjacent[step] = edgeTile
				edgeTile.Adjacent[negStep] = nodeTile

				// connect the edge Tile to a Tile corresponding to adj
				adjTile, seen := visited[adj]
				// The edge not existing is nessesary but sufficient for the
				// adj to have not been seen and created since we could have
				// seen the adj, but not popped it from the frontier yet. Thus,
				// this check is needed so we don't create the adj node twice.
				if !seen {
					adjTile = f(edgeOff.Add(step), true)
					// Since we found a new node, we should enqueue it now
					visited[adj] = adjTile
					frontier = append(frontier, adj)
				}
				edgeTile.Adjacent[step] = adjTile
				adjTile.Adjacent[negStep] = edgeTile
			}
		}
	}

	return visited[m.GetArbitraryNode()]
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

// addWalls connects each passable tile to a wall where needed. Walls are *not*
// connected, and their adjacency should never be used as a result.
func addWalls(origin *Tile, f MapGenBool) {
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
					wall = f(off, false)
				}
				// connect to the wall - note that this connection means that
				// nearby nodes will reuse this one when they are popped off
				// the queue in breadth first fashion.
				curr.Adjacent[step] = wall
				wall.Adjacent[step.Neg()] = curr
			}
		}
	}
}

// findWall queries the neighbors of origin for a Tile at the given dest.
func findWall(origin *Tile, dest Offset) (tile *Tile, ok bool) {
	for _, adj := range origin.Adjacent {
		adjStep := dest.Sub(adj.Offset)
		if tile, exists := adj.Adjacent[adjStep]; exists {
			return tile, true
		}
	}
	return nil, false
}

// TODO Add dungeon
// TODO Add caveify
// TODO Use writer interfaces instead of directly writing Tile

// TODO Add function to connect overworld with maze
// allow caller to specify entrance/exit criteria
// translate maze coordinates to match overworld coordinates
// fully connect maze entrance with overworld
// fully connect overworld connection neighbors with maze
