package pisc

import (
	"fmt"
	"strings"
)

var ModLocalsCore = Module{
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

func loadLocalCore(m *Machine) error {
	// Make sure we always have locals
	m.Locals = append(m.Locals, make(map[string]StackEntry))
	// ( val name -- )
	m.PredefinedWords["set-local"] = GoWord(func(m *Machine) error {
		varName := m.PopValue().(String)
		if len(m.Locals) <= 0 {
			return ErrNoLocalsExist
		}
		m.Locals[len(m.Locals)-1][varName.String()] = m.PopValue()
		return nil
	})
	// ( name -- val )
	m.PredefinedWords["get-local"] = GoWord(func(m *Machine) error {
		varName := m.PopValue().(String)
		if len(m.Locals) <= 0 {
			return ErrNoLocalsExist
		}
		if val, ok := m.Locals[len(m.Locals)-1][varName.String()]; ok {
			m.PushValue(val)
			return nil
		} else {
			fmt.Printf("Can't find %v", varName)
			return ErrLocalNotFound
		}
	})
	m.PredefinedWords["get-locals"] = NilWord(func(m *Machine) {
		m.Locals = append(m.Locals, make(map[string]StackEntry))
	})
	m.PredefinedWords["drop-locals"] = NilWord(func(m *Machine) {
		m.Locals = m.Locals[:len(m.Locals)-1]
	})
	// ( -- locals.. )
	// Run a Quotation for each local
	m.PredefinedWords["each-local"] = GoWord(func(m *Machine) error {
		quot := m.PopValue().(*Quotation)
		code := quot.toCode()
		for key, val := range m.Locals[len(m.Locals)-1] {
			m.PushValue(val)
			m.PushValue(String(key))
			code.Idx = 0
			err := m.execute(code)
			if err != nil {
				return err
			}
		}
		return nil
	})

	// TODO: compress this
	// incr-local-var is here to help ++ be faster
	m.AddGoWord("incr-local-var", "( name -- )", GoWord(func(m *Machine) error {
		varName := m.PopValue().(String)
		if len(m.Locals) <= 0 {
			return ErrNoLocalsExist
		}
		if val, ok := m.Locals[len(m.Locals)-1][varName.String()]; ok {
			// TODO: Clean up this cast
			v, canNumber := val.(Integer)
			if canNumber {
				m.Locals[len(m.Locals)-1][varName.String()] = v + 1
				return nil
			} else {
				return ErrAttemtToIncrementNonNumber
			}
		} else {
			return ErrLocalNotFound
		}
	}))
	m.AddGoWord("decr-local-var", "( name -- )", GoWord(func(m *Machine) error {
		varName := m.PopValue().(String)
		if len(m.Locals) <= 0 {
			return ErrNoLocalsExist
		}
		if val, ok := m.Locals[len(m.Locals)-1][varName.String()]; ok {
			// TODO: Clean up this cast
			v, canNumber := val.(Integer)
			if canNumber {
				m.Locals[len(m.Locals)-1][varName.String()] = v - 1
				return nil
			} else {
				return ErrAttemtToIncrementNonNumber
			}
		} else {
			return ErrLocalNotFound
		}
	}))
	return m.ImportPISCAsset("stdlib/locals.pisc")
}
