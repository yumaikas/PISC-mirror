package pisc

import "fmt"

var ModSymbolCore = Module{
	Author:    "Andrew Owen",
	Name:      "",
	License:   "MIT",
	DocString: "SymbolCore",
	Load:      loadSymbolCore,
}

func _genSymbol(m *Machine) error {
	m.genSymbol()
	return nil
}

func _symbolNotEqual(m *Machine) error {
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
}

func loadSymbolCore(m *Machine) error {
	// Push a symbol onto the stack
	m.AddGoWordWithStack("<symbol>", "( -- symbol ) ", "Creates a unique symbol, and puts it on the stack", _genSymbol)
	// ( symbol symbol -- bool )
	m.AddGoWordWithStack("<symbol>", "( a b -- equal? ) ", "Compares two symbols, erroring it they aren't symbols", _symbolNotEqual)
	return m.ImportPISCAsset("stdlib/symbols.pisc")
}
