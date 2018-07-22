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

func _iterateHelpinfo(m *Machine) error {
	perFunction := m.PopValue().(*Quotation)
	perModule := m.PopValue().(*Quotation)
	for _, module := range m.LoadedModules {
		perModuleArgs := Dict{
			"name":    String(module.Name),
			"count":   Integer(len(m.ModuleFunctions[module.Name])),
			"author":  String(module.Author),
			"license": String(module.License),
			"doc":     String(module.DocString),
		}
		m.PushValue(perModuleArgs)
		modErr := m.CallQuote(perModule)
		if modErr != nil {
			return modErr
		}
		for _, funcName := range m.ModuleFunctions[module.Name] {
			perFuncArgs := Dict{
				"name":         String(funcName),
				"stack-effect": String(m.DefinedStackComments[funcName]),
				"doc":          String(m.HelpDocs[funcName]),
			}
			m.PushValue(perFuncArgs)
			funcErr := m.CallQuote(perFunction)
			if funcErr != nil {
				return funcErr
			}
		}
	}
	return nil
}

// TODO: Pull from more Sources of docs, like word defs, not just
// :DOC directive
func loadHelpCore(m *Machine) error {
	// TODO: Implement cli docs browser?

	m.AddGoWordWithStack(
		"help",
		"( search-term:str -- )",
		"Search the docs for the given search term, printing out any entries that match",
		_help)
	m.AddGoWordWithStack(
		"iterate-help-info",
		"( per-module:func[ attrs:dict[ count:int name:str author:str license:str doc:str ] - ] per-function:func[ attrs:dict[ name stack doc ] ] -- ? )",
		"Iterate over all of the installed functions, running the callbacks provided as needed",
		_iterateHelpinfo,
	)
	return m.ImportPISCAsset("stdlib/help.pisc")

}
