package core

// mazenode is a single node (in the graph sense) in an abstractmaze
type mazenode struct {
	Pos   Offset
	Edges map[Offset]*mazenode
}

// abstractmaze represents the graph structure of a maze
type abstractmaze struct {
	Origin *mazenode
	Nodes  map[Offset][]*mazenode
}

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
func abstractPerfect(n int, runProb, weaveProb float64) abstractmaze {
	// set up bookkeeping for growing tree algorithm
	origin := &mazenode{Offset{}, make(map[Offset]*mazenode)}
	maze := abstractmaze{origin, map[Offset][]*mazenode{origin.Pos: {origin}}}
	frontier := []*mazenode{origin}
	nodesAdded := 1

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
func abstractBraid(n int, runProb, weaveProb, loopProb float64) abstractmaze {
	maze := abstractPerfect(n, runProb, weaveProb)
	removeDeadends(maze, loopProb)
	return maze
}

// removeDeadends removes a given percent of deadends from an abstractmaze.
// Deadends are removed by adding an edge to an unconnected but adjacent node.
// Deadends which have no unconnected adjacent node are simply removed.
func removeDeadends(m abstractmaze, loopProb float64) {
	// find all the dead ends - nodes which have only one edge
	deadends := []*mazenode{}
	frontier := []*mazenode{m.Origin}
	visited := map[*mazenode]struct{}{m.Origin: {}}
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
	for _, deadend := range deadends {
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

			negstep := step.Neg()
			for _, adj := range m.Nodes[deadend.Pos.Add(step)] {
				if _, used := adj.Edges[negstep]; !used {
					candidates = append(candidates, adj)
				}
			}
		}

		if len(candidates) == 0 {
			// since there was no valid edge to connect, we have a straggler
			// we just delete nodes until we no longer have a dead end
			for len(deadend.Edges) == 1 {

				// delete(m.Nodes, deadend.Pos) // FIXME delete *only* the deadend

				// find the node adjacent to the deadend, and delete its edge
				// to the deadend. it will either be the next dead end to prune
				// or will be left alone if it has 2+ edges remaining.
				for step, adj := range deadend.Edges {
					delete(adj.Edges, step.Neg())
					deadend = adj
					break // there is only one adjacent node to a deadend
				}
			}
		} else {
			// pick a neighboring node, and connect an edge to the neighbor
			neighbor := candidates[RandIntn(len(candidates))]
			step := neighbor.Pos.Sub(deadend.Pos)
			deadend.Edges[step] = neighbor
			neighbor.Edges[step.Neg()] = deadend
		}
	}
}

// PerfectMaze creates a set of Tile which form a perfect maze (meaning the
// maze has no loops). The value of n specifies the size of the underlying graph
// describing the maze, which is related to but not equal to the number of Tile
// in the result maze. The runProb specifies how often the algorithm will try
// to continue extending a corridor, as opposed to starting a new branch.
func PerfectMaze(n int, runProb, weaveProb float64) *Tile {
	return applyMaze(abstractPerfect(n, runProb, weaveProb))
}

// BraidMaze creates a set of Tile which form a braid maze (meaning the maze
// has no dead ends). The value of n specifies the size of the underlying graph
// describing the maze, which is related to but not equal to the number of Tile
// in the result maze. The runProb specifies how often the algorithm will try
// to continue extending a corridor, as opposed to starting a new branch.
func BraidMaze(n int, runProb, weaveProb float64) *Tile {
	return applyMaze(abstractBraid(n, runProb, weaveProb, 1))
}

// HalfBraidMaze creates a set of Tile which form a half-braid maze (meaning the
// maze will have some dead ends and some loops). The value of n specifies the
// size of the underlying graph describing the maze, which is related to but not
// equal to the number of Tile in the result maze. The runProb specifies how
// often the algorithm will try to continue extending a corridor, as opposed to
// starting a new branch. The loopProb is the probability of removing a
// deadend by creating a loop.
func HalfBraidMaze(n int, runProb, weaveProb, loopProb float64) *Tile {
	return applyMaze(abstractBraid(n, runProb, weaveProb, loopProb))
}

// applyMaze converts each node and edge an abstract maze to a single Tile.
// The origin Tile of the maze is returned.
func applyMaze(m abstractmaze) *Tile {
	origin := createPassTiles(m)
	connectDiagonals(origin)
	addWalls(origin)
	return origin
}

// createPassTiles returns a set of passable Tiles (given by the origin Tile)
// corresponding to the nodes and edges of an abstractmaze.
func createPassTiles(m abstractmaze) *Tile {
	// graph traversal bookkeeping
	frontier := []*mazenode{m.Origin}
	visited := map[*mazenode]*Tile{m.Origin: NewTile(m.Origin.Pos.Scale(2))}

	for len(frontier) != 0 {
		// pop the next node and grab its corresponding Tile
		node := frontier[len(frontier)-1]
		frontier = frontier[:len(frontier)-1]
		nodeTile := visited[node]

		for step, adj := range node.Edges {
			if _, edgeExists := nodeTile.Adjacent[step]; !edgeExists {
				negStep := step.Neg()

				// create a Tile corresponding to the edge between node and adj
				edgeOff := nodeTile.Offset.Add(step)
				edgeTile := NewTile(edgeOff)
				nodeTile.Adjacent[step] = edgeTile
				edgeTile.Adjacent[negStep] = nodeTile

				// connect the edge Tile to a Tile corresponding to adj
				adjTile, seen := visited[adj]
				// The edge not existing is nessesary but sufficient for the
				// adj to have not been seen and created since we could have
				// seen the adj, but not popped it from the frontier yet. Thus,
				// this check is needed so we don't create the adj node twice.
				if !seen {
					adjTile = NewTile(edgeOff.Add(step))
					// Since we found a new node, we should enqueue it now
					visited[adj] = adjTile
					frontier = append(frontier, adj)
				}
				edgeTile.Adjacent[step] = adjTile
				adjTile.Adjacent[negStep] = edgeTile
			}
		}
	}

	return visited[m.Origin]
}

func isDiag(o Offset) bool {
	return Abs(o.X) == 1 && Abs(o.Y) == 1
}

// connectDiagonals takes an orthogonally connected maze, and connects each Tile
// diagonally through its neighbors.
func connectDiagonals(origin *Tile) {
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

func addWalls(origin *Tile) {
	frontier := []*Tile{origin}
	visited := map[*Tile]struct{}{origin: {}}

	for len(frontier) != 0 {
		// pop in breadth first fashion
		curr := frontier[0]
		frontier = frontier[1:]

		for _, step := range cardinal {
			if adj, ok := curr.Adjacent[step]; ok {
				if _, seen := visited[adj]; !seen && adj.Pass {
					frontier = append(frontier, adj)
					visited[adj] = struct{}{}
				}
			} else {
				off := curr.Offset.Add(step)
				wall, ok := findWall(curr, off)
				if !ok {
					wall = newWall(off)
				}
				curr.Adjacent[step] = wall
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

func newWall(o Offset) *Tile {
	wall := NewTile(o)
	wall.Pass = false
	wall.Face = Glyph{'#', ColorWhite}
	return wall
}

// TODO Add dungeon
// TODO Add caveify
// TODO Use writer interfaces instead of directly writing Tile
