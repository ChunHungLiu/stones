package core

import (
	"math"
)

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
	minWeight, minOffset := f.weights[t], Offset{}

	for offset, adj := range t.Adjacent {
		if weight := f.weights[adj]; weight < minWeight {
			minWeight = weight
			minOffset = offset
		}
	}

	return minOffset
}

// computeAttractWeights computes the weights of a sparsefield which pull
// towards the given goals. The edge weigts will be 0, with the goals
// having a weight of -radius.
func computeAttractWeights(radius int, goals []*Tile) map[*Tile]float64 {
	// setup Djkstra's algorithm bookkeeping
	weights := make(map[*Tile]float64)
	queue := make([]*Tile, len(goals))
	for i, goal := range goals {
		weights[goal] = float64(-radius)
		queue[i] = goal
	}

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

	return weights
}

// AttractiveField computes a Field which pulls towards the goal Tile.
func AttractiveField(radius int, goals ...*Tile) Field {
	return &sparseField{computeAttractWeights(radius, goals)}
}

// ReplusiveField creates a Field which pulls towards the outermost edge of the
// field with the given ungoals as the sources. This is *not* the same as
// negating the weights of an attractive field, as the path towards the edge
// of the field may require a step towards an ungoal.
func ReplusiveField(radius int, ungoals ...*Tile) Field {
	attractWeights := computeAttractWeights(radius, ungoals)

	// compute the weight of the edge of the attractive field
	edgeWeight := math.Inf(-1)
	for _, weight := range attractWeights {
		edgeWeight = math.Max(edgeWeight, weight)
	}

	// create bookkeeping for djiksta's algorithm again
	weights := make(map[*Tile]float64)
	var queue []*Tile
	for tile, weight := range attractWeights {
		if weight == edgeWeight {
			weights[tile] = 0
			queue = append(queue, tile)
		}
	}

	// perform djikstra algorith, with the restriction that we stay in the
	// bounds of the original attractive field.
	for len(queue) > 0 {
		// pop the next Tile off the queue
		curr := queue[0]
		queue = queue[1:]

		cost := weights[curr] + 1
		for _, adj := range curr.Adjacent {
			_, seen := weights[adj]
			_, keep := attractWeights[adj]
			// only consider unseen nodes which are in the attractive field.
			if !seen && keep {
				weights[adj] = cost
				queue = append(queue, adj)
			}
		}
	}

	return &sparseField{weights}
}

// funcField is a Filed which is composed of only a single function call.
type funcField func(*Tile) Offset

// Follow simply calls the underlying field function.
func (f funcField) Follow(t *Tile) Offset {
	return f(t)
}

// randField is the underlying function for RandomField.
func randField(t *Tile) Offset {
	candidates := make([]Offset, 0, len(t.Adjacent))
	for offset, adj := range t.Adjacent {
		if adj.Pass {
			candidates = append(candidates, offset)
		}
	}
	if len(candidates) == 0 {
		return Offset{}
	}
	return candidates[RandIntn(len(candidates))]
}

// RandomField is a Field which generates random Offsets. The resulting Offset
// will always corespond to an adjacent Tile which is passable.
func RandomField() Field {
	return funcField(randField)
}
