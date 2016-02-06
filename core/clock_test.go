package core

import (
	"testing"
)

func checkAdvance(expected []Entity, actual map[Entity]struct{}) bool {
	if len(expected) != len(actual) {
		return false
	}
	for _, e := range expected {
		if _, ok := actual[e]; !ok {
			return false
		}
	}
	return true
}

func checkSchedule(t *testing.T, c *DeltaClock, schedule [][]Entity, speeds map[Entity]float64) {
	for i, expected := range schedule {
		actual := c.Advance()
		if !checkAdvance(expected, actual) {
			t.Errorf("Unexpected delta (%v!=%v) at delta=%d", expected, actual, i)
		}
		for e := range actual {
			c.Schedule(e, speeds[e])
		}
	}
}

func initClock(speeds map[Entity]float64) *DeltaClock {
	c := NewDeltaClock()
	for e, s := range speeds {
		c.Schedule(e, s)
	}
	return c
}

func TestDeltaClock_Schedule(t *testing.T) {
	e1, e2, e3 := &ComponentSlice{}, &ComponentSlice{}, &ComponentSlice{}
	speeds := map[Entity]float64{e1: 1, e2: 1.5, e3: 2}
	c := initClock(speeds)

	schedule := [][]Entity{
		{e1}, {e2},
		{e1, e3}, {e2},
		{e1}, {e2},
		{e1, e3}, {e2},
		{e1}, {e2},
	}
	checkSchedule(t, c, schedule, speeds)
}

func TestDeltaClock_Unschedule(t *testing.T) {
	e1, e2, e3 := &ComponentSlice{}, &ComponentSlice{}, &ComponentSlice{}
	speeds := map[Entity]float64{e1: 2, e2: 2, e3: 3}
	c := initClock(speeds)
	c.Unschedule(e2)
	c.Unschedule(e3)
	schedule := [][]Entity{{e1}, {}, {e1}}
	checkSchedule(t, c, schedule, speeds)
}
