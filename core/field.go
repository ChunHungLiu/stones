package core

// Field is discrete potential field (or Dijkstra map) used for navigation.
type Field interface {
	Follow(*Tile) Offset
}

// sparseField implemenents Field using a sparse map of Tile weights.
type sparseField struct {
	weights map[*Tile]float64
}

// Follow returns an Offset from the given Tile which will lead to the
// neighboring tile which has a lower weight. Goal weights are negative
// so that the default value of 0 is neutral.
func (f *sparseField) Follow(t *Tile) Offset {
	bestWeight, bestOffset := f.weights[t], Offset{}

	for offset, adj := range t.Adjacent {
		if weight := f.weights[adj]; weight < bestWeight {
			bestWeight = weight
			bestOffset = offset
		}
	}

	return bestOffset
}

// AttractiveField computes a Field which pulls towards the goal Tile.
func AttractiveField(radius int, goal *Tile) Field {
	// setup Djkstra's algorithm bookkeeping
	weights := map[*Tile]float64{goal: float64(-radius)}
	queue := []*Tile{goal}

	// run Djkstra's algorithm to compute attractive weights from the goal
	for len(queue) > 0 {
		// pop the next Tile off the queue
		curr := queue[0]
		queue = queue[1:]

		// if we've reached edge field, stop expanding the field
		cost := weights[curr] + 1
		if cost > 0 {
			continue
		}

		// expand the frontier using neighbors of curr
		for _, adj := range curr.Adjacent {
			if _, seen := weights[adj]; !seen && adj.Pass {
				weights[adj] = cost
				queue = append(queue, adj)
			}
		}
	}

	return &sparseField{weights}
}

// TODO RepulsiveField
// TODO RandomField
