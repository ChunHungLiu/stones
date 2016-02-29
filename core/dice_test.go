package core

import (
	"math"
	"testing"
)

// LowerIncompleteGamma approximates the following integral:
// \sum_{k=0}^\inf \frac{(-1)^k z^{s+k}}{k!(s+k)}
func LowerIncompleteGamma(s, z float64) float64 {
	integral := 0.0
	for k := 0.0; k < 500; k++ {
		inc := (math.Pow(-1, k) * math.Pow(z, s+k)) / (math.Gamma(k+1) * (s + k))
		integral += inc
		if math.Abs(inc) < 1e-20 {
			break
		}
	}
	return integral
}

// ChiSquaredCDF computes the cumulative distribution function for the
// Chi-squared distribution.
func ChiSquaredCDF(x float64, k int) float64 {
	return LowerIncompleteGamma(float64(k)/2, x/2) / math.Gamma(float64(k)/2)
}

// ComputePearonStats computes the test statistics needed to run Pearon's
// Chi-squared goodness of fit test.
func ComputePearsonStats(obs, exp []int) (X2 float64, df int) {
	// Both obs and exp need to have the same number of categories, and same
	// total count.
	if len(obs) != len(exp) {
		panic(Error("observed and expected dimension mismatch"))
	}
	diff := 0
	for i := 0; i < len(obs); i++ {
		diff += obs[i] - exp[i]
	}
	if diff != 0 {
		panic(Error("observed and expected total mismatch"))
	}

	// Calculate Chi^2 as \sum_{i=1}^n \frac{(O_i - E_i)^2}{E_i}
	X2 = 0
	for i := 0; i < len(obs); i++ {
		X2 += math.Pow(float64(obs[i]-exp[i]), 2) / float64(exp[i])
	}

	// Degrees of freedom is 1 less than the number of categories since final
	// count is constrained by the total.
	df = len(obs) - 1

	return X2, df
}

func PearsonGoodness(obs, exp []int, alpha float64) (accept bool) {
	X2, df := ComputePearsonStats(obs, exp)
	pValue := ChiSquaredCDF(X2, df)
	return pValue <= alpha
}

// TestPearsonsChiSquaredTest is a test of the other dice tests!
func TestPearsonsChiSquaredTest(t *testing.T) {
	cases := []struct {
		Obs, Exp []int
		X2       float64
		DF       int
		PValue   float64
	}{
		{
			[]int{44, 56}, []int{50, 50},
			1.44, 1, .7699,
		}, {
			[]int{5, 8, 9, 8, 10, 20}, []int{10, 10, 10, 10, 10, 10},
			13.4, 5, 0.9801,
		}, {
			[]int{5, 8, 9, 8, 10, 20}, []int{10, 10, 10, 10, 10, 10},
			13.4, 5, 0.9801,
		},
	}
	for _, c := range cases {
		X2, df := ComputePearsonStats(c.Obs, c.Exp)
		if X2 != c.X2 || df != c.DF {
			t.Errorf("ComputePearsonStats(%v, %v) = (%f, %d) != (%f, %d)", c.Obs, c.Exp, X2, df, c.X2, c.DF)
			t.FailNow()
		}
		if pValue := ChiSquaredCDF(X2, df); math.Abs(pValue-c.PValue) > 1e-4 {
			t.Errorf("ChiSquaredCDF(%f, %d) = %f != %f", X2, df, pValue, c.PValue)
		}
	}
}

// We will be using Pearson's to do a sanity check on our dice funtions. We
// could of course run a more complete statistical test of randomness (such as
// Diehard or similar), but really we just want to make sure that we aren't
// causing some pathological errors. Consequently, we'll be content to just
// run Pearson's Chi-Squared test on some of the functions we added to the rand
// package. Note that these tests are all run using the global Dice object,
// which of course tests the underlying Dice methods.

// We'll reuse a common set of random seeds in order to have deterministic tests.
var seeds = []int64{
	0xFEE15600D,
	0xFEE15BAD,
	0xFA7CA7,
	0xBADD06,
	0x0B5E55ED,
	0xFA151F1AB1E,
}

func TestRandSeed(t *testing.T) {
	n := 1000
	exp := make(map[int64][]int64)
	for _, seed := range seeds {
		exp[seed] = make([]int64, n)
		RandSeed(seed)
		for i := 0; i < n; i++ {
			exp[seed][i] = RandInt63()
		}
	}

	for _, seed := range seeds {
		RandSeed(seed)
		for i := 0; i < n; i++ {
			if obs := RandInt63(); obs != exp[seed][i] {
				t.Logf("Seed(%#X) failed to produce consistent results", seed)
				t.FailNow()
			}
		}
	}
}

func TestRandBool(t *testing.T) {
	for _, seed := range seeds {
		RandSeed(seed)
		n := 1000
		obs := []int{0, 0}
		for i := 0; i < n; i++ {
			if RandBool() {
				obs[0]++
			} else {
				obs[1]++
			}
		}
		exp := []int{n / 2, n / 2}
		if !PearsonGoodness(obs, exp, .99) {
			t.Errorf("RandBool() failed Chi-squared test (seed=%#X)", seed)
		}
	}
}

func TestRandChance(t *testing.T) {
	cases := []float64{.01, .25, .5, .75, .99}
	for _, chance := range cases {
		for _, seed := range seeds {
			RandSeed(seed)
			n := 1000
			obs := []int{0, 0}
			for i := 0; i < n; i++ {
				if RandChance(chance) {
					obs[0]++
				} else {
					obs[1]++
				}
			}
			expTrue := int(float64(n) * chance)
			exp := []int{expTrue, n - expTrue}
			if !PearsonGoodness(obs, exp, .99) {
				t.Errorf("RandChance(%f) failed Chi-squared test (%v !~= %v, seed=%#X)", chance, obs, exp, seed)
			}
		}
	}
}

func TestRandChance_bounds(t *testing.T) {
	for _, seed := range seeds {
		RandSeed(seed)
		for i := 0; i < 100; i++ {
			if RandChance(0) {
				t.Log("RandChance(0) returned true (seed=%#X)", seed)
				t.FailNow()
			}
			if !RandChance(1) {
				t.Log("RandChance(1) returned false (seed=%#X", seed)
				t.FailNow()
			}
		}
	}
}

func TestRandRange(t *testing.T) {
	cases := []struct {
		Min, Max int
	}{
		{-1, 1},
		{0, 10},
		{20, 28},
	}
	for _, c := range cases {
		for _, seed := range seeds {
			RandSeed(seed)
			spread := c.Max - c.Min + 1
			n := 1000

			obs := make([]int, spread)
			for i := 0; i < n; i++ {
				obs[RandRange(c.Min, c.Max)-c.Min]++
			}

			exp := make([]int, spread)
			expTotal := 0
			for i := 0; i < spread-1; i++ {
				exp[i] = n / spread
				expTotal += exp[i]
			}
			exp[spread-1] = n - expTotal

			if !PearsonGoodness(obs, exp, .99) {
				t.Errorf("RandRange(%d, %d) failed Chi-squared test (%v !~= %v, seed=%#X)", c.Min, c.Max, obs, exp, seed)
			}
		}
	}
}

func TestRollXdY(t *testing.T) {
	cases := []struct {
		X, Y int
		Dist []float64
	}{
		{
			3, 6, []float64{
				0.462962962963,
				1.38888888889,
				2.77777777778,
				4.62962962963,
				6.94444444444,
				9.72222222222,
				11.5740740741,
				12.5,
				12.5,
				11.5740740741,
				9.72222222222,
				6.94444444444,
				4.62962962963,
				2.77777777778,
				1.38888888889,
				0.462962962963,
			},
		}, {
			2, 4, []float64{
				6.25,
				12.5,
				18.75,
				25,
				18.75,
				12.5,
				6.25,
			},
		},
	}
	for _, c := range cases {
		for _, seed := range seeds {
			RandSeed(seed)
			n := 1000

			total := 0
			exp := make([]int, len(c.Dist))
			for i, prob := range c.Dist {
				exp[i] = int(prob * float64(n))
				total += exp[i]
			}

			obs := make([]int, len(c.Dist))
			for i := 0; i < total; i++ {
				obs[RollXdY(c.X, c.Y)-c.X]++
			}

			if !PearsonGoodness(obs, exp, .99) {
				t.Errorf("RollXdY(%d, %d) failed Chi-squared test (%v !~= %v, seed=%#X)", c.X, c.Y, obs, exp, seed)
			}
		}
	}
}
