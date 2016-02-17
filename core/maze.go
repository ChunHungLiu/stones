package core

type Mazenode struct {
	Pos   Offset
	Edges map[Offset]*Mazenode
}

type Abstractmaze map[Offset]*Mazenode

var orthogonal = []Offset{
	{0, 1},
	{0, -1},
	{1, 0},
	{-1, 0},
}

func NewPerfect(n int) Abstractmaze {
	start := &Mazenode{Offset{}, make(map[Offset]*Mazenode)}
	maze := Abstractmaze{Offset{}: start}
	frontier := []*Mazenode{start}

	for len(maze) < n {
		// TODO Fix selection strategy (and removal), so we get good maze shape
		curr := frontier[0]

		candidates := make([]Offset, 0, 4)
		for _, step := range orthogonal {
			if _, seen := maze[curr.Pos.Add(step)]; !seen {
				candidates = append(candidates, step)
			}
		}

		if len(candidates) == 0 {
			frontier = frontier[1:]
		} else {
			step := candidates[RandIntn(len(candidates))]
			adjpos := curr.Pos.Add(step)
			adjnode := &Mazenode{adjpos, make(map[Offset]*Mazenode)}

			maze[adjpos] = adjnode
			curr.Edges[step] = adjnode
			adjnode.Edges[step.Neg()] = curr

			frontier = append(frontier, adjnode)
		}
	}

	return maze
}

// TODO Add perfect maze
// TODO Add braid maze
// TODO Add dungeon
// TODO Add caveify
