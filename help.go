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
	DocString: "A function to look up help for words",
	Load:      loadHelpCore,
}

// TODO: Pull from more sources of docs, like word defs, not just
// :DOC directives
func loadHelpCore(m *machine) error {
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
			fmt.Println("No help available.")
			return nil
		}
	})

	return nil

}
