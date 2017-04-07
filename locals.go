package main

import (
	"fmt"
	"strings"
)

var ModLocalsCore = PISCModule{
	Author:    "Andrew Owen",
	Name:      "LocalsCore",
	License:   "MIT",
	DocString: "TODO: Fill this in",
	Load:      loadLocalCore,
}

// Used when parsing words
func isLocalWordPrefix(w word) bool {
	return strings.HasPrefix(w.str, ":") ||
		strings.HasPrefix(w.str, "$")
}

var ErrLocalNotFound = fmt.Errorf("Local variable not found!")
var ErrNoLocalsExist = fmt.Errorf("A local frame hasn't been allocated with get-locals!")
var ErrAttemtToIncrementNonNumber = fmt.Errorf("Attempted to increment a non-integer")

func loadLocalCore(m *machine) error {
	// Make sure we always have locals
	m.locals = append(m.locals, make(map[string]stackEntry))
	// ( val name -- )
	m.predefinedWords["set-local"] = GoWord(func(m *machine) error {
		varName := m.popValue().(String)
		if len(m.locals) <= 0 {
			return ErrNoLocalsExist
		}
		m.locals[len(m.locals)-1][varName.String()] = m.popValue()
		return nil
	})
	// ( name -- val )
	m.predefinedWords["get-local"] = GoWord(func(m *machine) error {
		varName := m.popValue().(String)
		if len(m.locals) <= 0 {
			return ErrNoLocalsExist
		}
		if val, ok := m.locals[len(m.locals)-1][varName.String()]; ok {
			m.pushValue(val)
			return nil
		} else {
			return ErrLocalNotFound
		}
	})
	m.predefinedWords["get-locals"] = NilWord(func(m *machine) {
		m.locals = append(m.locals, make(map[string]stackEntry))
	})
	m.predefinedWords["drop-locals"] = NilWord(func(m *machine) {
		m.locals = m.locals[:len(m.locals)-1]
	})
	// ( -- locals.. )
	// Run a quotation for each local
	m.predefinedWords["each-local"] = GoWord(func(m *machine) error {
		quot := m.popValue().(*quotation)
		code := quot.toCode()
		for key, val := range m.locals[len(m.locals)-1] {
			m.pushValue(val)
			m.pushValue(String(key))
			code.idx = 0
			err := m.execute(code)
			if err != nil {
				return err
			}
		}
		return nil
	})

	// TODO: compress this
	// incr-local-var is here to help ++ be faster
	m.addGoWord("incr-local-var", "( name -- )", GoWord(func(m *machine) error {
		varName := m.popValue().(String)
		if len(m.locals) <= 0 {
			return ErrNoLocalsExist
		}
		if val, ok := m.locals[len(m.locals)-1][varName.String()]; ok {
			// TODO: Clean up this cast
			v, canNumber := val.(Integer)
			if canNumber {
				m.pushValue(v + 1)
				return nil
			} else {
				return ErrAttemtToIncrementNonNumber
			}
		} else {
			return ErrLocalNotFound
		}
	}))
	m.addGoWord("decr-local-var", "( name -- )", GoWord(func(m *machine) error {
		varName := m.popValue().(String)
		if len(m.locals) <= 0 {
			return ErrNoLocalsExist
		}
		if val, ok := m.locals[len(m.locals)-1][varName.String()]; ok {
			// TODO: Clean up this cast
			v, canNumber := val.(Integer)
			if canNumber {
				m.pushValue(v - 1)
				return nil
			} else {
				return ErrAttemtToIncrementNonNumber
			}
		} else {
			return ErrLocalNotFound
		}
	}))
	return m.importPISCAsset("stdlib/locals.pisc")
}
