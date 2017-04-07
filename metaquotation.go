package main

import (
	"fmt"
	"strings"
)

var ModMetaQuoation = PISCModule{
	Author:    "Andrew Owen",
	Name:      "MetaQuotation",
	License:   "MIT",
	DocString: "Words that manipulate quotations. Use with care",
	Load:      loadQuotWords,
}

func loadQuotWords(m *machine) error {
	// This is probably a dangerous word...
	m.addGoWord("quot-set-var", " ( quot key val -- ) ", GoWord(func(m *machine) error {
		val := m.popValue()
		key := m.popValue().String()
		quot := m.popValue().(*quotation)
		if !strings.HasPrefix(key, "_") {
			return fmt.Errorf("attempted to set a private local (%v)", key)
		}
		if quot.locals == nil {
			return fmt.Errorf("attempted to add local to word that has no locals")
		}
		quot.locals[key] = val
		return nil
	}))
	m.addGoWord("quot-has-var", " ( quot key -- ? ) ", GoWord(func(m *machine) error {
		key := m.popValue().String()
		quot := m.popValue().(*quotation)
		if !strings.HasPrefix(key, "_") {
			return fmt.Errorf("attempted to get a private local (%v)", key)
		}
		if quot.locals == nil {
			m.pushValue(Boolean(false))
			return nil
		}
		_, found := quot.locals[key]
		m.pushValue(Boolean(found))
		return nil
	}))
	m.addGoWord("quot-get-var", " ( quot key -- val ) ", GoWord(func(m *machine) error {
		key := m.popValue().String()
		quot := m.popValue().(*quotation)
		if !strings.HasPrefix(key, "_") {
			return fmt.Errorf("attempted to get a private local (%v)", key)
		}
		if quot.locals == nil {
			return fmt.Errorf("Attempted to get a local on a word that has no locals")
		}
		m.pushValue(quot.locals[key])
		return nil
	}))
	return m.importPISCAsset("stdlib/with.pisc")
}
