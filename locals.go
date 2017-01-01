package main

import (
	"fmt"
	"strings"
)

// Used when parsing words
func isLocalWordPrefix(w word) bool {
	return strings.HasPrefix(string(w), ":") ||
		strings.HasPrefix(string(w), "$")
}

var ErrLocalNotFound = fmt.Errorf("Local variable not found!")
var ErrNoLocalsExist = fmt.Errorf("A local frame hasn't been allocated with get-locals!")

func (m *machine) loadLocalWords() {
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
		quot := m.popValue().(quotation)
		code := &codeQuotation{idx: 0, words: quot.code, codePosition: quot.codePosition}
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
}
