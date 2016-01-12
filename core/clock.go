package core

import (
	"math"
)

// deltanode stores Entity events for a particular delta in a DeltaClock.
type deltanode struct {
	delta  float64
	link   *deltanode
	events map[Entity]struct{}
}

// DeltaClock implements a data structure which allows for fast scheduling.
// See http://stonesrl.blogspot.com/2013/02/delta-clock.html for more info.
type DeltaClock struct {
	head  *deltanode
	nodes map[Entity]*deltanode
}

// NewDeltaClock creates an empty DeltaClock.
func NewDeltaClock() *DeltaClock {
	return &DeltaClock{nil, make(map[Entity]*deltanode)}
}

// Schedule adds an Entity to the queue at the given delta. Note that the delta
// is split into its integer and fractional part. The integer part is used to
// determine the amount of delay, while the fractional part is only used to
// ensure a unique scheduling delta.
//
// As an example, suppose we repeatedly schedule event A with deltas of 1,
// event B with deltas of 1.5, and event c with deltas of 2. It is *not* the
// case that A will fire 3 times for every 2 times B fires. Instead, both A and
// B will fire at the same rate, but A will always go first as it has a lower
// fractional part in its delta. It *is* the case that A and B will fire twice
// as often as C.
func (c *DeltaClock) Schedule(e Entity, delta float64) {
	var prev, curr *deltanode = nil, c.head

	// iterate over nodes, ensuring we haven't gone passed the end,
	// or passed the desired node
	for curr != nil && delta > curr.delta {
		delta -= math.Trunc(curr.delta)
		prev, curr = curr, curr.link
	}

	var node *deltanode
	if curr != nil && delta == curr.delta {
		// if the desired node already exists, just reuse it
		node = curr
	} else {
		// desired didnt exist, so create it, with a link to curr node
		node = &deltanode{delta, curr, make(map[Entity]struct{})}

		if prev == nil {
			// prev == nil iff we're at the beginning of the list
			c.head = node
		} else {
			// otherwise, insert the node in the middle of this list
			prev.link = node
		}

		// the next node needs to take the new node's delta into account
		if curr != nil {
			curr.delta -= math.Trunc(delta) // only subtract schedule time
		}
	}

	// ad the event to the node
	node.events[e] = struct{}{}
	c.nodes[e] = node
}

// Unschedule removes an Entity from the queue. If the Entity is not in the
// queue, no action is taken.
func (c *DeltaClock) Unschedule(e Entity) {
	if node, ok := c.nodes[e]; ok {
		delete(node.events, e)
		delete(c.nodes, e)
	}
}

// Advance returns the next set of Entity scheduled. If all Events have been
// removed from the next delta, then the set will be empty. If no deltas are
// remaining, then the result is nil.
func (c *DeltaClock) Advance() map[Entity]struct{} {
	// no events were scheduled
	if c.head == nil {
		return nil
	}

	// delete the events from the scheduler, since the head is being removed
	for e := range c.head.events {
		delete(c.nodes, e)
	}

	// gets the events to return, and remove the head node
	events := c.head.events
	c.head = c.head.link
	return events
}
