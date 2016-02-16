package core

import (
	"strconv"
	"testing"
	"unicode"
)

type FieldCase struct {
	Name      string
	Grid      StrGrid
	Radius    int
	ConvertFn func(StrGrid) ([]*Tile, map[*Tile]int)
	FieldFn   func(int, ...*Tile) Field
}

func (c FieldCase) Run(t *testing.T, casenum int) {
	goals, weights := c.ConvertFn(c.Grid)
	actual := c.FieldFn(c.Radius, goals...)
	for tile, weight := range weights {
		off := actual.Follow(tile)
		adj := tile.Adjacent[off]
		if adjweight := weights[adj]; adjweight > weight || (adjweight == weight && weight != 0) {
			t.Errorf("%s failed case %d (stepped from %d to %d)", c.Name, casenum, adjweight, weight)
		}
	}
}

func AttractiveFieldCase(g StrGrid) (goals []*Tile, weights map[*Tile]int) {
	goals, weights = make([]*Tile, 0), make(map[*Tile]int)
	callback := func(t *Tile, c byte) {
		t.Face = Glyph{rune(c), ColorWhite}
		switch c {
		case '#':
			t.Pass = false
		case '@':
			goals = append(goals, t)
			weights[t] = 0
		default:
			if weight, err := strconv.Atoi(string(c)); err == nil {
				weights[t] = weight
			}
		}
	}
	g.Convert(callback)
	return goals, weights
}

func TestAttractiveField(t *testing.T) {
	cases := []struct {
		g StrGrid
		r int
	}{
		{
			StrGrid{
				"#######",
				"#@1234#",
				"#11234#",
				"#22234#",
				"#33334#",
				"#44444#",
				"#######",
			}, 10,
		}, {
			StrGrid{
				"########",
				"#987666#",
				"#####56#",
				"#544456#",
				"#543####",
				"#54321@#",
				"########",
			}, 10,
		}, {
			StrGrid{
				"########",
				"#000333#",
				"#####23#",
				"#21@123#",
				"#211####",
				"#222230#",
				"########",
			}, 3,
		}, {
			StrGrid{
				"#############",
				"#@1234555555#",
				"#11234444444#",
				"#22234333333#",
				"#33334322222#",
				"#44444321112#",
				"#55554321@12#",
				"#66654321112#",
				"#############",
			}, 10,
		}, {
			StrGrid{
				"########",
				"#@12300#",
				"#####00#",
				"#000000#",
				"#003####",
				"#00321@#",
				"########",
			}, 3,
		},
	}
	for i, c := range cases {
		FieldCase{"AttractiveField", c.g, c.r, AttractiveFieldCase, AttractiveField}.Run(t, i)
	}
}

func ReplusiveFieldCase(g StrGrid) (ungoals []*Tile, weights map[*Tile]int) {
	ungoals, weights = make([]*Tile, 0), make(map[*Tile]int)
	callback := func(t *Tile, c byte) {
		switch c {
		case '#':
			t.Pass = false
		default:
			if r := rune(c); unicode.IsLetter(r) {
				weights[t] = int(unicode.ToLower(r) - 'a')
				if !unicode.IsLower(r) {
					ungoals = append(ungoals, t)
				}
			} else if unicode.IsDigit(r) {
				weights[t] = int(r - '0')
			}
		}
	}
	g.Convert(callback)
	return ungoals, weights
}

func TestRepusliveField(t *testing.T) {
	cases := []struct {
		g StrGrid
		r int
	}{
		{
			StrGrid{
				"#######",
				"#Edcba#",
				"#ddcba#",
				"#cccba#",
				"#bbbba#",
				"#aaaaa#",
				"#######",
			}, 10,
		}, {
			StrGrid{
				"########",
				"#mlkjii#",
				"######H#",
				"######g#",
				"######f#",
				"#abcdef#",
				"########",
			}, 10,
		}, {
			StrGrid{
				"########",
				"#000aaa#",
				"#####ba#",
				"#edDcba#",
				"#edc####",
				"#edcba0#",
				"########",
			}, 3,
		}, {
			StrGrid{
				"#############",
				"#Ggggggggghi#",
				"#ffffffffghi#",
				"#eeeeeeefghi#",
				"#ddddddefghi#",
				"#cccccdefghi#",
				"#bbbbcdefGhi#",
				"#aaabcdefghi#",
				"#############",
			}, 10,
		},
	}
	for i, c := range cases {
		FieldCase{"ReplusiveField", c.g, c.r, ReplusiveFieldCase, ReplusiveField}.Run(t, i)
	}
}
