package main

var ModBoolCore = PISCModule{
	Author:    "Andrew Owen",
	Name:      "BoolCore",
	License:   "MIT",
	DocString: "The 3 basic boolean operation, and/or and not",
	Load:      loadBoolCore,
}

func loadBoolCore(m *machine) error {
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
	return m.importPISCAsset("stdlib/bools.pisc")
}
