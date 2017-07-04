package pisc

// "strings"

var ModDictionaryCore = Module{
	Author:  "Andrew Owen",
	Name:    "DictionaryCore",
	License: "MIT", // TODO: Clarify here
	Load:    loadDictMod,
	// Possible: indicate PISC files to be loaded?
}

func loadDictMod(m *Machine) error {
	return m.loadDictWords()
}

func (m *Machine) loadDictWords() error {

	// Push a dictionary to the stack.
	m.PredefinedWords["<dict>"] = NilWord(func(m *Machine) {
		dict := make(map[string]StackEntry)
		m.PushValue(Dict(dict))
	})

	// ( dict key -- bool )
	m.PredefinedWords["dict-has-key?"] = NilWord(func(m *Machine) {
		key := m.PopValue().(String).String()
		dict := m.PopValue().(Dict)
		_, ok := dict[key]
		m.PushValue(Boolean(ok))
	})

	// ( dict value key -- dict )
	m.PredefinedWords["dict-set"] = NilWord(func(m *Machine) {
		key := m.PopValue().(String).String()
		value := m.PopValue()
		// Peek, since we have no intention of popping here.
		dict := m.PopValue().(Dict)
		dict[string(key)] = value
	})

	// ( dict key -- value )
	m.PredefinedWords["dict-get"] = NilWord(func(m *Machine) {
		key := m.PopValue().(String).String()
		m.PushValue(m.PopValue().(Dict)[key])
	})

	m.PredefinedWords["push-dict-keys"] = NilWord(func(m *Machine) {
		dic := m.PopValue().(Dict)
		for k, _ := range dic {
			m.PushValue(String(k))
		}
	})

	// ( dict -- key value )
	m.PredefinedWords["dict-get-rand"] = GoWord(func(m *Machine) error {
		dic := m.PopValue().(Dict)
		// Rely on random key ordering
		for k, v := range dic {
			m.PushValue(String(k))
			m.PushValue(v)
			break
		}
		return nil
	})
	return m.importPISCAsset("stdlib/dicts.pisc")
}
