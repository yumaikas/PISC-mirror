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

func _setLocal(m *Machine) error {
	varName := m.PopValue().(String)
	if len(m.Locals) <= 0 {
		return ErrNoLocalsExist
	}
	m.Locals[len(m.Locals)-1][varName.String()] = m.PopValue()
	return nil
}

func _getLocal(m *Machine) error {
	varName := m.PopValue().(String)
	if len(m.Locals) <= 0 {
		return ErrNoLocalsExist
	}
	if val, ok := m.Locals[len(m.Locals)-1][varName.String()]; ok {
		m.PushValue(val)
		return nil
	} else {
		fmt.Printf("ERROR: Can't find %v\n", varName)
		return ErrLocalNotFound
	}
}

func _getLocals(m *Machine) error {
	m.Locals = append(m.Locals, make(map[string]StackEntry))
	return nil
}

func _dropLocals(m *Machine) error {
	m.Locals = m.Locals[:len(m.Locals)-1]
	return nil
}

func _eachLocal(m *Machine) error {
	quot := m.PopValue().(*Quotation)
	for key, val := range m.Locals[len(m.Locals)-1] {
		m.PushValue(val)
		m.PushValue(String(key))
		err := m.CallQuote(quot)
		if IsLoopError(err) && LoopShouldEnd(err) {
			return nil
		}
		if IsLoopError(err) && !LoopShouldEnd(err) {
			continue
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func _incrLocal(m *Machine) error {
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
}

func _decrLocal(m *Machine) error {
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
}

func loadLocalCore(m *Machine) error {
	// Make sure we always have locals
	m.Locals = append(m.Locals, make(map[string]StackEntry))

	m.AddGoWordWithStack(
		"set-local", "( value name -- ) ",
		"Usually handled by :var, but also can be used directly",
		_setLocal)

	m.AddGoWordWithStack(
		"get-local", "( name -- value ) ",
		"Usually handled by $var, but also can be used directly",
		_getLocal)

	m.AddGoWordWithStack(
		"get-locals", "( -- ) ",
		"Pushes a new frame onto the locals stack",
		_getLocals)

	m.AddGoWordWithStack(
		"drop-locals", "( -- ) ",
		"Drops the current locals from the locals stack",
		_dropLocals)

	m.AddGoWordWithStack(
		"each-local", "( quot -- ? ) ",
		"Run a function for each local in the current locals frame",
		_eachLocal)

	m.AddGoWordWithStack(
		"incr-local-var", "( name --  ) ",
		"Increment the integer at :name, if it exists. Error otherwise",
		_incrLocal)

	m.AddGoWordWithStack(
		"decr-local-var", "( name --  ) ",
		"Decrement the integer at :name, if it exists. Error otherwise",
		_decrLocal)
	return m.ImportPISCAsset("stdlib/locals.pisc")
}
