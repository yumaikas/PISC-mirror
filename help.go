package main

import (
	"fmt"
	"strings"
)

// TODO: Indicate deps modules?
var ModHelpCore = PISCModule{
	Author:    "Andrew Owen",
	Name:      "HelpCore",
	License:   "MIT",
	DocString: "TODO: Fill this in",
	Load:      loadHelpModule,
}

func loadHelpModule(m *machine) error {
	return m.loadHelpWords()
}

// These are the help words

func (m *machine) loadHelpWords() error {
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

	return nil

}
