package core

import (
	"fmt"
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
	fmt.Println(X2)
	for i := 0; i < len(obs); i++ {
		fmt.Println(obs[i] - exp[i])
		X2 += math.Pow(float64(obs[i]-exp[i]), 2) / float64(exp[i])
		fmt.Println(X2)
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
// run Pearson's Chi-Squared test on each of the functions we added to the rand
// package. Note that these tests are all run using the global Dice object,
// which of course tests the underlying Dice methods.

// We'll reuse a common set of random seeds in order to have deterministic tests.
var seeds = []int64{
	0xDEADBEAF,
	0xBA5EBA11,
	0xF00BA55,
	0xDEAD10CC,
	0xFACEB00C,
	0x5EED,
}

func TestRandBool(t *testing.T) {
	for _, seed := range seeds {
		RandSeed(seed)
		n := 10
		exp := []int{n / 2, n / 2}
		obs := []int{0, 0}
		for i := 0; i < n; i++ {
			if RandBool() {
				obs[0]++
			} else {
				obs[1]++
			}
		}
		if !PearsonGoodness(obs, exp, .95) {
			t.Errorf("RandBool failed Chi-squared test (seed=%#X)", seed)
		}
	}
}
