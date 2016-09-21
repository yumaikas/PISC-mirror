package main

import (
// "strings"
)

func (m *machine) loadDictWords() error {

	// Push a dictionary to the stack.
	m.predefinedWords["<dict>"] = NilWord(func(m *machine) {
		dict := make(map[string]stackEntry)
		m.pushValue(Dict(dict))
	})

	// TODO: Consider deprecating this
	// ( dict quot -- vals.. )
	m.predefinedWords["dict-get-many"] = NilWord(func(m *machine) {
		q := m.popValue().(quotation)
		// Keep from popping the dictionary
		dict := m.values[len(m.values)-1].(Dict)
		for _, innerWord := range q.code {
			m.pushValue(dict[string(innerWord)])
		}
	})

	// ( dict key -- dict bool )
	m.predefinedWords["dict-has-key"] = NilWord(func(m *machine) {
		key := m.popValue().(String).String()
		dict := m.popValue().(Dict)
		_, ok := dict[key]
		m.pushValue(dict)
		m.pushValue(Boolean(ok))
	})

	// ( dict value key -- dict )
	m.predefinedWords["dict-set"] = NilWord(func(m *machine) {
		key := m.popValue().(String).String()
		value := m.popValue()
		// Peek, since we have no intention of popping here.
		dict := m.values[len(m.values)-1].(Dict)
		dict[string(key)] = value
	})

	// ( dict key -- dict value )
	m.predefinedWords["dict-get"] = NilWord(func(m *machine) {
		key := m.popValue().(String).String()
		m.pushValue(m.values[len(m.values)-1].(Dict)[key])
	})
	return nil
}
