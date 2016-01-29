package core

import (
	"math/rand"
	"time"
)

// Similar to the math/rand package, we use a global instance of Rand. However,
// ours uses a superior xorshift source and is seeded using the current time.
var globalRand = rand.New(newXorshift(time.Now().UnixNano()))

// RandBool returns true with probability .5 and false otherwise.
func RandBool() bool {
	return globalRand.Int63()%2 == 1
}

// RandChance returns true with the given probability.
// The probability p should be in [0, 1], but RandChance will simple return
// false if p < 0 or true if p > 1.
func RandChance(p float64) bool {
	return globalRand.Float64() < p
}

// RandFloat returns, as a float64, a pseudo-random number in [0, 1).
func RandFloat() float64 {
	return globalRand.Float64()
}

// RandInt returns, as an int, a non-negative pseudo-random number in [0,n).
func RandInt(n int) int {
	return globalRand.Intn(n)
}

// RandRange returns, as an int, a pseudo-random number in [min, max].
func RandRange(min, max int) int {
	return RandInt(max-min+1) + min
}

// RandCoord returns two ints in [0, cols) and [0, rows) respectively.
func RandCoord(cols, rows int) (x, y int) {
	return RandInt(cols), RandInt(rows)
}

// xorshift is a rand.Source implementing the xorshift1024* algorithm. It works
// by scrambling the output of an xorshift generator with a 64-bit invertible
// multiplier. While not cryptographically secure, this algorithm is faster
// and produces better output than the famous MT19937-64 algorithm and should
// be higher quality than the built in GFSR source found in the rand package.
type xorshift struct {
	state [16]uint64
	index int
}

// Seed initialize the state array based on a given int64 seed value.
func (x *xorshift) Seed(seed int64) {
	// Since we use a single 64 bit seed, we use an xorshift64* generator
	// to get the 1024 bits we need to seed the xorshift1024* generator.
	s := uint64(seed)
	for i := 0; i < 16; i++ {
		s ^= s >> 12
		s ^= s << 25
		s ^= s >> 27
		s *= 2685821657736338717
		x.state[i] = s
	}
}

// Int63 gets the next positive int64 from the sequence.
func (x *xorshift) Int63() int64 {
	a := x.state[x.index]
	a ^= a >> 30

	x.index = (x.index + 1) & 15
	b := x.state[x.index]
	b ^= b << 31
	b ^= b >> 11

	c := a ^ b
	x.state[x.index] = c
	c = c * 1783497276652981

	return int64(c >> 1)
}

// newXorshift returns a rand.Source which implements the xorshift1024*
// algorithm. While not cryptographically secure, this source should be
// superior in both speed and randomness to the default GFSR source found in
// the rand package.
func newXorshift(seed int64) rand.Source {
	var x xorshift
	x.Seed(seed)
	return &x
}

// TODO Completely wrap math/rand
