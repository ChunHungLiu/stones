package core

import (
	"fmt"
)

// We use these tables to cheaply approximate FOV, but we cache the tables so
// we only have to compute them once.
var tableCache = make(map[int]map[Offset]map[Offset]struct{})

// FoV uses a simple heuristic to approximate shadowcasting field of view
// calculation. The offsets in the resulting field are reletive to the given
// origin.
func FoV(origin *Tile, radius int) map[Offset]*Tile {
	// Retrieve (or create and cache) the table for the given radius.
	// This table maps a particular offset to a set of offsets which can seen
	// if the given one is transparent. Using this table, we basically just do
	// a recursive search using the table to guide us. Thus, we get a field of
	// view algorithm which performs minimal computation, never revisits tiles,
	// and short circuits on closed maps.
	table, cached := tableCache[radius]
	if !cached {
		table = computeTable(radius)
		tableCache[radius] = table
		fmt.Println(table)
	}

	fov := map[Offset]*Tile{Offset{0, 0}: origin}
	stack := []Offset{{0, 0}}

	for len(stack) > 0 {
		// Pop an offset from the stack
		off := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		// Get the tile  for that offset
		tile := fov[off]

		for adj := range table[off] {
			// Add all the adjacent tiles to the field of view.
			neighbor := tile.Adjacent[adj.Diff(off)]
			fov[adj] = neighbor

			// If the neighbor is passable, push it onto the stack to continue
			// exploration. Since we already added it to fov, when we pop it,
			// we'll be able to access the position again.
			if neighbor.Pass {
				stack = append(stack, adj)
			}
		}
	}

	// fix some artifacts related to standing next to a long wall.
	wallfix(fov, radius)
	return fov
}

// computeTable gets the table for a particular radius. This table will allow
// us to approxmiate shadowcasting using FoV.
func computeTable(radius int) map[Offset]map[Offset]struct{} {
	table := make(map[Offset]map[Offset]struct{})

	// We start at the origin, and will compute a single octant.
	addEntry(table, Offset{0, 0}, Offset{1, 0})
	addEntry(table, Offset{0, 0}, Offset{1, 1})

	// The following algorithm is better described in the blog post at:
	// http://stonesrl.blogspot.com/2013/02/pre-computed-fov.html
	// Basically there is a pattern in which tiles spawn both diagonally and
	// horizontally. Each row, the distance between these tiles increases by 1.
	// Everything below such a tile continues diagoanlly, Everything else goes
	// horizontally. A picture is worth a thousand words, so check out the
	// blog post...
	currBreak := 0
	breakCount := 0
	for x := 1; x < radius; x++ {
		nextY := 0
		for y := 0; y <= x; y++ {
			pos := Offset{x, y}
			if y == currBreak {
				addEntry(table, pos, Offset{x + 1, nextY})
				addEntry(table, pos, Offset{x + 1, nextY + 1})
				nextY += 2
			} else {
				addEntry(table, pos, Offset{x + 1, nextY})
				nextY++
			}
		}
		breakCount--
		if breakCount < 0 {
			breakCount = currBreak + 1
			currBreak++
		}
	}

	// Now that we've computed one octant, reflect and rotate to complete the
	// other 7 octants.
	completeTable(table)

	return table
}

// addEntry places a link between two offsets, adding the set keyed by src
// if it is not already present.
func addEntry(table map[Offset]map[Offset]struct{}, src, dst Offset) {
	neighbors, ok := table[src]
	if !ok {
		neighbors = map[Offset]struct{}{}
		table[src] = neighbors
	}
	neighbors[dst] = struct{}{}
}

// completeTable uses reflection and rotation to take a table with a single
// octant and extend it to all 8 octants.
func completeTable(table map[Offset]map[Offset]struct{}) {
	for key := range table {
		from := Offset{key.Y, key.X}
		for pos := range table[key] {
			addEntry(table, from, Offset{pos.Y, pos.X})
		}
	}

	for key := range table {
		negX := Offset{-key.X, key.Y}
		negY := Offset{key.X, -key.Y}
		negXY := Offset{-key.X, -key.Y}
		for pos := range table[key] {
			addEntry(table, negX, Offset{-pos.X, pos.Y})
			addEntry(table, negY, Offset{pos.X, -pos.Y})
			addEntry(table, negXY, Offset{-pos.X, -pos.Y})
		}
	}
}

// wallfix fills in some missing wall artifacts in a field of view.
func wallfix(fov map[Offset]*Tile, radius int) {
	for dx := 0; dx <= radius; dx++ {
		if pos, ok := fov[Offset{dx, 0}]; ok {
			fov[Offset{dx, 1}] = pos.Adjacent[Offset{0, 1}]
			fov[Offset{dx, -1}] = pos.Adjacent[Offset{0, -1}]
		} else {
			break
		}
	}
	for dx := 0; dx >= -radius; dx-- {
		if pos, ok := fov[Offset{dx, 0}]; ok {
			fov[Offset{dx, 1}] = pos.Adjacent[Offset{0, 1}]
			fov[Offset{dx, -1}] = pos.Adjacent[Offset{0, -1}]
		} else {
			break
		}
	}
	for dy := 0; dy <= radius; dy++ {
		if pos, ok := fov[Offset{0, dy}]; ok {
			fov[Offset{1, dy}] = pos.Adjacent[Offset{1, 0}]
			fov[Offset{-1, dy}] = pos.Adjacent[Offset{-1, 0}]
		} else {
			break
		}
	}
	for dy := 0; dy >= -radius; dy-- {
		if pos, ok := fov[Offset{0, dy}]; ok {
			fov[Offset{1, dy}] = pos.Adjacent[Offset{1, 0}]
			fov[Offset{-1, dy}] = pos.Adjacent[Offset{-1, 0}]
		} else {
			break
		}
	}
}
