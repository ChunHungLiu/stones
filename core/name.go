package core

import (
	"strings"
)

// boundary is a special marker for beginning and end of names.
var boundary = "\x00"

// counter implements a categorical distribution over byte events. Event counts
// are stored unnormalized, but we also store the count total so we can
// normalize when needed.
type counter struct {
	Counts map[byte]float64
	Total  float64
}

// newCounter creates a new counter with a Dirichlet prior over the support.
func newCounter(support map[byte]struct{}, prior float64) *counter {
	counts := make(map[byte]float64)
	for x := range support {
		counts[x] = prior
	}
	return &counter{counts, float64(len(counts)) * prior}
}

// Observe adds mass (given by count) for a particular byte event.
func (c *counter) Observe(b byte, count float64) {
	c.Counts[b] += count
	c.Total += count
}

// Sample generates a byte event from the normalized categorical distribution
// represented by the current counts.
func (c *counter) Sample() byte {
	sample := RandFloat() * c.Total
	for b, count := range c.Counts {
		if sample <= count {
			return b
		}
		sample -= count
	}

	return 0
}

// NameGen is a random generator based on an interpolated Markov process with a
// simplified Katz back-off scheme.
//
// The NameGen is trained by feeding name data through Observe. Depending on the
// parameterization and enough training data, calls to Generate will result in
// similar sounding names.
type NameGen struct {
	support map[byte]struct{}
	counts  map[string]*counter
	order   int
	prior   float64
}

// NewNameGen creates an empty NameGen with the given parameters.
//
// The order determines how many characters to consider at any given place in a
// name. Higher orders will result in more consistent results, as we can better
// model character and syllabic dependencies, but we will also require
// exponentially more data to properly learn the patterns. The prior determines
// how much we trust our data. Typically it is fairly low (less than 1). Higher
// values allow the generator to deviate more from the learned name model.
func NewNameGen(order int, prior float64) *NameGen {
	support := map[byte]struct{}{boundary[0]: {}}
	counts := make(map[string]*counter)
	return &NameGen{support, counts, order, prior}
}

// Observe updates the model based on the observed names.
func (g *NameGen) Observe(names ...string) {
	for _, name := range names {
		g.observe(name)
	}
}

// observe updates the model based on the observed name.
func (g *NameGen) observe(name string) {
	seq := strings.Repeat(boundary, g.order) + name + boundary
	for i := g.order; i < len(seq); i++ {
		context := seq[i-g.order : i]
		event := seq[i]

		if _, ok := g.support[event]; !ok {
			g.augmentSupport(event)
		}

		for j := 0; j <= len(context); j++ {
			g.getCounter(context[j:]).Observe(event, 1)
		}
	}
}

// augmentSupport adds an event to the support and the support of the
// underlying categorical distributions.
func (g *NameGen) augmentSupport(b byte) {
	g.support[b] = struct{}{}
	for _, counter := range g.counts {
		counter.Observe(b, g.prior)
	}
}

// getCounter gets the categorical distribution given a context.
func (g *NameGen) getCounter(context string) *counter {
	if _, ok := g.counts[context]; !ok {
		g.counts[context] = newCounter(g.support, g.prior)
	}
	return g.counts[context]
}

// Generate creates a random name using the observed name data.
func (g *NameGen) Generate() string {
	seq := string(g.sample(""))
	for seq[len(seq)-1] != boundary[0] {
		seq += string(g.sample(seq))
	}
	return seq[:len(seq)-1]
}

// backoff computes Katz back-off context with a threshold of 0 and weight of 1.
func (g *NameGen) backoff(context string) string {
	if len(context) > g.order {
		context = context[len(context)-g.order:]
	} else if len(context) < g.order {
		context = strings.Repeat(boundary, g.order-len(context)) + context
	}

	for len(context) > 0 {
		if _, ok := g.counts[context]; ok {
			return context
		}
		context = context[1:]
	}

	return context
}

// sample generates a byte event from the given context after back-off.
func (g *NameGen) sample(context string) byte {
	return g.getCounter(g.backoff(context)).Sample()
}

// TODO Add dialog cache
