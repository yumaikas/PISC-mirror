package main

// "strings"

var ModDictionaryCore = PISCModule{
	Author:  "Andrew Owen",
	Name:    "DictionaryCore",
	License: "MIT", // TODO: Clarify here
	Load:    loadDictMod,
	// Possible: indicate PISC files to be loaded?
}

func loadDictMod(m *machine) error {
	return m.loadDictWords()
}

func (m *machine) loadDictWords() error {

	// Push a dictionary to the stack.
	m.predefinedWords["<dict>"] = NilWord(func(m *machine) {
		dict := make(map[string]stackEntry)
		m.pushValue(Dict(dict))
	})

	// ( dict key -- bool )
	m.predefinedWords["dict-has-key?"] = NilWord(func(m *machine) {
		key := m.popValue().(String).String()
		dict := m.popValue().(Dict)
		_, ok := dict[key]
		m.pushValue(Boolean(ok))
	})

	// ( dict value key -- dict )
	m.predefinedWords["dict-set"] = NilWord(func(m *machine) {
		key := m.popValue().(String).String()
		value := m.popValue()
		// Peek, since we have no intention of popping here.
		dict := m.popValue().(Dict)
		dict[string(key)] = value
	})

	// ( dict key -- value )
	m.predefinedWords["dict-get"] = NilWord(func(m *machine) {
		key := m.popValue().(String).String()
		m.pushValue(m.popValue().(Dict)[key])
	})

	m.predefinedWords["push-dict-keys"] = NilWord(func(m *machine) {
		dic := m.popValue().(Dict)
		for k, _ := range dic {
			m.pushValue(String(k))
		}
	})

	// ( dict -- key value )
	m.predefinedWords["dict-get-rand"] = GoWord(func(m *machine) error {
		dic := m.popValue().(Dict)
		// Rely on random key ordering
		for k, v := range dic {
			m.pushValue(String(k))
			m.pushValue(v)
			break
		}
		return nil
	})
	return nil
}
