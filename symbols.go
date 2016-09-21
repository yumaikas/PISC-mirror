package main

import "fmt"

func (m *machine) loadSymbolWords() error {

	// Push a symbol onto the stack
	// ( -- #symbol )
	m.predefinedWords["<symbol>"] = NilWord(func(m *machine) {
		m.genSymbol()
	})

	// ( symbol symbol -- bool )
	m.predefinedWords["symb-neq"] = GoWord(func(m *machine) error {
		a, ok := m.popValue().(Symbol)
		b, bOk := m.popValue().(Symbol)
		if ok && bOk {
			m.pushValue(Boolean(a != b))

		} else if ok || bOk {
			// If one of them is symbol, but they aren't equal,then return true
			m.pushValue(Boolean(true))
		} else {
			return fmt.Errorf("Symb-neq called on two non-symbols!")
		}
		return nil

	})

	return nil
}