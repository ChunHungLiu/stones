package core

import (
	"fmt"
)

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
			adj := curr.Pos.Add(step)
			if _, seen := maze[adj]; !seen {
				candidates = append(candidates, adj)
			}
		}
		fmt.Println(curr.Pos, candidates)

		if len(candidates) == 0 {
			frontier = frontier[1:]
		} else {
			adj := candidates[RandIntn(len(candidates))]
			// TODO clean up the edge insertion
			maze[adj] = &Mazenode{adj, make(map[Offset]*Mazenode)}
			maze[adj].Edges[curr.Pos.Sub(adj)] = curr
			curr.Edges[adj.Sub(curr.Pos)] = maze[adj]
			frontier = append(frontier, maze[adj])
		}
	}

	return maze
}

// TODO Add perfect maze
// TODO Add braid maze
// TODO Add dungeon
// TODO Add caveify
