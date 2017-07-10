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

func loadRandy(m *Machine) error {

	m.HelpDocs["seed-rand-time"] = "( -- ) Seeds the PSRNG with the current time"
	m.PredefinedWords["seed-rand-time"] = NilWord(func(m *Machine) {
		rand.Seed(time.Now().UTC().UnixNano())
	})

	m.PredefinedWords["rand-int"] = NilWord(func(m *Machine) {
		m.PushValue(Integer(rand.Int()))
	})

	m.HelpDocs["range-rand"] = "( min max -- value ) Take a min and max, give a random integer"
	m.PredefinedWords["range-rand"] = NilWord(func(m *Machine) {
		max := m.PopValue().(Integer)
		min := m.PopValue().(Integer)
		m.PushValue(min + Integer(rand.Intn(int(max-min))))
	})
	return m.ImportPISCAsset("stdlib/random.pisc")
}
