package pisc

import (
	"fmt"
	"strings"
)

// TODO: Indicate deps modules?
var ModHelpCore = Module{
	Author:    "Andrew Owen",
	Name:      "HelpCore",
	License:   "MIT",
	DocString: "A function to look up help for words",
	Load:      loadHelpCore,
}

func _help(m *Machine) error {
	searchTerm := m.PopValue().(String).String()
	if len(m.HelpDocs) > 0 {
		for k, v := range m.HelpDocs {
			if strings.Contains(string(k), searchTerm) || strings.Contains(string(v), searchTerm) {
				fmt.Println("Word: ", k, "\n\tDescription: ", v)
			}
		}
		return nil
	} else {
		fmt.Println("No help available.")
		return nil
	}
}

// TODO: Pull from more Sources of docs, like word defs, not just
// :DOC directive
func loadHelpCore(m *Machine) error {
	// TODO: Implement cli docs browser?

	m.AddGoWordWithStack(
		"help",
		"( search-terim -- )",
		"Search the docs for the given search term, printing out any entries that match",
		_help)
	return nil

}
