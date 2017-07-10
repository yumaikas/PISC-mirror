package pisc

import (
	"fmt"
	"strings"
)

var ModMetaQuoation = Module{
	Author:    "Andrew Owen",
	Name:      "MetaQuotation",
	License:   "MIT",
	DocString: "Words that manipulate Quotations. Use with care",
	Load:      loadQuotWords,
}

func loadQuotWords(m *Machine) error {
	// This is probably a dangerous word...
	m.AddGoWord("quot-set-var", " ( quot key val -- ) ", GoWord(func(m *Machine) error {
		val := m.PopValue()
		key := m.PopValue().String()
		quot := m.PopValue().(*Quotation)
		if !strings.HasPrefix(key, "_") {
			return fmt.Errorf("attempted to set a private local (%v)", key)
		}
		if quot.locals == nil {
			return fmt.Errorf("attempted to add local to word that has no locals")
		}
		quot.locals[key] = val
		return nil
	}))
	m.AddGoWord("quot-has-var", " ( quot key -- ? ) ", GoWord(func(m *Machine) error {
		key := m.PopValue().String()
		quot := m.PopValue().(*Quotation)
		if !strings.HasPrefix(key, "_") {
			return fmt.Errorf("attempted to get a private local (%v)", key)
		}
		if quot.locals == nil {
			m.PushValue(Boolean(false))
			return nil
		}
		_, found := quot.locals[key]
		m.PushValue(Boolean(found))
		return nil
	}))
	m.AddGoWord("quot-get-var", " ( quot key -- val ) ", GoWord(func(m *Machine) error {
		key := m.PopValue().String()
		quot := m.PopValue().(*Quotation)
		if !strings.HasPrefix(key, "_") {
			return fmt.Errorf("attempted to get a private local (%v)", key)
		}
		if quot.locals == nil {
			return fmt.Errorf("Attempted to get a local on a word that has no locals")
		}
		m.PushValue(quot.locals[key])
		return nil
	}))
	return m.ImportPISCAsset("stdlib/with.pisc")
}
