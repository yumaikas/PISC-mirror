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

func _quotSetVar(m *Machine) error {
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
}

func _quotGetVar(m *Machine) error {
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
}

func _quotHasVar(m *Machine) error {
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
}

func loadQuotWords(m *Machine) error {
	// This is probably a dangerous word...
	m.AddGoWordWithStack("quot-set-var",
		" ( quot key val -- ) ",
		"Set local for key in quot to val",
		_quotSetVar)
	m.AddGoWordWithStack(
		"quot-has-var",
		" ( quot key -- ? ) ",
		"Check to see if `quot` has a var a `key`",
		_quotHasVar)
	m.AddGoWordWithStack("quot-get-var",
		" ( quot key -- val ) ",
		"Get get `quot` local value for `key`",
		_quotGetVar)
	return m.ImportPISCAsset("stdlib/with.pisc")
}
