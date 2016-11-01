package main

import (
	"fmt"
	"time"
)

func (m *machine) loadDebugWords() error {
	// ( -- )
	m.predefinedWords["show-prefix-words"] = NilWord(func(m *machine) {
		for name, _ := range m.prefixWords {
			fmt.Println(name)
		}
	})
	// ( quot -- .. time )
	m.predefinedWords["time"] = GoWord(func(m *machine) error {
		words := &codeQuotation{
			idx:   0,
			words: []word{"call"},
		}
		start := time.Now()
		executeWordsOnMachine(m, words)
		elapsed := time.Since(start)
		m.pushValue(String(fmt.Sprint("Code took ", elapsed)))
		return nil
	})
	return nil
}
