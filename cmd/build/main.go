package main

import (
	"pisc"
	"pisc/libs/shell"
)

func main() {
	m := &pisc.Machine{
		Values:               make([]pisc.StackEntry, 0),
		DefinedWords:         make(map[string]*pisc.CodeQuotation),
		DefinedStackComments: make(map[string]string),
		PredefinedWords:      make(map[string]pisc.GoWord),
		PrefixWords:          make(map[string]*pisc.CodeQuotation),
		HelpDocs:             make(map[string]string),
	}
	m.LoadModules(append(pisc.StandardModules, shell.ModShellUtils)...)
}
