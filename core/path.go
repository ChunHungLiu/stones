package core

import (
	"container/heap"
	"math"
)

// DistFn computes a cost or heuristic between two Tiles.
type DistFn func(*Tile, *Tile) float64

// score tracks the f-score and path for a particular Tile
type score struct {
	GCost float64
	HEst  float64
	Prev  *Tile
}

// FScore computes the sum of the g-cost and h-estimate.
func (s *score) FScore() float64 {
	return s.GCost + s.HEst
}

// scorer tracks the scores for a set of nodes.
type scorer struct {
	nodes     map[*Tile]*score
	heuristic DistFn
	goal      *Tile
}

// newscorer creates a new scorer, with the origin having a 0 cost.
func newscorer(origin, goal *Tile, heuristic DistFn) *scorer {
	s := &scorer{make(map[*Tile]*score), heuristic, goal}
	s.Score(origin).GCost = 0
	return s
}

// Score gets the score for a particular tile. The scores are lazily
// instantiated so that the graph can be arbitrarily large and the scorer is
// still memory efficient.
func (s *scorer) Score(t *Tile) *score {
	// if we've already initialized the score, just return it
	if score, ok := s.nodes[t]; ok {
		return score
	}

	// initialize the score with an infinite cost (so that any path at all is
	// accepted), as well as precomputing the heuristic.
	score := &score{math.Inf(1), s.heuristic(t, s.goal), nil}
	s.nodes[t] = score
	return score
}

// Path recursively constructs the path from the given Tile to the origin.
func (s *scorer) Path(t *Tile) []*Tile {
	prev := s.Score(t).Prev

	// if we are at the origin, the path to the origin is empty.
	if prev == nil {
		return nil
	}

	// add ourselves to the path from the origin to our prev
	return append(s.Path(prev), t)
}

// tilequeue implements heap.Interface, using the FScore to sort.
type tilequeue struct {
	queue  []*Tile
	scores *scorer
}

// Len returns the number of Tiles in the queue.
func (q *tilequeue) Len() int {
	return len(q.queue)
}

// Less compares the FScore of the ith and jth Tiles, and returns true if the
// score for the ith Tile is less than that of the jth Tile.
func (q *tilequeue) Less(i, j int) bool {
	return q.scores.Score(q.queue[i]).FScore() < q.scores.Score(q.queue[j]).FScore()
}

// Swap switches the values of the ith and jth Tile in the queue.
func (q *tilequeue) Swap(i, j int) {
	q.queue[i], q.queue[j] = q.queue[j], q.queue[i]
}

// Push pushes a *Tile onto the queue, panicing if the data is not a *Tile.
func (q *tilequeue) Push(x interface{}) {
	q.queue = append(q.queue, x.(*Tile))
}

// Pop removes and returns the last *Tile in the queue as an interface{}.
func (q *tilequeue) Pop() interface{} {
	n := len(q.queue) - 1
	x := q.queue[n]
	q.queue = q.queue[:n]
	return x
}

// GraphSearch performs a generic graph search from the origin to the goal
// using the given heuristic and cost. If the heuristic is admissible, meaning
// it never underestimates the final path cost, then the resulting path will be
// optimal with respect to cost.
func GraphSearch(origin, goal *Tile, cost, heuristic DistFn) []*Tile {
	scores := newscorer(origin, goal, heuristic)
	frontier := &tilequeue{[]*Tile{origin}, scores}
	closed := make(map[*Tile]struct{})

	for frontier.Len() > 0 {
		// get the next tile to explore, skip if we've already closed it
		curr := heap.Pop(frontier).(*Tile)
		if _, seen := closed[curr]; seen {
			continue
		}

		// mark current as seen
		closed[curr] = struct{}{}

		// if we find the goal, we've already found the best path
		if curr == goal {
			return scores.Path(goal)
		}

		// for each neighbor, see if we've found a better path, then enqueue it
		currscore := scores.Score(curr)
		for _, adj := range curr.Adjacent {
			if !adj.Pass {
				continue
			}

			if _, seen := closed[adj]; !seen {
				// compute the cost of that path to adj going through curr
				cost := currscore.GCost + cost(curr, adj)

				// we found a better path for the adjacent tile
				if adjscore := scores.Score(adj); cost < adjscore.GCost {
					adjscore.GCost = cost
					adjscore.Prev = curr
					heap.Push(frontier, adj)
				}
			}
		}
	}

	// if we exhaust the frontier, and didn't find the goal, there is no path
	return nil
}

// NewGraphSearch creates a GraphSearch function with the given DistFns.
func NewGraphSearch(cost, heuristic DistFn) func(*Tile, *Tile) []*Tile {
	return func(a, b *Tile) []*Tile {
		return GraphSearch(a, b, cost, heuristic)
	}
}

// multiplying heuristic values by this constant should still lead to an
// admissible heuristic, but will hopefully shorten search times and yield more
// direct looking paths by favoring nodes with smaller heuristics.
const tiebreak = 1 + 1e-10

// euclidean computes the tiebreak version of Euclidean distance between nodes.
func euclidean(a, b *Tile) float64 {
	return b.Offset.Sub(a.Offset).Euclidean() * tiebreak
}

// zero is a DistFn which simply returns 0.
func zero(*Tile, *Tile) float64 {
	return 0
}

// AStar computes a minimum cost path between two Tiles.
func AStarPath(origin, goal *Tile) []*Tile {
	return GraphSearch(origin, goal, euclidean, euclidean)
}

// Greedy computes a greedy path between two Tiles.
func GreedyPath(origin, goal *Tile) []*Tile {
	return GraphSearch(origin, goal, zero, euclidean)
}
