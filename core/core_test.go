package core

import (
	"testing"
)

func TestMax(t *testing.T) {
	cases := []struct {
		x, y     int
		expected int
	}{
		{-4, 4, 4},
		{0, 4, 4},
		{4, 4, 4},
		{8, 4, 8},
		{4, 0, 4},
		{4, -4, 4},
	}
	for _, c := range cases {
		actual := Max(c.x, c.y)
		if c.expected != actual {
			t.Errorf("Max(%d, %d) = %d != %d", c.x, c.y, actual, c.expected)
		}
	}
}

func TestMin(t *testing.T) {
	cases := []struct {
		x, y     int
		expected int
	}{
		{-4, 4, -4},
		{0, 4, 0},
		{4, 4, 4},
		{8, 4, 4},
		{4, 0, 0},
		{4, -4, -4},
	}
	for _, c := range cases {
		actual := Min(c.x, c.y)
		if c.expected != actual {
			t.Errorf("Min(%d, %d) = %d != %d", c.x, c.y, actual, c.expected)
		}
	}
}

func TestMod(t *testing.T) {
	cases := []struct {
		x, y     int
		expected int
	}{
		{-10, 5, 0},
		{-9, 5, 1},
		{-8, 5, 2},
		{-7, 5, 3},
		{-6, 5, 4},
		{-5, 5, 0},
		{-4, 5, 1},
		{-3, 5, 2},
		{-2, 5, 3},
		{-1, 5, 4},
		{0, 5, 0},
		{1, 5, 1},
		{2, 5, 2},
		{3, 5, 3},
		{4, 5, 4},
		{5, 5, 0},
		{6, 5, 1},
		{7, 5, 2},
		{8, 5, 3},
		{9, 5, 4},
		{10, 5, 0},
	}
	for _, c := range cases {
		actual := Mod(c.x, c.y)
		if c.expected != actual {
			t.Errorf("Mod(%d, %d) = %d != %d", c.x, c.y, actual, c.expected)
		}
	}
}

func TestAbs(t *testing.T) {
	cases := []struct {
		x        int
		expected int
	}{
		{1, 1},
		{2, 2},
		{3, 3},
		{0, 0},
		{-1, 1},
		{-2, 2},
		{-3, 3},
	}
	for _, c := range cases {
		actual := Abs(c.x)
		if c.expected != actual {
			t.Errorf("Abs(%d) = %d != %d", c.x, actual, c.expected)
		}
	}
}

func TestClamp(t *testing.T) {
	cases := []struct {
		min, val, max int
		expected      int
	}{
		{-1, -2, 1, -1},
		{-1, -1, 1, -1},
		{-1, 0, 1, 0},
		{-1, 1, 1, 1},
		{-1, 2, 1, 1},

		{0, -1, 10, 0},
		{0, 0, 10, 0},
		{0, 1, 10, 1},
		{0, 5, 10, 5},
		{0, 9, 10, 9},
		{0, 10, 10, 10},
		{0, 11, 10, 10},
	}
	for _, c := range cases {
		actual := Clamp(c.min, c.val, c.max)
		if c.expected != actual {
			t.Errorf("Clamp(%d, %d, %d) = %d != %d", c.min, c.val, c.max, actual, c.expected)
		}
	}
}

func TestInBounds(t *testing.T) {
	cases := []struct {
		x, y, w, h int
		expected   bool
	}{
		{-1, 0, 80, 24, false},
		{0, -1, 80, 24, false},
		{-1, -1, 80, 24, false},
		{0, 0, 80, 24, true},
		{80, 23, 80, 24, false},
		{79, 24, 80, 24, false},
		{80, 24, 80, 24, false},
		{79, 23, 80, 24, true},
		{40, 12, 80, 24, true},
	}
	for _, c := range cases {
		actual := InBounds(c.x, c.y, c.w, c.h)
		if c.expected != actual {
			t.Errorf("InBounds(%d, %d, %d, %d) = %d != %d", c.x, c.y, c.w, c.h, actual, c.expected)
		}
	}
}

func TestRound(t *testing.T) {
	cases := []struct {
		x        float64
		n        int
		expected float64
	}{
		{0, 0, 0},
		{.1, 0, 0},
		{.49, 0, 0},
		{.5, 0, 1},
		{.51, 0, 1},
		{.9, 0, 1},
		{1, 0, 1},

		{.620, 2, .62},
		{.621, 2, .62},
		{.624, 2, .62},
		{.625, 2, .63},
		{.626, 2, .63},
		{.629, 2, .63},
		{.630, 2, .63},

		{-.1, 0, 0},
		{-.49, 0, 0},
		{-.5, 0, -1},
		{-.51, 0, -1},
		{-.9, 0, -1},
		{-1, 0, -1},

		{-.620, 2, -.62},
		{-.621, 2, -.62},
		{-.624, 2, -.62},
		{-.625, 2, -.63},
		{-.626, 2, -.63},
		{-.629, 2, -.63},
		{-.630, 2, -.63},
	}
	for _, c := range cases {
		actual := Round(c.x, c.n)
		if c.expected != actual {
			t.Errorf("Round(%f, %d) = %f != %f", c.x, c.n, actual, c.expected)
		}
	}
}
