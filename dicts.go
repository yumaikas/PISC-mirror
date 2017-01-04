package main

// "strings"

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
	m.predefinedWords["dict-keys"] = NilWord(func(m *machine) {
		dic := m.popValue().(Dict)
		// Rely on random key ordering
		for k, _ := range dic {
			m.pushValue(String(k))
			break
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
