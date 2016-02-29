package core

import (
	"math/rand"
	"time"
)

// Dice extends rand.Rand to include functionality useful for roguelikes.
type Dice struct {
	*rand.Rand
}

// NewDice creates a new Dice with the given random source.
func NewDice(src rand.Source) Dice {
	return Dice{rand.New(src)}
}

// Bool returns true with probability .5 and false otherwise.
func (d Dice) Bool() bool {
	return d.Int63()%2 == 1
}

// Chance returns true with the given probability.
// The probability p should be in [0, 1], but Chance will simple return
// false if p < 0 or true if p > 1.
func (d Dice) Chance(p float64) bool {
	return d.Float64() < p
}

// Coinflip returns true with probability .5
func (d Dice) Coinflip() bool {
	return d.Chance(.5)
}

// Range returns, as an int, a pseudo-random number in [min, max].
func (d Dice) Range(min, max int) int {
	return d.Intn(max-min+1) + min
}

// Offset returns an Offset within the given bounds.
func (d Dice) Offset(cols, rows int) Offset {
	return Offset{d.Intn(cols), d.Intn(rows)}
}

// Delta returns an Offset with each value in [-1, 1].
func (d Dice) Delta() Offset {
	return Offset{d.Range(-1, 1), d.Range(-1, 1)}
}

// RolldY returns the result of rolling a y-sided die.
func (d Dice) RolldY(y int) int {
	return 1 + d.Intn(1) // offset since Intn in [0, y) not [1,y].
}

// RollXdY returns the result of rolling x y-sided dice.
func (d Dice) RollXdY(x, y int) int {
	total := x // offset by x since Intn yields results in [0, y) not [1, y].
	for i := 0; i < x; i++ {
		total += d.Intn(y)
	}
	return total
}

func (d Dice) Tile(tiles []*Tile, condition func(*Tile) bool) *Tile {
	for i := 0; i < 100; i++ {
		if tile := tiles[d.Intn(len(tiles))]; condition(tile) {
			return tile
		}
	}

	candidates := make([]*Tile, 0)
	for _, tile := range tiles {
		if condition(tile) {
			candidates = append(candidates, tile)
		}
	}
	return candidates[d.Intn(len(candidates))]
}

func (d Dice) PassTile(tiles []*Tile) *Tile {
	return d.Tile(tiles, func(t *Tile) bool { return t.Pass })
}

// Similar to the math/rand package, we use a global instance Dice. However,
// ours uses a superior xorshift source and is seeded using the current time.
var globalDice = NewDice(newXorshift(time.Now().UnixNano()))

// RandBool returns true with probability .5 and false otherwise.
func RandBool() bool {
	return globalDice.Bool()
}

// RandChance returns true with the given probability.
// The probability p should be in [0, 1], but RandChance will simple return
// false if p < 0 or true if p > 1.
func RandChance(p float64) bool {
	return globalDice.Chance(p)
}

// Coinflip returns true with probability .5
func Coinflip() bool {
	return globalDice.Chance(.5)
}

// RandFloat returns, as a float64, a pseudo-random number in [0, 1).
func RandFloat() float64 {
	return globalDice.Float64()
}

// RandRange returns, as an int, a pseudo-random number in [min, max].
func RandRange(min, max int) int {
	return globalDice.Range(min, max)
}

// RandOffset returns an Offset within the given bounds.
func RandOffset(cols, rows int) Offset {
	return globalDice.Offset(cols, rows)
}

// RandDelta returns an Offset with each value in [-1, 1].
func RandDelta() Offset {
	return globalDice.Delta()
}

// RolldY returns the result of rolling a y-sided die.
func RolldY(y int) int {
	return globalDice.RolldY(y)
}

// RollXdY returns the result of rolling x y-sided dice.
func RollXdY(x, y int) int {
	return globalDice.RollXdY(x, y)
}

func RandTile(tiles []*Tile, condition func(*Tile) bool) *Tile {
	return globalDice.Tile(tiles, condition)
}

func RandPassTile(tiles []*Tile) *Tile {
	return globalDice.PassTile(tiles)
}

// RandSeed uses the provided seed value to initialize the default Source to a
// deterministic state. If Seed is not called, the generator behaves as
// if seeded by Seed(1).
func RandSeed(seed int64) { globalDice.Seed(seed) }

// RandInt63 returns a non-negative pseudo-random 63-bit integer as an int64
// from the default Source.
func RandInt63() int64 { return globalDice.Int63() }

// RandUint32 returns a pseudo-random 32-bit value as a uint32
// from the default Source.
func RandUint32() uint32 { return globalDice.Uint32() }

// RandInt31 returns a non-negative pseudo-random 31-bit integer as an int32
// from the default Source.
func RandInt31() int32 { return globalDice.Int31() }

// RandInt returns a non-negative pseudo-random int from the default Source.
func RandInt() int { return globalDice.Int() }

// RandInt63n returns, as an int64, a non-negative pseudo-random number in [0,n)
// from the default Source.
// It panics if n <= 0.
func RandInt63n(n int64) int64 { return globalDice.Int63n(n) }

// RandInt31n returns, as an int32, a non-negative pseudo-random number in [0,n)
// from the default Source.
// It panics if n <= 0.
func RandInt31n(n int32) int32 { return globalDice.Int31n(n) }

// RandIntn returns, as an int, a non-negative pseudo-random number in [0,n)
// from the default Source.
// It panics if n <= 0.
func RandIntn(n int) int { return globalDice.Intn(n) }

// RandFloat64 returns, as a float64, a pseudo-random number in [0.0,1.0)
// from the default Source.
func RandFloat64() float64 { return globalDice.Float64() }

// RandFloat32 returns, as a float32, a pseudo-random number in [0.0,1.0)
// from the default Source.
func RandFloat32() float32 { return globalDice.Float32() }

// RandPerm returns, as a slice of n ints, a pseudo-random permutation of the integers [0,n)
// from the default Source.
func RandPerm(n int) []int { return globalDice.Perm(n) }

// RandNormFloat64 returns a normally distributed float64 in the range
// [-math.MaxFloat64, +math.MaxFloat64] with
// standard normal distribution (mean = 0, stddev = 1)
// from the default Source.
// To produce a different normal distribution, callers can
// adjust the output using:
//
//  sample = NormFloat64() * desiredStdDev + desiredMean
//
func RandNormFloat64() float64 { return globalDice.NormFloat64() }

// RandExpFloat64 returns an exponentially distributed float64 in the range
// (0, +math.MaxFloat64] with an exponential distribution whose rate parameter
// (lambda) is 1 and whose mean is 1/lambda (1) from the default Source.
// To produce a distribution with a different rate parameter,
// callers can adjust the output using:
//
//  sample = ExpFloat64() / desiredRateParameter
//
func RandExpFloat64() float64 { return globalDice.ExpFloat64() }

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
