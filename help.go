package main

import (
	"fmt"
	"strings"
)

// These are the help words

func (m *machine) loadHelpWords() {
	// ( val name -- )
	m.predefinedWords["help"] = GoWord(func(m *machine) error {
		searchTerm := m.popValue().(String).String()
		if len(m.helpDocs) > 0 {
			for k, v := range m.helpDocs {
				if strings.Contains(string(k), searchTerm) || strings.Contains(string(v), searchTerm) {
					fmt.Println("Word: ", k, "\n\tDescription: ", v)
				}
			}
			return nil
		} else {
			fmt.Println("No help avaialble.")
			return nil
		}
	})

}
