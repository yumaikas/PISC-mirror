package pisc

import "fmt"

var ModSymbolCore = Module{
	Author:    "Andrew Owen",
	Name:      "",
	License:   "MIT",
	DocString: "SymbolCore",
	Load:      loadSymbolCore,
}

func loadSymbolCore(m *Machine) error {

	// Push a symbol onto the stack
	// ( -- #symbol )
	m.PredefinedWords["<symbol>"] = NilWord(func(m *Machine) {
		m.genSymbol()
	})

	// ( symbol symbol -- bool )
	m.PredefinedWords["symb-neq"] = GoWord(func(m *Machine) error {
		a, ok := m.PopValue().(Symbol)
		b, bOk := m.PopValue().(Symbol)
		if ok && bOk {
			m.PushValue(Boolean(a != b))

		} else if ok || bOk {
			// If one of them is symbol, but they aren't equal,then return true
			m.PushValue(Boolean(true))
		} else {
			return fmt.Errorf("Symb-neq called on two non-symbols!")
		}
		return nil

	})
	return m.importPISCAsset("stdlib/symbols.pisc")
}
