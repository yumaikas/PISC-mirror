package main

import (
	"math/rand"
	"time"
)

func (m *machine) loadRandyWords() {

	m.helpDocs[word("seed-rand-time")] = "( -- ) Seeds the PSRNG with the current time"
	m.predefinedWords["seed-rand-time"] = NilWord(func(m *machine) {
		rand.Seed(time.Now().UTC().UnixNano())
	})

	m.predefinedWords["rand-int"] = NilWord(func(m *machine) {
		m.pushValue(Integer(rand.Int()))
	})

	m.helpDocs[word("range-rand")] = "( min max -- value ) Take a min and max, give a random integer"
	m.predefinedWords["range-rand"] = NilWord(func(m *machine) {
		max := m.popValue().(Integer)
		min := m.popValue().(Integer)
		m.pushValue(min + Integer(rand.Intn(int(max-min))))
	})
}
