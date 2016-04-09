package main

import (
	"strings"
)

// TODO: add words for printing and such here.
func isDictWord(w word) bool {
	return w == ">>" ||
		w == "<dict>" ||
		strings.HasSuffix(string(w), "#get") ||
		strings.HasSuffix(string(w), "#set")

}

func (m *machine) executeDictWord(w word) error {
	switch {
	// Push a dictionary to the stack.
	case w == "<dict>":
		dict := make(map[string]stackEntry)
		m.pushValue(Dict(dict))
	// The is a multi-read operation
	case strings.HasSuffix(string(w), ">>"):
		q := m.popValue().(quotation)
		// Keep from popping the dictionary
		dict := m.values[len(m.values)-1].(Dict)
		for _, innerWord := range q {
			m.pushValue(dict[string(innerWord)])
		}
	// This is a set operation
	case strings.HasSuffix(string(w), "#set"):
		key := w[:len(w)-len("#set")]
		value := m.popValue()
		// Peek, since we have no intention of popping here.
		dict := m.values[len(m.values)-1].(Dict)
		dict[string(key)] = value

		// This is a get operation
	case strings.HasSuffix(string(w), "#get"):
		key := string(w[:len(w)-len("#get")])
		m.pushValue(m.values[len(m.values)-1].(Dict)[key])
	}
	return nil
}
