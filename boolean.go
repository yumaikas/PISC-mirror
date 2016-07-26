package main

func (m *machine) loadBooleanWords() error {
	m.predefinedWords["and"] = NilWord(func(m *machine) {
		a := m.popValue().(Boolean)
		b := m.popValue().(Boolean)
		m.pushValue(Boolean(a && b))
	})
	m.predefinedWords["or"] = NilWord(func(m *machine) {
		a := m.popValue().(Boolean)
		b := m.popValue().(Boolean)
		m.pushValue(Boolean(a || b))
	})
	m.predefinedWords["not"] = NilWord(func(m *machine) {
		a := m.popValue().(Boolean)
		m.pushValue(Boolean(!a))
	})
	return nil
}
