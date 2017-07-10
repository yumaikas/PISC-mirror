package pisc

var ModBoolCore = Module{
	Author:    "Andrew Owen",
	Name:      "BoolCore",
	License:   "MIT",
	DocString: "The 3 basic boolean operation, and/or and not",
	Load:      loadBoolCore,
}

func loadBoolCore(m *Machine) error {
	m.PredefinedWords["and"] = NilWord(func(m *Machine) {
		a := m.PopValue().(Boolean)
		b := m.PopValue().(Boolean)
		m.PushValue(Boolean(a && b))
	})
	m.PredefinedWords["or"] = NilWord(func(m *Machine) {
		a := m.PopValue().(Boolean)
		b := m.PopValue().(Boolean)
		m.PushValue(Boolean(a || b))
	})
	m.PredefinedWords["not"] = NilWord(func(m *Machine) {
		a := m.PopValue().(Boolean)
		m.PushValue(Boolean(!a))
	})
	return m.ImportPISCAsset("stdlib/bools.pisc")
}
