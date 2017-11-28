package pisc

import (
	"math/rand"
	"time"
)

var ModRandomCore = Module{
	Author:    "Andrew Owen",
	Name:      "RandomCore",
	License:   "MIT",
	DocString: "Functions for random choice",
	Load:      loadRandy,
}

func _seedRandomTime(m *Machine) error {
	rand.Seed(time.Now().UTC().UnixNano())
	return nil
}

func _getRandInt(m *Machine) error {
	m.PushValue(Integer(rand.Int()))
	return nil
}

func _rangeRandomInt(m *Machine) error {
	max := m.PopValue().(Integer)
	min := m.PopValue().(Integer)
	m.PushValue(min + Integer(rand.Intn(int(max-min))))
	return nil
}

func loadRandy(m *Machine) error {

	m.AddGoWordWithStack("seed-rand-time", "( -- )", "Seeds the PSRNG with the current time", _seedRandomTime)
	m.AddGoWordWithStack("rand-int", "( -- random-int )", "Get a random positive int", _getRandInt)
	m.AddGoWordWithStack(
		"range-rand",
		"( min max -- random-int )",
		"Get a random positive int between min and max", _rangeRandomInt)

	return m.ImportPISCAsset("stdlib/random.pisc")
}
