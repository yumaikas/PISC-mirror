package main

import (
	"math/rand"
	"time"
)

var ModRandomCore = PISCModule{
	Author:    "Andrew Owen",
	Name:      "RandomCore",
	License:   "MIT",
	DocString: "Functions for random choice",
	Load:      loadRandy,
}

func loadRandy(m *machine) error {

	m.helpDocs["seed-rand-time"] = "( -- ) Seeds the PSRNG with the current time"
	m.predefinedWords["seed-rand-time"] = NilWord(func(m *machine) {
		rand.Seed(time.Now().UTC().UnixNano())
	})

	m.predefinedWords["rand-int"] = NilWord(func(m *machine) {
		m.pushValue(Integer(rand.Int()))
	})

	m.helpDocs["range-rand"] = "( min max -- value ) Take a min and max, give a random integer"
	m.predefinedWords["range-rand"] = NilWord(func(m *machine) {
		max := m.popValue().(Integer)
		min := m.popValue().(Integer)
		m.pushValue(min + Integer(rand.Intn(int(max-min))))
	})
	return m.importPISCAsset("stdlib/random.pisc")
}
